/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/10/2021
 * @Desc: user service 层
 */

package service

import (
	"fmt"
	"gin_template/app/data_struct/requests"
	"gin_template/app/data_struct/responses"
	"gin_template/app/libs"
	"gin_template/app/libs/jwt"
	"gin_template/app/libs/serror"
	"gin_template/app/models"
	"sync"
)

type (
	UserSrv interface {
		// 获取用户信息
		UserInfo(userId int64) (data *models.UserModel, err error)
		// 注册用户
		Register(req *requests.RegisterReq) (token *responses.LoginResp, err error)
	}

	defaultUserSrv struct {
	}
)

var (
	_userSrv     UserSrv
	_userSrvOnce sync.Once
)

func NewUserSrv() UserSrv {
	_userSrvOnce.Do(func() {
		_userSrv = &defaultUserSrv{}
	})

	return _userSrv
}

func (s *defaultUserSrv) UserInfo(userId int64) (data *models.UserModel, err error) {
	userInfo, queryErr := models.NewUserModel().FindOne(userId)
	if queryErr != nil {
		err = serror.NewErr().SetMsg("服务器错误").SetErr(err)
	}

	data = userInfo
	fmt.Println(userInfo, queryErr)
	return
}

func (s *defaultUserSrv) Register(req *requests.RegisterReq) (token *responses.LoginResp, err error) {
	// 创建用户
	var (
		user *models.UserModel
	)

	user = new(models.UserModel)
	err = libs.Struct2Struct(req, user)
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return
	}
	user.OpenId = req.OpenId
	user.Username = req.Nickname
	user.AvatarUrl = req.AvatarUrl

	if err = models.NewUserModel().Create(user); err != nil {
		err = serror.NewErr().SetErr(err)
	}

	tokenStr, claims, generateErr := jwt.NewJwtToken().GenerateToken(user.Id)
	if generateErr != nil {
		err = serror.NewErr().SetErr(generateErr)
		return
	}
	refreshToken, _, generateErr := jwt.NewJwtRefreshToken().GenerateToken(user.Id)
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
