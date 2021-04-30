/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ValidateToken(ctx *gin.Context) {
	fmt.Println("validate token")

	ctx.Next()
}
