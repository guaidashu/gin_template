package ws

import (
	"gin_template/app/libs"
	"github.com/gorilla/websocket"
)

func ServeWs(conn *websocket.Conn, handle func(name string, data []byte, close func())) {
	defer func() {
		r := recover()
		if r != nil {
			libs.Logger.Error("serverWs error: ================> ", r)
		}
	}()

	libs.Logger.Info("开始实例化websocket对象")
	// 实例化一个socket对象
	client := NewClient(conn)
	libs.Logger.Info("实例化websocket对象成功")
	pool := NewClientPool()
	libs.Logger.Info("正在建立连接，name为：", client.name, "当前总连接数：", pool.GetCount())
	// 注册到client池
	pool.register <- client
	// 注册处理事件
	client.Run(handle)
}
