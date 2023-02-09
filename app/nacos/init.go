package nacos

import (
	"errors"
	"fmt"
	"gin_template/app/config"
	"gin_template/app/data_struct/_interface"
	"gin_template/app/enum"
	"os"
)

type (
	NacosInit struct{}
)

func NewNacosInit() _interface.ComponentsInit {
	return &NacosInit{}
}

func (n *NacosInit) Init(*_interface.ServiceParam) error {
	return InitNacosConfig()
}

func (n *NacosInit) ComponentName() enum.BootModuleType {
	return enum.NacosInit
}

func (n *NacosInit) Close() error {
	Close()
	return nil
}


// 初始化 Nacos配置中心
// 环境变量为 NACOS_CONFIG
// 这个为必须的配置
// NACOS_CONFIG = {"default": {"endpoints":["127.0.0.1:8848"],"server_address":"127.0.0.1","server_port":8848,"namespace":"控制台生成的Namespace Id","username":"nacos分配的用户名","password":"nacos分配的用户对应的密码"}}
// 这个为非必须的配置
// NACOS_CONFIG = {"default": {"endpoints":["127.0.0.1:8848"],"server_address":"127.0.0.1","server_port":8848,"namespace":"控制台生成的Namespace Id","username":"nacos分配的用户名","password":"nacos分配的用户对应的密码","log_dir":"日志目录","cache_dir":"缓存目录","timeout_ms":"请求超时时间，整数(ms)"}}
// 本地启动可以添加如下配置，可略作更改(nacos默认账户名密码是 nacos nacos)：
// {"default": {"endpoints":["127.0.0.1:8848"],"server_address":"127.0.0.1","server_port":8848,"namespace":"自行创建添加","username":"nacos","password":"nacos"}}
func InitNacosConfig() error {
	// 改为从环境变量获取参数 环境变量数据格式为json字符串
	nacosConfigStr := os.Getenv("NACOS_CONFIG")
	if nacosConfigStr == "" {
		return errors.New("NACOS_CONFIG 环境变量配置未设置")
	}

	err := InitNacosByConfig(nacosConfigStr)
	if err != nil {
		return err
	}

	return nil
}

type (
	ConfigInit struct{}
)

func NewConfigInit() *ConfigInit {
	return &ConfigInit{}
}

func (c ConfigInit) Init(params *_interface.ServiceParam) error {
	return InitFromNacos()
}

func (c ConfigInit) ComponentName() enum.BootModuleType {
	return enum.ConfigInit
}

func (c ConfigInit) Close() error {
	return nil
}

// 从nacos读取配置
func InitFromNacos() error {
	err := GetYamlConfig("config", "env", &config.Config, false)
	if err != nil {
		return err
	}
	fmt.Println(config.Config)

	return nil
}
