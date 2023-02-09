/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 05/07/2021
 * @Desc: desc
 */

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_template/app/config"
	"gin_template/app/data_struct/_interface"
	"gin_template/app/enum"
	"gin_template/app/libs"
)

var Kafka *Kq

type (
	Kq struct {
		producers map[string]*Producer
		consumers map[string]*Consumer
	}

	KqInit struct{}
)

func NewKqInit() *KqInit {
	return &KqInit{}
}

func (k *KqInit) Init(*_interface.ServiceParam) error {
	InitKafka()
	return nil
}

func (k *KqInit) ComponentName() enum.BootModuleType {
	return enum.KafkaInit
}

func (k *KqInit) Close() error {
	Kafka.Stop()
	return nil
}

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
				// m, err := k.consumers[key].consumer.ReadMessage(context.Background())
				m, err := k.consumers[key].consumer.FetchMessage(context.Background())
				if err != nil {
					break
				}

				err = k.consumers[key].handle.Consume(string(m.Key), string(m.Value))
				if err != nil {
					libs.Logger.Error(err.Error())
					continue
				}

				// 当用了group的时候, ReadMessage函数会自动执行 CommitMessages函数提交事务,
				// 这样不安全,所以,使用FetchMessage
				_ = k.consumers[key].consumer.CommitMessages(context.Background(), m)
			}
		})

	}
}

func (k *Kq) Start() {
	k.startConsumer()
}

func (k *Kq) Stop() {
	for key := range k.producers {
		_ = k.producers[key].producer.Close()
	}

	for key := range k.consumers {
		_ = k.consumers[key].consumer.Close()
	}
}
