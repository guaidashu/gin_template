/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: 小程序相关操作
 */

package miniprogram

import (
	"fmt"
	"gin_template/app/config"
	"gin_template/app/data_struct"
	"gin_template/app/libs"

	"github.com/guaidashu/go_helper"
)

// 获取OpenId
func GetOpenId(code string) (openInfo *data_struct.MiniProGramLoginInfo, err error) {
	var (
		miniProGramLoginInfo data_struct.MiniProGramLoginInfo
	)

	// 登录需要请求的 url
	requestUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code"

	requestUrl = fmt.Sprintf(requestUrl,
		config.Config.MiniProgram.AppId,
		config.Config.MiniProgram.Secret,
		code)

	if err = go_helper.Get(requestUrl).ToJSON(&miniProGramLoginInfo); err != nil {
		err = libs.NewReportError(err)
		return
	}

	openInfo = &miniProGramLoginInfo

	return
}
