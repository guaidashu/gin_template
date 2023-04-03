package ws

import (
	"gin_template/app/enum"
	"fmt"
	"sync"
)

var (
	_clientPoolOnce sync.Once
	_clientPool     *ClientPool
)

type (
	// 客户端连接池
	ClientPool struct {
		// 客户端连接池
		clients map[string]*Client
		// channel客户端连接池 (精确到每一个channel组)， 暂时不启用
		clientsChannel map[string]map[string]*Client
		// 发送的数据
		send chan []byte
		// 注册客户端(存入连接池)
		register chan *Client
		// 销毁客户端(从连接池移除)
		unregister chan *Client
	}
)

func NewClientPool() *ClientPool {
	_clientPoolOnce.Do(func() {
		_clientPool = &ClientPool{
			clients:        make(map[string]*Client),
			clientsChannel: make(map[string]map[string]*Client),
			send:           make(chan []byte),
			register:       make(chan *Client),
			unregister:     make(chan *Client),
		}

		go _clientPool.registerClient()
		go _clientPool.unregisterClient()
	})

	return _clientPool
}

func (c *ClientPool) registerClient() {
	for {
		select {
		case register := <-c.register:
			c.clients[register.name] = register
		}
	}
}

func (c *ClientPool) unregisterClient() {
	for {
		select {
		case unregister := <-c.unregister:
			delete(c.clients, unregister.name)
		}
	}
}

func (c *ClientPool) Get(name string) *Client {
	return c.clients[name]
}

func (c *ClientPool) GetChannel(channel enum.WsChannelEnum) map[string]*Client {
	return c.clientsChannel[string(channel)]
}

func (c *ClientPool) GetByChannel(name string, channel enum.WsChannelEnum) *Client {
	if _, ok := c.clientsChannel[string(channel)]; !ok {
		return nil
	}

	return c.clientsChannel[string(channel)][name]
}

// 设置client到channel的map连接池子
func (c *ClientPool) SetChannel(name, channel string) {
	if _, ok := c.clientsChannel[channel]; !ok {
		c.clientsChannel[channel] = make(map[string]*Client)
	}

	c.clients[name].channel = channel
	c.clientsChannel[channel][name] = c.clients[name]
}

// 结束项目时执行
func (c *ClientPool) CloseAllClients() {
	for _, v := range c.clients {
		v.close()
	}
}

func (c *ClientPool) RemoveClient(name string, channel string) {
	// 移除clients的client
	c.unregister <- c.clients[name]

	// 移除channel池的client
	if _, ok := c.clientsChannel[channel]; !ok {
		return
	}
	delete(c.clientsChannel[channel], name)
}

// 广播
func (c *ClientPool) Broadcast() {
	for {
		select {
		case send := <-c.send:
			go c.broadcast(send)
		}
	}
}

func (c *ClientPool) broadcast(data []byte) {
	fmt.Println("广播数据：", string(data))
}
