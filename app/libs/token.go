/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package libs

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type jwtCustomClaims struct {
	jwt.StandardClaims

	// 追加自己需要的信息
	Uid int `json:"uid"`
}

// 生成token
func GenerateToken(uid int, secretKey, issuer string, expireTime time.Duration) (token string, err error) {
	secretKeyByte := []byte(secretKey)

	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * expireTime).Unix(),
			Issuer:    issuer,
		},
		uid,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if token, err = tokenClaims.SignedString(secretKeyByte); err != nil {
		err = NewReportError(err)
	}

	return
}

func ParseToken(token, secretKey string) (customClaims *jwtCustomClaims, err error) {
	var (
		tokenClaims *jwt.Token
	)

	if tokenClaims, err = jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secretKey), nil
	}); err != nil {
		err = NewReportError(err)
		return
	}

	if claims, ok := tokenClaims.Claims.(*jwtCustomClaims); ok && tokenClaims.Valid {
		customClaims = claims
	} else {
		err = NewReportError(err)
	}

	return
}

func GetUidByToken(token, secretKey string) (uid int, err error) {
	var (
		claims *jwtCustomClaims
	)

	if claims, err = ParseToken(token, secretKey); err != nil {
		err = NewReportError(err)
		return
	}

	uid = claims.Uid

	return
}
