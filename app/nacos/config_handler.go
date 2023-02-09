package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strings"
)

type (
	// 这里的结构体只供查看，不做任何用途

	// constant.ClientConfig
	ClientConfig struct {
		TimeoutMs   uint64 // timeout for requesting Nacos server, default value is 10000ms
		NamespaceId string // the namespaceId of Nacos
		Endpoint    string // the endpoint for ACM. https://help.aliyun.com/document_detail/130146.html
		RegionId    string // the regionId for ACM & KMS
		AccessKey   string // the AccessKey for ACM & KMS
		SecretKey   string // the SecretKey for ACM & KMS
		OpenKMS     bool   // it's to open KMS, default is false. https://help.aliyun.com/product/28933.html
		// , to enable encrypt/decrypt, DataId should be start with "cipher-"
		CacheDir             string // the directory for persist nacos service info,default value is current path
		UpdateThreadNum      int    // the number of goroutine for update nacos service info,default value is 20
		NotLoadCacheAtStart  bool   // not to load persistent nacos service info in CacheDir at start time
		UpdateCacheWhenEmpty bool   // update cache when get empty service instance from server
		Username             string // the username for nacos auth
		Password             string // the password for nacos auth
		LogDir               string // the directory for log, default is current path
		RotateTime           string // the rotate time for log, eg: 30m, 1h, 24h, default is 24h
		MaxAge               int64  // the max age of a log file, default value is 3
		LogLevel             string // the level of log, it's must be debug,info,warn,error, default value is info
	}

	// constant.ServerConfig
	ServerConfig struct {
		Scheme      string // the nacos server scheme,defaut=http,this is not required in 2.0
		ContextPath string // the nacos server contextpath,defaut=/nacos,this is not required in 2.0
		IpAddr      string // the nacos server address
		Port        uint64 // nacos server port
		GrpcPort    uint64 // nacos server grpc port, default=server port + 1000, this is not required
	}
)

// 初始化Nacos配置中心客户端(不包含服务发现的初始化)
func InitNacos(params *Params) (sc *Box) {
	sc = &Box{
		Server: []constant.ServerConfig{
			{
				// Scheme:      "http",
				IpAddr:      params.ServerAddress,
				Port:        params.ServerPort,
				ContextPath: "/nacos",
			},
		},
		Client: constant.ClientConfig{
			Endpoint:    strings.Join(params.Endpoints, ","),
			NamespaceId: params.Namespace,
			AccessKey:   params.AccessKey,
			SecretKey:   params.SecretKey,
			Username:    params.Username,
			Password:    params.Password,
			CacheDir:    params.CacheDir,
			LogDir:      params.LogDir,
			TimeoutMs:   params.TimeoutMs,
			// ListenInterval: 30 * 1000, 已弃用
		},
	}

	return
}

// 生成配置中心客户端
func (box *Box) CreateConfigClient() (configHandler *ConfigHandler, err error) {
	// client, err := clients.CreateConfigClient(map[string]interface{}{
	// 	// Nacos Golang SDK 2.0里面必须传入serverConfigs这个参数
	// 	"serverConfigs": box.Server,
	// 	"clientConfig":  box.Client,
	// })
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: &box.Client,
		// Nacos Golang SDK 2.0里面必须传入serverConfigs这个参数
		ServerConfigs: box.Server,
	})
	if err != nil {
		return nil, err
	}

	configHandler = &ConfigHandler{}
	configHandler.Client = client

	return configHandler, err
}

func NewNacosConfigHandler(params *Params) (*ConfigHandler, error) {
	sc := InitNacos(params)

	configHandler, err := sc.CreateConfigClient()
	if err != nil {
		fmt.Println("初始化nacos客户端失败")
		return nil, err
	}

	return configHandler, nil
}

// 查询指定dataId内容
func (c *ConfigHandler) GetConfig(dataId, group string) (content string, err error) {
	content, err = c.Client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	return
}

// 发布配置
func (c *ConfigHandler) PublishConfig(dataId, group, content, fmtType string) (published bool, err error) {
	conf := vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	}
	if fmtType != "" {
		conf.Type = fmtType
	}
	published, err = c.Client.PublishConfig(conf)

	return
}

// 删除配置
func (c *ConfigHandler) DeleteConfig(dataId, group string) (deleted bool, err error) {
	deleted, err = c.Client.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})

	return
}

// 监听配置
func (c *ConfigHandler) ListenConfig(dataId, group string,
	callback func(namespace, group, dataId, data string)) (err error) {
	err = c.Client.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: callback,
	})

	return
}

// 取消监听
func (c *ConfigHandler) CancelListen(dataId, group string) (err error) {
	err = c.Client.CancelListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})

	return
}

// 搜索配置
// search  require search=accurate--精确搜索  search=blur--模糊搜索
func (c *ConfigHandler) SearchConfig(searchType, dataId, group string, pageNo, pageSize int) (
	searchPage *model.ConfigPage, err error) {
	searchPage, err = c.Client.SearchConfig(vo.SearchConfigParam{
		Search:   searchType,
		DataId:   dataId,
		Group:    group,
		PageNo:   pageNo,
		PageSize: pageSize,
	})

	return
}

func (c *ConfigHandler) Close() {
	c.Client.CloseClient()
}
