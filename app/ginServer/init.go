/**
  create by yy on 2019-07-02
*/

package ginServer

import (
	"fmt"
	"gin_template/app/config"
	"github.com/gin-gonic/gin"
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
	_ = Router.Run(address)
}
