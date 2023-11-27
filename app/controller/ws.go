package controller

import (
	"gin_template/app/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: time.Second * 5,
}

type (
	WsController struct {
	}
)

func (c *WsController) WsHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("cant upgrade connection:", err)
		return
	}

	ws.ServeWs(conn, ws.NewWsSrv().Handler)
}
