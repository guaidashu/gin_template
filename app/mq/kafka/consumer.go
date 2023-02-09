/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 05/07/2021
 * @Desc: kafka consumer
 */

package kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	defaultCommitInterval = time.Second
	defaultMaxWait        = time.Second
	defaultQueueCapacity  = 1000
)

type (
	ConsumeHandle func(key, value string) error

	ConsumeHandler interface {
		Consume(key, value string) error
	}

	Consumer struct {
		consumer *kafka.Reader
		handle   ConsumeHandler
	}

	innerConsumer struct {
		handle ConsumeHandle
	}
)

func NewConsumer(address []string, topic string, handle ConsumeHandler) *Consumer {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	return &Consumer{consumer: kafka.NewReader(kafka.ReaderConfig{
		Brokers:        address,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		StartOffset:    -2,
		GroupID:        "gin_template_queue",
		MaxWait:        defaultMaxWait,
		CommitInterval: defaultCommitInterval,
		QueueCapacity:  defaultQueueCapacity,
	}),
		handle: handle,
	}
}

func (ic *innerConsumer) Consume(key, value string) error {
	return ic.handle(key, value)
}

func WithHandle(handle ConsumeHandle) ConsumeHandler {
	return &innerConsumer{
		handle: handle,
	}
}
