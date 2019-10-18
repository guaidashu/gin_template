/**
  create by yy on 2019-07-02
*/

package ginServer

import (
	"context"
	"errors"
	"fmt"
	"gin_template/app/config"
	"gin_template/app/libs"
	"gin_template/app/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	Router *gin.Engine
)

func init() {
	Router = gin.Default()
}

func Run(addr ...string) {
	if len(addr) > 0 || config.Config.App.RunAddress == "" {
		_ = Router.Run(addr...)
		return
	}
	address := fmt.Sprintf("%v:%v", config.Config.App.RunAddress, config.Config.App.RunPort)
	err := startServer(address)
	if err != nil {
		libs.DebugPrint(libs.GetErrorString(err))
	}
}

func startServer(address string) error {
	var err error
	server := &http.Server{
		Addr:    address,
		Handler: Router,
	}
	libs.DebugPrint("Listening and serving HTTP on %s\n", address)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		ctx, cancel := context.WithTimeout(&ShutdownContext{}, 3*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown failed, err: %v\n", libs.NewReportError(err))
		}
		libs.DebugPrint("server gracefully")
	}()
	err = server.ListenAndServe()
	// 关闭数据库
	models.CloseDB()

	if http.ErrServerClosed != err {
		err = errors.New(fmt.Sprintf("server not gracefully shutdown, err :%v", err))
	}

	return err
}

// 实现 context.Context 接口
// Deadline() (deadline time.Time, ok bool)
// Done() <-chan struct{}
// Err() error
// Value(key interface{}) interface{}
// 让 http 优雅退出(graceful)
type ShutdownContext struct {
	Chan         chan struct{}
	DeadLineTime time.Time
}

func (s *ShutdownContext) Deadline() (deadline time.Time, ok bool) {
	deadline = s.DeadLineTime
	ok = true
	return
}

func (s *ShutdownContext) Done() <-chan struct{} {
	return nil
}

func (s *ShutdownContext) Err() error {
	return nil
}

func (s *ShutdownContext) Value(key interface{}) interface{} {
	return nil
}
