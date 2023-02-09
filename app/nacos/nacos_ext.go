package nacos

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gin_template/app/libs"
	jsoniter "github.com/json-iterator/go"
	"github.com/magiconair/properties"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/young2j/gocopy"
	"reflect"
	"sync"
)

const (
	DefaultNamespace = "default"
)

type (
	Config struct {
		Data     interface{} `json:"data"`
		DataType string      `json:"data_type"` // 数据类型，用于作数据类型校验
	}

	// 存储实例，自动监听的map对应数据，目前无实际意义，主要用于map判断是否存在对应数据
	// 因为自动监听主需要启动一次监听，所以当map内已经存在对应这个实例，表示不需要再启动监听
	NacosConfigStore struct{}

	EventType int8
)

const (
	PUT    EventType = 0
	DELETE EventType = 1
)

var (
	// 配置监听map，用于判断是否加入了监听 map[group]map[dataId]interface{}
	NacosConfigMap map[string]map[string]*NacosConfigStore
	NacosClient    map[string]*ConfigHandler
	lock           sync.Mutex
)

// 初始化 Nacos配置中心
// 环境变量为 NACOS_CONFIG
// 这个为必须的配置
// NACOS_CONFIG = {"default": {"endpoints":["127.0.0.1:8848"],"server_address":"127.0.0.1","server_port":8848,"namespace":"控制台生成的Namespace Id","username":"nacos分配的用户名","password":"nacos分配的用户对应的密码"}}
// 这个为非必须的配置
// NACOS_CONFIG = {"default": {"endpoints":["127.0.0.1:8848"],"server_address":"127.0.0.1","server_port":8848,"namespace":"控制台生成的Namespace Id","username":"nacos分配的用户名","password":"nacos分配的用户对应的密码","log_dir":"日志目录","cache_dir":"缓存目录","timeout_ms":"请求超时时间，整数(ms)"}}
// 本地启动可以添加如下配置，可略作更改(nacos默认账户名密码是 nacos nacos)：
// {"default": {"endpoints":["127.0.0.1:8848"],"server_address":"127.0.0.1","server_port":8848,"namespace":"自行创建添加","username":"nacos","password":"nacos"}}
func InitNacosByConfig(config string) error {
	NacosConfigMap = make(map[string]map[string]*NacosConfigStore)
	NacosClient = make(map[string]*ConfigHandler)

	// Params具体参数示例见 common/pkg/nacos/nacos_test.go文件
	nacosConfigs := make(map[string]*Params)
	err := json.Unmarshal(libs.String2Bytes(config), &nacosConfigs)
	if err != nil {
		fmt.Println("json.Unmarshal nacos链接参数失败")
		return err
	}
	for k, v := range nacosConfigs {
		NacosClient[k], err = NewNacosConfigHandler(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// 一体化获取配置和监听等
// 此方法传入一个 变量指针(data) 具体用法见 nacos_test.go TestGetConfig方法
//
// dataId 相当于 etcd的key
// group 是分组ID，以服务块划分，可自定义
// isAutoListen 为是否自动监听，默认是会监听，传入false则不监听
//
// 调用此方法后会自动发起一个监听协程
// 配置发生改变时，会自动更新传入指针地址对应的值，所以不需要另外启动监听程序
// 除非有其他监听操作，否则不需要调用 ListenConfig 方法
//
// data 必须传入指针，且如果是结构体则成员必须加上 json tag，具体用法见 nacos_test.go TestGetConfig方法
func GetConfigByNamespace(dataId, group, namespace string, data interface{}, isAutoListen ...bool) error {
	isListen := false

	lock.Lock()
	defer func() {
		_ = recover()
		lock.Unlock()
	}()
	if _, ok := NacosConfigMap[group]; !ok {
		NacosConfigMap[group] = make(map[string]*NacosConfigStore)
	}

	// 如果NacosConfigMap[group][dataId]对应数据 不存在,则进入监听流程，存在则直接进行获取
	if _, ok := NacosConfigMap[group][dataId]; !ok {
		NacosConfigMap[group][dataId] = &NacosConfigStore{}

		if len(isAutoListen) > 0 {
			isListen = isAutoListen[0]
		}
		// 发起监听
		if isListen {
			libs.RunSafe(func() {
				listenConfigByNamespace(dataId, group, namespace, data)
			})
		}
	}

	content, err := NacosClient[namespace].GetConfig(dataId, group)
	if err != nil {
		libs.Logger.Error("Nacos配置数据获取失败")
		return err
	}

	err = jsoniter.Unmarshal(libs.String2Bytes(content), data)
	if err != nil {
		libs.Logger.Error("Nacos配置数据解析失败")
	}

	return err
}

// 发布配置
// dataId 相当于 etcd的key
// group 是分组ID，以服务块划分，可自定义
// content 为要发布的配置数据，传入一个指针变量(且不需要在外序列化) 必须为指针
func PublishConfigByNamespace(dataId, group, namespace string, content interface{}) error {
	if content == nil {
		return errors.New("数据不能为空")
	}

	result, err := jsoniter.MarshalIndent(content, "", "    ")
	if err != nil {
		libs.Logger.Error("序列化发布配置数据失败")
		return err
	}

	ok, err := NacosClient[namespace].PublishConfig(dataId, group, libs.Bytes2String(result), JSONType)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("publish失败")
	}

	return nil
}

// 监听配置，特殊情况进行调用
// 返回源数据可供其他操作，但会自动改变 传入的data变量底层值
// PS: 这里返回的callbackData是源数据
func ListenConfigByNamespace(dataId, group, namespace string, fn func(eventType EventType, callbackData string)) error {
	return NacosClient[namespace].ListenConfig(dataId, group, func(namespace, groupId, dataId, data string) {
		if !(groupId == group && dataId == dataId) {
			libs.Logger.Error("groupId或者dataId不匹配")
			return
		}

		if data == "" {
			fn(DELETE, data)
		} else {
			fn(PUT, data)
		}

	})
}

func CancelListenByNamespace(dataId, group, namespace string) error {
	err := NacosClient[namespace].CancelListen(dataId, group)
	if err != nil {
		libs.Logger.Error("Nacos取消监听失败")
	}

	return err
}

// 关闭客户端
func CloseByNamespace(namespace string) {
	NacosClient[namespace].Close()
}

// 获取配置数据的类型
func getContentType(content interface{}) string {
	vt := reflect.TypeOf(content)
	return vt.String()
}

func listenConfigByNamespace(dataId, group, namespace string, data interface{}) {
	defer func() {
		rcv := recover()
		if rcv != nil {
			libs.Logger.Error(fmt.Sprintf("自动监听配置panic, recover error: %v", rcv))
		}
	}()

	err := NacosClient[namespace].ListenConfig(dataId, group, func(namespace, callbackGroupId, callbackDataId, callbackData string) {
		if !(callbackGroupId == group && callbackDataId == dataId) {
			return
		}

		e := jsoniter.Unmarshal([]byte(callbackData), data)
		if e != nil {
			libs.Logger.Error("Nacos配置监听数据解析出错")
		}
	})
	if err != nil {
		libs.Logger.Error("Nacos监听配置失败")
		return
	}
}

func GetPropertiesConfigByNamespace(dataId, group, namespace string, data interface{}, isAutoListen ...bool) error {
	isListen := false

	lock.Lock()
	defer func() {
		_ = recover()
		lock.Unlock()
	}()
	if _, ok := NacosConfigMap[group]; !ok {
		NacosConfigMap[group] = make(map[string]*NacosConfigStore)
	}

	if _, ok := NacosConfigMap[group][dataId]; !ok {
		NacosConfigMap[group][dataId] = &NacosConfigStore{}

		if len(isAutoListen) > 0 {
			isListen = isAutoListen[0]
		}
		// 发起监听
		if isListen {
			libs.RunSafe(func() {
				listenPropertiesConfigByNamespace(dataId, group, namespace, data)
			})
		}
	}

	content, err := NacosClient[namespace].GetConfig(dataId, group)
	if err != nil {
		libs.Logger.Error("Nacos配置数据获取失败")
		return err
	}

	// or from a string
	p := properties.MustLoadString(content)
	m := p.Map()
	gocopy.Copy(data, &m)

	return err
}

func listenPropertiesConfigByNamespace(dataId, group, namespace string, data interface{}) {
	defer func() {
		rcv := recover()
		if rcv != nil {
			libs.Logger.Error(fmt.Sprintf("自动监听配置panic, recover error: %v", rcv))
		}
	}()

	err := NacosClient[namespace].ListenConfig(dataId, group, func(namespace, callbackGroupId, callbackDataId, callbackData string) {
		if !(callbackGroupId == group && callbackDataId == dataId) {
			return
		}
		// or from a string
		p := properties.MustLoadString(callbackData)
		m := p.Map()
		gocopy.Copy(data, &m)
	})
	if err != nil {
		libs.Logger.Error("Nacos监听配置失败")
		return
	}
}

func GetYamlConfigByNamespace(dataId, group, namespace string, data interface{}, isAutoListen ...bool) error {
	isListen := false

	lock.Lock()
	defer func() {
		_ = recover()
		lock.Unlock()
	}()
	if _, ok := NacosConfigMap[group]; !ok {
		NacosConfigMap[group] = make(map[string]*NacosConfigStore)
	}

	if _, ok := NacosConfigMap[group][dataId]; !ok {
		NacosConfigMap[group][dataId] = &NacosConfigStore{}

		if len(isAutoListen) > 0 {
			isListen = isAutoListen[0]
		}
		// 发起监听
		if isListen {
			libs.RunSafe(func() {
				listenYamlConfigByNamespace(dataId, group, namespace, data)
			})
		}
	}

	content, err := NacosClient[namespace].GetConfig(dataId, group)
	if err != nil {
		libs.Logger.Error("Nacos配置数据获取失败")
		return err
	}

	err = xml.Unmarshal([]byte(content), data)

	return err
}


func listenYamlConfigByNamespace(dataId, group, namespace string, data interface{}) {
	defer func() {
		rcv := recover()
		if rcv != nil {
			libs.Logger.Error(fmt.Sprintf("自动监听配置panic, recover error: %v", rcv))
		}
	}()

	err := NacosClient[namespace].ListenConfig(dataId, group, func(namespace, callbackGroupId, callbackDataId, callbackData string) {
		if !(callbackGroupId == group && callbackDataId == dataId) {
			return
		}

		e := xml.Unmarshal([]byte(callbackData), data)
		if e != nil {
			libs.Logger.Error("Nacos监听配置失败")
			return
		}
	})
	if err != nil {
		libs.Logger.Error("Nacos监听配置失败")
		return
	}
}

// 搜索配置
// search  require search=accurate--精确搜索  search=blur--模糊搜索
// example:
// data, err := SearchConfigByNamespace("blur", "*language_*", "i18n", namespace, 1, 10)
func SearchConfigByNamespace(searchType, dataId, group, namespace string, pageNo, pageSize int) (
	data *model.ConfigPage, err error) {
	data, err = NacosClient[namespace].SearchConfig(searchType, dataId, group, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	return data, nil
}
