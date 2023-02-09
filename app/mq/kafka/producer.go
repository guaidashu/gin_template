/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 05/07/2021
 * @Desc: kafka producer
 */

package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"strconv"
	"time"
)

type (
	Producer struct {
		producer *kafka.Writer
	}
)

func NewProducer(address []string, topic string) *Producer {
	return &Producer{producer: kafka.NewWriter(kafka.WriterConfig{
		Brokers:          address,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		CompressionCodec: snappy.NewCompressionCodec(),
	})}
}

func (p *Producer) Push(v string) error {
	return p.producer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.FormatInt(time.Now().UnixNano(), 10)),
		Value: []byte(v),
	})
}
