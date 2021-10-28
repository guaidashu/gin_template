/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/10/2021
 * @Desc: user 相关返回数据结构
 */

package responses

type (
	LoginResp struct {
		Token        string `json:"token"`         // token
		RefreshToken string `json:"refresh_token"` // 刷新token 用于重新获取token和刷新token
		ExpireAt     int64  `json:"expire_at"`     // 过期时间
		Status       int64  `json:"status"`        // 0 正常登录或其他正常流程 1未注册
		OpenId       string `json:"open_id"`       // 用户唯一ID
	}

	GetTokenResp struct {
		Token             string `json:"token"`                // token
		RefreshToken      string `json:"refresh_token"`        // 刷新token
		ExpireAt          int64  `json:"expire_at"`            // 过期时间
		IsGetRefreshToken int64  `json:"is_get_refresh_token"` // 是否重新获取刷新token 1需要重新获取 0 不需要重新获取
	}

	GetRefreshTokenResp struct {
		Token    string `json:"token"`     // token
		ExpireAt int64  `json:"expire_at"` // 过期时间
	}
)
