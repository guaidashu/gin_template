package internal

import (
	"fmt"
	"gin_template/app/libs/datetime"
	"gin_template/app/rds"
	"strconv"
	"time"

	"github.com/go-redis/redis"
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
		codeKey:        "gogo:sms:code:%s",
		ipLimitKey:     "gogo:sms:ipLimit:%s",
		mobileLimitKey: "gogo:sms:mobileLimit:%s",
	}
}

func (c *defaultSmsCache) client() *redis.Client {
	return rds.Redis
}

// SetCode 设置短信验证码缓存
func (c *defaultSmsCache) SetCode(mobile, code string) error {
	codeKey := fmt.Sprintf(c.codeKey, mobile)

	_, err := c.client().Pipelined(func(pipe redis.Pipeliner) error {
		pipe.HSet(codeKey, "code", code)
		pipe.HSet(codeKey, "time", time.Now().Unix())
		pipe.HSet(codeKey, "try", 0)
		pipe.Expire(codeKey, 5*time.Minute)
		_, e := pipe.Exec()
		return e
	})
	if err != nil {
		return err
	}

	return nil
}

// GetCodeExpire 获取短信验证码过期时间
func (c *defaultSmsCache) GetCodeExpire(mobile string) (int, error) {
	t, err := c.client().TTL(fmt.Sprintf(c.codeKey, mobile)).Result()
	return int(t), err
}

// GetCode 获取短信验证码
func (c *defaultSmsCache) GetCode(mobile string) (string, error) {
	val, err := c.client().HGet(fmt.Sprintf(c.codeKey, mobile), "code").Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	return val, nil
}

// DelCode 删除短信验证码
func (c *defaultSmsCache) DelCode(mobile string) error {
	_, err := c.client().Del(fmt.Sprintf(c.codeKey, mobile)).Result()
	return err
}

// IncrCodeTry 增加验证码试错次数
func (c *defaultSmsCache) IncrCodeTry(mobile string) (int64, error) {
	return c.client().HIncrBy(fmt.Sprintf(c.codeKey, mobile), "try", 1).Result()
}

// IncrMobileLimit 增加手机号次数限制次数
func (c *defaultSmsCache) IncrMobileLimit(mobile string) (int, error) {
	var (
		deadline       = datetime.TodayLastSecond()
		mobileLimitKey = fmt.Sprintf(c.mobileLimitKey, mobile)
		count          int64
	)

	_, err := c.client().Pipelined(func(pipe redis.Pipeliner) error {
		count = pipe.Incr(mobileLimitKey).Val()
		pipe.ExpireAt(mobileLimitKey, deadline)
		_, e := pipe.Exec()
		return e
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetMobileLimit 获取手机号限制次数
func (c *defaultSmsCache) GetMobileLimit(mobile string) (int, error) {
	countStr, err := c.client().Get(fmt.Sprintf(c.mobileLimitKey, mobile)).Result()
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

	_, err := c.client().Pipelined(func(pipe redis.Pipeliner) error {
		count = pipe.Incr(ipLimitKey).Val()
		pipe.Expire(ipLimitKey, time.Minute)
		_, e := pipe.Exec()
		return e
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetIpLimit 获取IP限制次数
func (c *defaultSmsCache) GetIpLimit(ip string) (int, error) {
	countStr, err := c.client().Get(fmt.Sprintf(c.ipLimitKey, ip)).Result()
	if err != nil {
		return 0, err
	}

	if countStr == "" {
		return 0, nil
	}

	return strconv.Atoi(countStr)
}
