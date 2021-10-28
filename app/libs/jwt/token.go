/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package jwt

import (
	"gin_template/app/config"
	"gin_template/app/libs"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	secretKey  string // 密钥
	expireTime int64  // 过期时间 单位是小时
	Issuer     string // 发行人
}

type CustomClaims struct {
	jwt.StandardClaims
	// 追加自己需要的信息
	UserId  int64 `json:"user_id"`
	GenType int64 `json:"gen_type"`
}

var (
	_jwtTokenOnce         sync.Once
	_jwtToken             *Token
	_jwtAdminTokenOnce    sync.Once
	_jwtAdminToken        *Token
	_jwtMerchantTokenOnce sync.Once
	_jwtMerchantToken     *Token
)

func NewJwtToken() *Token {
	_jwtTokenOnce.Do(func() {
		_jwtToken = &Token{
			secretKey:  config.Config.App.TokenKey,
			expireTime: config.Config.App.TokenExpireTime,
			Issuer:     "guaidashu",
		}
	})

	return _jwtToken
}

func NewJwtAdminToken() *Token {
	_jwtAdminTokenOnce.Do(func() {
		_jwtAdminToken = &Token{
			secretKey:  config.Config.App.AdminTokenKey,
			expireTime: config.Config.App.TokenExpireTime,
			Issuer:     "guaidashu",
		}
	})

	return _jwtAdminToken
}

func NewJwtMerchantToken() *Token {
	_jwtMerchantTokenOnce.Do(func() {
		_jwtMerchantToken = &Token{
			secretKey:  config.Config.App.TokenMerchantKey,
			expireTime: config.Config.App.TokenExpireTime,
			Issuer:     "guaidashu",
		}
	})

	return _jwtMerchantToken
}

func NewJwtRefreshToken() *Token {
	_jwtMerchantTokenOnce.Do(func() {
		_jwtMerchantToken = &Token{
			secretKey:  config.Config.App.RefreshTokenKey,
			expireTime: config.Config.App.RefreshTokenExpireTime,
			Issuer:     "guaidashu",
		}
	})

	return _jwtMerchantToken
}

// 生成token
// expireTime 应当传入 小时数, 默认是7小时
func (j *Token) GenerateToken(userId int64, generateType ...int64) (token string, claims *CustomClaims, err error) {
	var (
		genType int64
	)

	if len(generateType) > 0 {
		genType = generateType[0]
	}

	secretKeyByte := []byte(j.secretKey)
	claims = &CustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(j.expireTime)).Unix(),
			Issuer:    j.Issuer,
		},
		userId,
		genType,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err = tokenClaims.SignedString(secretKeyByte); err != nil {
		err = libs.NewReportError(err)
	}

	return
}

func (j *Token) ParseToken(token string) (customClaims *CustomClaims, err error) {
	var (
		tokenClaims *jwt.Token
	)

	if tokenClaims, err = jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (
		i interface{}, err error) {
		return []byte(j.secretKey), nil
	}); err != nil {
		err = libs.NewReportError(err)
		return
	}

	if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
		customClaims = claims
	} else {
		err = libs.NewReportError(err)
	}

	return
}

func (j *Token) GetUidByToken(token string) (userId int64, err error) {
	var (
		claims *CustomClaims
	)

	if claims, err = j.ParseToken(token); err != nil {
		err = libs.NewReportError(err)
		return
	}

	userId = claims.UserId

	return
}
