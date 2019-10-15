/**
  create by yy on 2019-08-23
*/

package main

import (
	"gin_template_origin/app/ginServer"
	_ "gin_template_origin/app/init"
	_ "gin_template_origin/app/router"
)

func main() {
	ginServer.Run()
}
