package enum

type WsEventEnum string

// websocket相关事件
const (
	WsCreateConnect       WsEventEnum = "ws_create_connect"         // 建立连接
	WsSubscribeEvent      WsEventEnum = "ws_subscribe_event"        // 订阅
	WsCloseEvent          WsEventEnum = "ws_close_event"            // 关闭消息
	WsPongEvent           WsEventEnum = "ws_pong_event"             // 心跳
)

type WsChannelEnum string

// channel 例子
const (
	WsExampleChannel WsChannelEnum = "ws_example_channel"
)

