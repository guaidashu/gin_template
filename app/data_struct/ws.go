package data_struct

type (
	WsRequest struct {
		Event string `json:"event"` // 事件
		Data  string `json:"data"`  // 数据
	}

	WsSubscribe struct {
		Channel string `json:"channel"` // 要订阅的channel
	}

	WsResponse struct {
		Event string      `json:"event"` // 事件
		Data  interface{} `json:"data"`  // 数据
	}
)

