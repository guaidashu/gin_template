/**
  create by yy on 2019-08-23
*/

package controller

import (
	"gin_template/app/libs"
	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	libs.Success(ctx, "index")
}
