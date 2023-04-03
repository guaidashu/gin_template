package service

import (
	"encoding/json"
	"gin_template/app/data_struct"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"gin_template/app/ws"
	"sync"
)

type (
	WsSrv interface {
		// 总处理handler
		Handler(name string, data []byte, close func())
	}

	defaultWsSrv struct {
		clientPool *ws.ClientPool
	}
)

var (
	_WsSrv     WsSrv
	_WsSrvOnce sync.Once
)

func NewWsSrv() WsSrv {
	_WsSrvOnce.Do(func() {
		_WsSrv = &defaultWsSrv{
			clientPool: ws.NewClientPool(),
		}
	})

	return _WsSrv
}

func (s *defaultWsSrv) Handler(name string, data []byte, close func()) {
	req := &data_struct.WsRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		libs.Logger.Error(err)
		return
	}

	switch enum.WsEventEnum(req.Event) {
	case enum.WsSubscribeEvent:
		// 初次订阅，需要写入channel
		wsData := &data_struct.WsSubscribe{}
		err = json.Unmarshal([]byte(req.Data), wsData)
		if err != nil {
			s.clientPool.Get(name).Send(enum.WsSubscribeEvent, "错误的订阅格式")
			return
		}

		// 写入channel
		s.clientPool.SetChannel(name, wsData.Channel)
		s.clientPool.Get(name).Send(enum.WsSubscribeEvent, "ok")
	case enum.WsCloseEvent: // 关闭连接
		close()
	case enum.WsPongEvent:
		s.clientPool.Get(name).Send(enum.WsPongEvent, "ok")
	}
}

