package ws

import (
	"encoding/json"
	"gin_template/app/data_struct"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"gin_template/app/libs/serror"
	"sync"
)

type (
	WsSrv interface {
		// 总处理handler
		Handler(name string, data []byte, close func())
		// 注册路由
		Register(eventHandler *WsHandler) *WsEventHandler
	}

	defaultWsSrv struct {
		clientPool *ClientPool
		handlers   map[enum.WsEventEnum]*WsEventHandler
		lock       sync.RWMutex
	}
)

var (
	_WsSrv     WsSrv
	_WsSrvOnce sync.Once
)

func NewWsSrv() WsSrv {
	_WsSrvOnce.Do(func() {
		_WsSrv = &defaultWsSrv{
			clientPool: NewClientPool(),
			handlers:   make(map[enum.WsEventEnum]*WsEventHandler),
		}
	})

	return _WsSrv
}

// name 为ws链接的唯一ID，由雪花算法生成
// data 为客户端发来的 源数据
func (s *defaultWsSrv) Handler(name string, data []byte, close func()) {
	req := &data_struct.WsRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		libs.Logger.Error(err)
		return
	}

	// 不需要认证的
	// 连通逻辑
	switch req.Event {
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
		return
	case enum.WsCloseEvent: // 关闭连接
		close()
		return
	case enum.WsPongEvent:
		s.clientPool.Get(name).Send(enum.WsPongEvent, "ok")
		return
	}

	req.WsId = name
	// 以下相当于是路由了
	s.handler(req)
}

func (s *defaultWsSrv) Register(eventHandler *WsHandler) *WsEventHandler {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()

	libs.DebugPrint("ws监听：--> %v 事件监听开始", eventHandler.EventName)
	handler := &WsEventHandler{
		Handler:     eventHandler,
		Middlewares: make([]WsHandlerMiddleware, 0),
	}
	s.handlers[eventHandler.EventName] = handler

	return handler
}

func (s *defaultWsSrv) handler(req *data_struct.WsRequest) {
	client := s.clientPool.Get(req.WsId)
	ctx := NewContext()
	ctx.SetClient(req.WsId, client)
	ctx.SetData([]byte(req.Data))
	ctx.SetEvent(req.Event)

	s.lock.RLock()
	defer func() {
		s.lock.RUnlock()
	}()

	// 如果没有注册，则不处理
	if _, ok := s.handlers[req.Event]; !ok {
		ctx.Error(serror.NewErr().SetErr(serror.ErrMethodNotExist))
		return
	}

	ctx.Set("token", req.Token)
	// 注册到了handlers里，取出进行调用
	handler := s.handlers[req.Event]

	// 遍历执行中间件
	for _, v := range handler.Middlewares {
		err := v(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	// 调用最终注册的方法
	handler.Handler.Handler(ctx)
}
