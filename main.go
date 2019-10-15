/**
  create by yy on 2019-08-23
*/

package main

import (
	"gin_template/app/ginServer"
	_ "gin_template/app/init"
	_ "gin_template/app/router"
)

func main() {
	ginServer.Run()
}
