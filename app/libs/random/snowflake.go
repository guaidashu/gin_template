package random

import (
	"context"
	"google.golang.org/appengine/log"
	"sync"
	"time"
)

const (
	epoch             = int64(1577808000000)                           // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits     = uint(41)                                       // 时间戳占用位数
	datacenteridBits  = uint(2)                                        // 数据中心id所占位数
	workeridBits      = uint(7)                                        // 机器id所占位数
	sequenceBits      = uint(12)                                       // 序列所占的位数
	timestampMax      = int64(-1 ^ (-1 << timestampBits))              // 时间戳最大值
	datacenteridMax   = int64(-1 ^ (-1 << datacenteridBits))           // 支持的最大数据中心id数量
	workeridMax       = int64(-1 ^ (-1 << workeridBits))               // 支持的最大机器id数量
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))               // 支持的最大序列id数量
	workeridShift     = sequenceBits                                   // 机器id左移位数
	datacenteridShift = sequenceBits + workeridBits                    // 数据中心id左移位数
	timestampShift    = sequenceBits + workeridBits + datacenteridBits // 时间戳左移位数
)

type (
	Snowflake struct {
		sync.Mutex         // 锁
		timestamp    int64 // 时间戳 ，毫秒
		workerId     int64 // 工作节点
		datacenterId int64 // 数据中心机房id
		sequence     int64 // 序列号
	}
)

var (
	_snowflake     *Snowflake
	_snowflakeOnce sync.Once
)

func NewSnowflake() *Snowflake {
	_snowflakeOnce.Do(func() {
		_snowflake = &Snowflake{
			Mutex: sync.Mutex{},
		}
	})

	return _snowflake
}

func (s *Snowflake) NextVal() int64 {
	s.Lock()
	defer func() {
		s.Unlock()
	}()

	now := time.Now().UnixNano() / 1000000 // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}

	t := now - epoch
	if t > timestampMax {
		log.Errorf(context.Background(), "epoch must be between 0 and %d", timestampMax-1)
		return 0
	}
	s.timestamp = now
	r := (t)<<timestampShift | (s.datacenterId << datacenteridShift) | (s.workerId << workeridShift) | (s.sequence)

	return r
}

func GetSnowflake() int64 {
	return NewSnowflake().NextVal()
}
