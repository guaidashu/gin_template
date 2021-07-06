/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 05/07/2021
 * @Desc: desc
 */

package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_template/app/config"
	"gin_template/app/enum"
	"gin_template/app/libs"
)

var Kafka *Kq

type (
	Kq struct {
		producers map[string]*Producer
		consumers map[string]*Consumer
	}
)

func InitKafka() {
	Kafka = NewKq()
	Kafka.SetKafka(enum.UserTopic)
	Kafka.RegKafkaConsumer(enum.UserTopic, func(key, value string) error {
		fmt.Println(key, value)
		return nil
	})

	Kafka.Start()
}

func NewKq() *Kq {
	return &Kq{
		producers: make(map[string]*Producer),
		consumers: make(map[string]*Consumer),
	}
}

// 通过topic设置kafka writer
func (k *Kq) SetKafka(topic string) {
	k.producers[topic] = NewProducer(config.Config.Kafka.Hosts, topic)
}

func (k *Kq) RegKafkaConsumer(topic string, handle func(key, value string) error) {
	if len(config.Config.Kafka.Hosts) == 0 {
		panic("无可用的节点信息")
	}

	k.consumers[topic] = NewConsumer(config.Config.Kafka.Hosts, topic, WithHandle(handle))
}

// 发送消息
func (k *Kq) SendMsg(v interface{}, topic string) (err error) {
	var (
		data []byte
	)

	data, err = json.Marshal(v)
	if err != nil {
		return
	}

	return k.producers[topic].Push(string(data))
}

func (k *Kq) startConsumer() {
	for key := range k.consumers {
		libs.RunSafe(func() {
			for {
				m, err := k.consumers[key].consumer.ReadMessage(context.Background())
				if err != nil {
					break
				}

				err = k.consumers[key].handle.Consume(string(m.Key), string(m.Value))
				if err != nil {
					libs.Logger.Error(err.Error())
					continue
				}

				_ = k.consumers[key].consumer.CommitMessages(context.Background(), m)
			}
		})

	}
}

func (k *Kq) Start() {
	k.startConsumer()
}
