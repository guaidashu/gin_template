package ws

import (
	"fmt"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"sync"
)

var (
	_clientPoolOnce sync.Once
	_clientPool     *ClientPool
)

type (
	// 客户端连接池
	ClientPool struct {
		// 客户端连接池的读写锁
		clientsLock sync.RWMutex
		// 客户端连接池
		clients map[string]*Client
		// channel客户端连接池读写锁
		clientsChannelLock sync.RWMutex
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
			c.clientsLock.Lock()
			c.clients[register.name] = register
			libs.Logger.Info("成功建立连接，name为：", register.name, "成功后的总连接数：", len(c.clients))
			c.clientsLock.Unlock()
		}
	}
}

func (c *ClientPool) unregisterClient() {
	for {
		select {
		case unregister := <-c.unregister:
			c.clientsLock.Lock()
			c.removeUnregister(unregister.name)
			c.clientsLock.Unlock()
		}
	}
}

func (c *ClientPool) removeUnregister(name string) {
	defer func() {
		rc := recover()
		if rc != nil {
			libs.Logger.Error("removeUnregister ====================> 移除client失败")
		}
	}()

	delete(c.clients, name)
}

func (c *ClientPool) Get(name string) *Client {
	c.clientsLock.RLock()
	defer func() {
		c.clientsLock.RUnlock()
	}()

	if _, ok := c.clients[name]; !ok {
		return nil
	}

	return c.clients[name]
}

func (c *ClientPool) GetChannel(channel enum.WsChannelEnum) map[string]*Client {
	c.clientsChannelLock.RLock()
	defer func() {
		c.clientsChannelLock.RUnlock()
	}()

	if _, ok := c.clientsChannel[string(channel)]; !ok {
		return nil
	}

	return c.clientsChannel[string(channel)]
}

func (c *ClientPool) GetByChannel(name string, channel enum.WsChannelEnum) *Client {
	c.clientsChannelLock.RLock()
	defer func() {
		c.clientsChannelLock.RUnlock()
	}()

	if _, ok := c.clientsChannel[string(channel)]; !ok {
		return nil
	}

	return c.clientsChannel[string(channel)][name]
}

// 设置client到channel的map连接池子
func (c *ClientPool) SetChannel(name, channel string) {
	c.clientsChannelLock.Lock()
	defer func() {
		c.clientsChannelLock.Unlock()
	}()

	if _, ok := c.clientsChannel[channel]; !ok {
		c.clientsChannel[channel] = make(map[string]*Client)
	}

	c.clientsLock.RLock()
	defer func() {
		c.clientsLock.RUnlock()
	}()
	if _, ok := c.clients[name]; !ok {
		return
	}

	c.clients[name].channel = append(c.clients[name].channel, channel)
	c.clientsChannel[channel][name] = c.clients[name]
}

// 结束项目时执行
func (c *ClientPool) CloseAllClients() {
	c.clientsLock.Lock()
	defer func() {
		c.clientsLock.Unlock()
	}()

	for _, v := range c.clients {
		v.close()
	}
}

func (c *ClientPool) RemoveClient(name string, channels []string) {
	c.clientsLock.RLock()
	defer func() {
		c.clientsLock.RUnlock()
	}()

	c.clientsChannelLock.Lock()
	defer func() {
		c.clientsChannelLock.Unlock()
	}()

	// 移除channel池的client
	for _, channel := range channels {
		if _, ok := c.clientsChannel[channel]; !ok {
			continue
		}

		if _, ok := c.clientsChannel[channel][name]; !ok {
			continue
		}

		delete(c.clientsChannel[channel], name)
	}

	if _, ok := c.clients[name]; !ok {
		return
	}
	// 移除clients的client
	c.unregister <- c.clients[name]
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

func (c *ClientPool) GetCount() int {
	c.clientsLock.RLock()
	defer func() {
		c.clientsLock.RUnlock()
	}()

	return len(c.clients)
}
