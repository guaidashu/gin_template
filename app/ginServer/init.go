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
	var address string

	if len(addr) > 0 || config.Config.App.RunAddress == "" {
		address = resolveAddress(addr)
	} else {
		address = fmt.Sprintf("%v:%v", config.Config.App.RunAddress, config.Config.App.RunPort)
	}

	processed := make(chan struct{}, 1)

	err := startServer(address, processed)

	if err != nil {
		libs.DebugPrint(err.Error())
	}

	<-processed
}

func startServer(address string, processed chan struct{}) error {
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

		// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		ctx, cancel := context.WithTimeout(&ShutdownContext{Chan: nil}, 3*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown failed, err: %v\n", libs.NewReportError(err))
		}
		libs.DebugPrint("server gracefully")

		// 这里的processed  是防止在还未关闭完成的情况下，主程序就down掉了，因为文档说过，关闭途中主程序是不能结束的
		// 但是疑惑的是，我不用 processed来阻塞也不会出错， 目前怀疑是 因为通过子函数来调用，server.ListenAndServe()
		// 会阻塞在此， 然后关闭之后才会触发err 接收error， 并执行后续操作， 若有问题希望能指出
		close(processed)
	}()

	err = server.ListenAndServe()
	// 关闭数据库
	models.CloseDB()

	if http.ErrServerClosed != err {
		err = errors.New(fmt.Sprintf("server not gracefully shutdown, err :%v", err))
	}

	return err
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			libs.DebugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		libs.DebugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
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
	return s.Chan
}

func (s *ShutdownContext) Err() error {
	return nil
}

func (s *ShutdownContext) Value(key interface{}) interface{} {
	return nil
}
