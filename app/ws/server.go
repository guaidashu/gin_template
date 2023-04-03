package ws

import (
	"github.com/gorilla/websocket"
)

func ServeWs(conn *websocket.Conn, handle func(name string, data []byte, close func())) {
	// 实例化一个socket对象
	client := NewClient(conn)
	// 注册到client池
	NewClientPool().register <- client
	// 注册处理事件
	client.Run(handle)
}
