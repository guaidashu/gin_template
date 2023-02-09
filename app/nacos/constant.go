package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

type (
	Box struct {
		Server []constant.ServerConfig
		Client constant.ClientConfig
	}

	Params struct {
		Endpoints     []string `json:"endpoints"`      // Nacos地址
		ServerAddress string   `json:"server_address"` // 服务端地址
		ServerPort    uint64   `json:"server_port"`    // 服务端端口
		Namespace     string   `json:"namespace"`      // 命名空间 id
		AccessKey     string   `json:"access_key"`     // 当要上阿里云时，阿里云上面的一个云账号名
		SecretKey     string   `json:"secret_key"`     // 当要上阿里云时，阿里云上面的一个云账号密码
		Username      string   `json:"username"`       // 用户名
		Password      string   `json:"password"`       // 用户密码
		LogDir        string   `json:"log_dir"`        // 日志目录
		CacheDir      string   `json:"cache_dir"`      // 缓存目录
		TimeoutMs     uint64   `json:"timeout_ms"`     // 请求超时时间 Nacos 默认为10000ms
	}
)

const (
	JSONType = "json"
)

type (
	ConfigHandler struct {
		Client config_client.IConfigClient
	}
)
