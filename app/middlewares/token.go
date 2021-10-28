/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package middlewares

import (
	"gin_template/app/enum"
	"gin_template/app/libs"
	jwt2 "gin_template/app/libs/jwt"
	"gin_template/app/libs/serror"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 用户token验证
func ValidateUser() gin.HandlerFunc {
	return validateToken(enum.ValidateUser)
}

// 管理员token验证
func ValidateAdmin() gin.HandlerFunc {
	return validateToken(enum.ValidateAdmin)
}

// 商家token验证
func ValidateMerchant() gin.HandlerFunc {
	return validateToken(enum.ValidateMerchant)
}

func ValidateRefreshToken() gin.HandlerFunc {
	return validateToken(enum.ValidateRefreshToken)
}

func validateToken(validateType int64) gin.HandlerFunc {
	var (
		auth *jwt2.Token
	)

	switch validateType {
	case enum.ValidateUser:
		auth = jwt2.NewJwtToken()
	case enum.ValidateAdmin:
		auth = jwt2.NewJwtAdminToken()
	case enum.ValidateMerchant:
		auth = jwt2.NewJwtMerchantToken()
	case enum.ValidateRefreshToken:
		auth = jwt2.NewJwtRefreshToken()
	}

	return func(ctx *gin.Context) {
		var (
			userId int64
			err    error
		)

		token := ctx.Request.Header.Get("Access-Token")
		if token == "" {
			err = serror.NewErr().SetMsg("未登录")
			libs.Error(ctx, err, http.StatusForbidden)
			return
		}

		if userId, err = auth.GetUidByToken(token); err != nil {
			if e, ok := err.(*jwt.ValidationError); ok {
				switch e.Errors {
				case jwt.ValidationErrorExpired:
					err = serror.NewErr().SetMsg("授权已过期")
					libs.Error(ctx, err, http.StatusForbidden)
					return
				case jwt.ValidationErrorClaimsInvalid,
					jwt.ValidationErrorSignatureInvalid,
					jwt.ValidationErrorNotValidYet,
					jwt.ValidationErrorId,
					jwt.ValidationErrorIssuedAt,
					jwt.ValidationErrorIssuer:
					err = serror.NewErr().SetMsg("无效授权")
					libs.Error(ctx, err, http.StatusForbidden)
					return
				}
			}
		}

		if userId == 0 {
			err = serror.NewErr().SetMsg("无效授权")
			libs.Error(ctx, err, http.StatusForbidden)
			return
		}

		// 设置用户ID
		ctx.Set("userId", userId)

		ctx.Next()
	}
}
