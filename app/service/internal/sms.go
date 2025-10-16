package internal

import (
	"context"
	"fmt"
	"gin_template/app/libs/datetime"
	"gin_template/app/rds"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	SmsCache interface {
		// SetCode 设置短信验证码缓存
		SetCode(mobile, code string) error
		// GetCode 获取短信验证码
		GetCode(mobile string) (string, error)
		// DelCode 删除短信验证码
		DelCode(mobile string) error
		// GetCodeExpire 获取短信验证码过期时间(秒)
		GetCodeExpire(mobile string) (int, error)
		// IncrCodeTry 增加验证码试错次数
		IncrCodeTry(mobile string) (int64, error)
		// IncrMobileLimit 增加手机号次数限制次数
		IncrMobileLimit(mobile string) (int, error)
		// GetMobileLimit 获取手机号限制次数
		GetMobileLimit(mobile string) (int, error)
		// IncrIpLimit 增加IP限制次数
		IncrIpLimit(ip string) (int, error)
		// GetIpLimit 获取IP限制次数
		GetIpLimit(ip string) (int, error)
	}

	defaultSmsCache struct {
		abstract
		codeKey        string // 短信验证码(%s:mobile)
		ipLimitKey     string // IP限制(%s:ip)
		mobileLimitKey string // 手机号每天次数限制(%s:mobile)
	}
)

// NewSmsCache 创建短信缓存类
func NewSmsCache() SmsCache {
	return &defaultSmsCache{
		codeKey:        "gin_template:sms:code:%s",
		ipLimitKey:     "gin_template:sms:ipLimit:%s",
		mobileLimitKey: "gin_template:sms:mobileLimit:%s",
	}
}

func (c *defaultSmsCache) client() *redis.Client {
	return rds.Redis
}

// SetCode 设置短信验证码缓存
func (c *defaultSmsCache) SetCode(mobile, code string) error {
	codeKey := fmt.Sprintf(c.codeKey, mobile)
	ctx := context.Background()

	_, err := c.client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, codeKey, "code", code)
		pipe.HSet(ctx, codeKey, "time", time.Now().Unix())
		pipe.HSet(ctx, codeKey, "try", 0)
		pipe.Expire(ctx, codeKey, 5*time.Minute)
		_, e := pipe.Exec(ctx)
		return e
	})
	if err != nil {
		return err
	}

	return nil
}

// GetCodeExpire 获取短信验证码过期时间
func (c *defaultSmsCache) GetCodeExpire(mobile string) (int, error) {
	ctx := context.Background()
	t, err := c.client().TTL(ctx, fmt.Sprintf(c.codeKey, mobile)).Result()
	return int(t), err
}

// GetCode 获取短信验证码
func (c *defaultSmsCache) GetCode(mobile string) (string, error) {
	ctx := context.Background()
	val, err := c.client().HGet(ctx, fmt.Sprintf(c.codeKey, mobile), "code").Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	return val, nil
}

// DelCode 删除短信验证码
func (c *defaultSmsCache) DelCode(mobile string) error {
	ctx := context.Background()
	_, err := c.client().Del(ctx, fmt.Sprintf(c.codeKey, mobile)).Result()
	return err
}

// IncrCodeTry 增加验证码试错次数
func (c *defaultSmsCache) IncrCodeTry(mobile string) (int64, error) {
	ctx := context.Background()
	return c.client().HIncrBy(ctx, fmt.Sprintf(c.codeKey, mobile), "try", 1).Result()
}

// IncrMobileLimit 增加手机号次数限制次数
func (c *defaultSmsCache) IncrMobileLimit(mobile string) (int, error) {
	var (
		deadline       = datetime.TodayLastSecond()
		mobileLimitKey = fmt.Sprintf(c.mobileLimitKey, mobile)
		count          int64
	)
	ctx := context.Background()

	_, err := c.client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		count = pipe.Incr(ctx, mobileLimitKey).Val()
		pipe.ExpireAt(ctx, mobileLimitKey, deadline)
		_, e := pipe.Exec(ctx)
		return e
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetMobileLimit 获取手机号限制次数
func (c *defaultSmsCache) GetMobileLimit(mobile string) (int, error) {
	ctx := context.Background()
	countStr, err := c.client().Get(ctx, fmt.Sprintf(c.mobileLimitKey, mobile)).Result()
	if err != nil {
		return 0, err
	}

	if countStr == "" {
		return 0, nil
	}

	return strconv.Atoi(countStr)
}

// IncrIpLimit 增加IP限制次数
func (c *defaultSmsCache) IncrIpLimit(ip string) (int, error) {
	var (
		ipLimitKey = fmt.Sprintf(c.ipLimitKey, ip)
		count      int64
	)
	ctx := context.Background()

	_, err := c.client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		count = pipe.Incr(ctx, ipLimitKey).Val()
		pipe.Expire(ctx, ipLimitKey, time.Minute)
		_, e := pipe.Exec(ctx)
		return e
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetIpLimit 获取IP限制次数
func (c *defaultSmsCache) GetIpLimit(ip string) (int, error) {
	ctx := context.Background()
	countStr, err := c.client().Get(ctx, fmt.Sprintf(c.ipLimitKey, ip)).Result()
	if err != nil {
		return 0, err
	}

	if countStr == "" {
		return 0, nil
	}

	return strconv.Atoi(countStr)
}
