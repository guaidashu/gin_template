/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/10/2021
 * @Desc: 用户相关
 */

package requests

type (
	UserInfoReq struct {
		UserId int64
	}

	RegisterReq struct {
		OpenId     string `json:"openid"`      // 小程序 open_id
		SessionKey string `json:"session_key"` // 小程序 session_key
		Token      string `json:"token"`       // 登录token，备用
		AvatarUrl  string `json:"avatarUrl"`   // 用户头像
		Country    string `json:"country"`     // 国家
		Province   string `json:"province"`    // 省份
		City       string `json:"city"`        // 城市
		Nickname   string `json:"nickName"`    // 用户名
		Language   string `json:"language"`    // 语言
		Sex        uint64 `json:"gender"`      // 性别
		Email      string `json:"email"`       // 邮箱
	}
)
