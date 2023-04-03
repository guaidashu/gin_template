package ws

import (
	"fmt"
	"gin_template/app/data_struct"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"gin_template/app/libs/random"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type (
	Client struct {
		// 客户端链接
		conn *websocket.Conn
		// 终止标识channel
		stop chan int
		// 要发送给客户端的信息
		send chan *data_struct.WsResponse
		// websocket处理函数
		handle func(name string, data []byte, close func())
		// 等待组
		wait *sync.WaitGroup
		// 此客户端所存在的channel
		channel string
		// 此客户端名(由雪花算法计算而来)
		name string
	}
)

func (c *Client) reader() {
	defer func() {
		_ = recover() // 这句话单纯为了防止崩溃
		_ = c.conn.Close()
		close(c.stop)
		c.wait.Done()
		libs.Logger.Info("close reader! ")
	}()

	// 设置心跳
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				libs.Logger.Info(fmt.Printf("read message error: %v", err))
			}
			break
		}

		go c.handle(c.name, data, c.close)
	}
}

func (c *Client) writer() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		_ = recover() // 防止崩溃
		_ = c.conn.Close()
		ticker.Stop()
		c.wait.Done()
		libs.Logger.Info("close writer! ")
	}()

	for {
		select {
		case send := <-c.send:
			// 正常收到数据并发送
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteJSON(send)
			if err != nil {
				libs.Logger.Error("writer error: ", err)
			}
		case <-c.stop:
			// 如果终止了，告诉客户端终止信息
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case <-ticker.C:
			// 心跳
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) close() {
	_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
	_ = c.conn.Close()
}

func (c *Client) Send(event enum.WsEventEnum, data interface{}) {
	resp := &data_struct.WsResponse{
		Event: string(event),
		Data:  data,
	}

	c.send <- resp
}

func (c *Client) Run(handle func(name string, data []byte, close func())) {
	c.handle = handle

	go c.reader()
	go c.writer()
	c.Send(enum.WsCreateConnect, "success")
	c.wait.Wait()
	NewClientPool().RemoveClient(c.name, c.channel)

	libs.Logger.Info("close client: ", c.name)
}

func NewClient(conn *websocket.Conn) *Client {
	client := &Client{
		conn: conn,
		stop: make(chan int),
		send: make(chan *data_struct.WsResponse),
		wait: &sync.WaitGroup{},
		name: fmt.Sprintf("%v", random.GetSnowflake()),
	}
	client.wait.Add(2)

	return client
}
