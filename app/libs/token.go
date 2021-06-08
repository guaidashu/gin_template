/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package libs

import (
	"gin_template/app/config"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

type JwtToken struct {
	secretKey  string
	expireTime int64
}

type jwtCustomClaims struct {
	jwt.StandardClaims
	// 追加自己需要的信息
	UserId int64 `json:"user_id"`
}

var (
	_jwtTokenOnce sync.Once
	_jwtToken     *JwtToken
)

func NewJwtToken() *JwtToken {
	_jwtTokenOnce.Do(func() {
		_jwtToken = &JwtToken{
			secretKey:  config.Config.App.TokenKey,
			expireTime: config.Config.App.TokenExpireTime,
		}
	})

	return _jwtToken
}

// 生成token
// expireTime 应当传入 小时数, 默认是7小时
func (j *JwtToken) GenerateToken(userId int64, secretKey, issuer string) (token string, err error) {
	secretKeyByte := []byte(secretKey)

	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(j.expireTime)).Unix(),
			Issuer:    issuer,
		},
		userId,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if token, err = tokenClaims.SignedString(secretKeyByte); err != nil {
		err = NewReportError(err)
	}

	return
}

func (j *JwtToken) ParseToken(token, secretKey string) (customClaims *jwtCustomClaims, err error) {
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

func (j *JwtToken) GetUidByToken(token, secretKey string) (userId int64, err error) {
	var (
		claims *jwtCustomClaims
	)

	if claims, err = j.ParseToken(token, secretKey); err != nil {
		err = NewReportError(err)
		return
	}

	userId = claims.UserId

	return
}
