package data_struct

import "gin_template/app/enum"

type (
	WsRequest struct {
		Event enum.WsEventEnum `json:"event"` // 事件
		Token string           `json:"token"` // token
		Data  string           `json:"data"`  // 数据
		WsId  string           `json:"ws_id"` // ws id
	}

	WsSubscribe struct {
		Channel  string `json:"channel"`  // 要订阅的channel
		Username string `json:"username"` // 用户名
		Password string `json:"password"` // 密码
	}

	WsResponse struct {
		Event enum.WsEventEnum `json:"event"` // 事件
		Data  interface{}      `json:"data"`  // 数据
	}
)
