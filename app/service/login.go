/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/10/2021
 * @Desc: TODO 文件说明
 */

package service

import (
	"errors"
	"gin_template/app/config"
	"gin_template/app/data_struct"
	"gin_template/app/data_struct/requests"
	"gin_template/app/data_struct/responses"
	"gin_template/app/enum"
	"gin_template/app/libs/jwt"
	"gin_template/app/libs/miniprogram"
	"gin_template/app/libs/serror"
	"gin_template/app/models"
	"sync"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/guaidashu/go_helper/crypto_tool"
	"gorm.io/gorm"
)

type (
	LoginSrv interface {
		// 用户登录
		Login(code string) (token *responses.LoginResp, err error)
		// 管理员登录
		AdminLogin(req *requests.AdminLoginReq) (token *responses.LoginResp, err error)
		// 获取token
		GetToken(req *requests.GetTokenReq) (token *responses.GetTokenResp, err error)
		// 获取刷新token
		GetRefreshToken(req *requests.GetTokenReq) (token *responses.GetRefreshTokenResp, err error)
	}

	defaultLoginSrv struct {
	}
)

var (
	_loginSrv     LoginSrv
	_loginSrvOnce sync.Once
)

func NewLoginSrv() LoginSrv {
	_loginSrvOnce.Do(func() {
		_loginSrv = &defaultLoginSrv{}
	})

	return _loginSrv
}

// 小程序登录
func (s *defaultLoginSrv) Login(code string) (token *responses.LoginResp, err error) {
	var (
		openInfo *data_struct.MiniProGramLoginInfo
		user     *models.UserModel
	)

	if openInfo, err = miniprogram.GetOpenId(code); err != nil {
		err = serror.NewErr().SetErr(err)
		return
	}

	if openInfo.OpenId == "" {
		err = serror.NewErr().SetErr("got empty openid")
		return
	}

	// 否则进行判断，首先通过open id 获取用户信息
	user, err = models.NewUserModel().GetUserByOpenId(openInfo.OpenId)
	switch err {
	case gorm.ErrRecordNotFound:
		return &responses.LoginResp{
			Status: 1,
			OpenId: openInfo.OpenId,
		}, nil
	case nil:
	default:
		err = serror.NewErr().SetErr(err)
		return
	}

	tokenStr, claims, generateErr := jwt.NewJwtToken().GenerateToken(user.UserId)
	if generateErr != nil {
		err = serror.NewErr().SetErr(generateErr)
		return
	}
	refreshToken, _, generateErr := jwt.NewJwtRefreshToken().GenerateToken(user.UserId, enum.ValidateUser)
	if generateErr != nil {
		err = serror.NewErr().SetErr(generateErr)
		return
	}

	token = &responses.LoginResp{
		Token:        tokenStr,
		RefreshToken: refreshToken,
		ExpireAt:     claims.ExpiresAt,
	}

	return
}

func (s *defaultLoginSrv) AdminLogin(req *requests.AdminLoginReq) (token *responses.LoginResp, err error) {
	// 验证手机号和密码
	// 查询管理员信息
	user, queryErr := models.NewManageUserModel().GetByPhoneNumber(req.PhoneNumber)
	if queryErr != nil {
		err = serror.NewErr().SetMsg("此账号不存在,请联系超级管理员").SetErr(queryErr)
		return
	}

	if crypto_tool.Md5(req.Password) != user.Password {
		err = serror.NewErr().SetMsg("账号或密码错误").SetErr(queryErr)
		return
	}

	tokenStr, claims, generateErr := jwt.NewJwtAdminToken().GenerateToken(user.ManageUserId)
	if generateErr != nil {
		err = serror.NewErr().SetErr(generateErr)
		return
	}
	refreshToken, _, generateErr := jwt.NewJwtRefreshToken().GenerateToken(user.ManageUserId, enum.ValidateAdmin)
	if generateErr != nil {
		err = serror.NewErr().SetErr(generateErr)
		return
	}

	token = &responses.LoginResp{
		Token:        tokenStr,
		RefreshToken: refreshToken,
		ExpireAt:     claims.ExpiresAt,
	}

	return
}

func (s *defaultLoginSrv) GetToken(req *requests.GetTokenReq) (token *responses.GetTokenResp, err error) {
	var (
		auth              *jwt.Token
		claims            *jwt.CustomClaims
		refreshClaims     *jwt.CustomClaims
		isGetRefreshToken int64
		tokenStr          string
		ok                bool
	)

	nowTime := time.Now()
	refreshClaims, ok, err = s.validateRefreshToken(req.RefreshToken)
	if !ok || err != nil {
		err = serror.NewErr().SetErr(err)
		return
	}

	// 生成token 并返回
	switch req.Type {
	case enum.ValidateUser:
		auth = jwt.NewJwtToken()
	case enum.ValidateAdmin:
		auth = jwt.NewJwtAdminToken()
	case enum.ValidateMerchant:
		auth = jwt.NewJwtMerchantToken()
	default:
		err = serror.NewErr().SetMsg("token类型错误")
		return
	}

	tokenStr, claims, err = auth.GenerateToken(refreshClaims.UserId)
	if err != nil {
		err = serror.NewErr().SetMsg("token生成出错").SetErr(err)
		return
	}

	// 判断刷新token是否需要刷新
	if (refreshClaims.ExpiresAt-nowTime.Unix())/3600 <= config.Config.App.TokenExpireTime {
		isGetRefreshToken = 1
	}

	token = &responses.GetTokenResp{
		Token:             tokenStr,
		RefreshToken:      req.RefreshToken,
		ExpireAt:          claims.ExpiresAt,
		IsGetRefreshToken: isGetRefreshToken,
	}

	return
}

func (s *defaultLoginSrv) GetRefreshToken(req *requests.GetTokenReq) (token *responses.GetRefreshTokenResp, err error) {
	var (
		claims        *jwt.CustomClaims
		refreshClaims *jwt.CustomClaims
		tokenStr      string
		ok            bool
	)

	refreshClaims, ok, err = s.validateRefreshToken(req.RefreshToken)
	if !ok || err != nil {
		err = serror.NewErr().SetErr(err)
		return
	}

	tokenStr, claims, err = jwt.NewJwtRefreshToken().GenerateToken(refreshClaims.UserId, req.Type)
	if err != nil {
		err = serror.NewErr().SetMsg("token生成出错").SetErr(err)
		return
	}

	token = &responses.GetRefreshTokenResp{
		Token:    tokenStr,
		ExpireAt: claims.ExpiresAt,
	}

	return
}

func (s *defaultLoginSrv) validateRefreshToken(refreshToken string) (
	claims *jwt.CustomClaims, isValid bool, err error) {
	var (
		auth *jwt.Token
	)

	auth = jwt.NewJwtRefreshToken()
	claims, err = auth.ParseToken(refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt2.ErrTokenExpired):
			err = serror.NewErr().SetMsg("授权已过期").SetErr(err)
			return
		case errors.Is(err, jwt2.ErrTokenNotValidYet),
			errors.Is(err, jwt2.ErrTokenMalformed),
			errors.Is(err, jwt2.ErrTokenSignatureInvalid):
			err = serror.NewErr().SetMsg("无效授权").SetErr(err)
			return
		default:
			err = serror.NewErr().SetMsg("无效授权").SetErr(err)
			return
		}
	}

	isValid = true
	return
}
