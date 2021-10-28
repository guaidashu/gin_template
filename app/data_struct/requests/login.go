/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/10/2021
 * @Desc: 登录相关
 */

package requests

type (
	MiniProgramLoginReq struct {
		Code string `json:"code"` // 小程序登录返回所需的code
	}

	AdminLoginReq struct {
		PhoneNumber string `json:"phone_number"` // 电话号码
		Password    string `json:"password"`     // 密码
	}

	GetTokenReq struct {
		RefreshToken string `json:"refresh_token"` // 刷新token
		Type         int64  `json:"type"`          // 请求token 的类型, 1 普通用户 2 商家 3 管理员
	}
)
