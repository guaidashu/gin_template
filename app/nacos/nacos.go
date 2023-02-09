package nacos

import "github.com/nacos-group/nacos-sdk-go/v2/model"

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
func GetConfig(dataId, group string, data interface{}, isAutoListen ...bool) error {
	return GetConfigByNamespace(dataId, group, DefaultNamespace, data, isAutoListen...)
}

// 一体化获取配置和监听等
// 此方法传入一个 变量指针(data) 具体用法见 nacos_test.go TestGetPropertiesConfig方法
//
// dataId 相当于 etcd的key
// group 是分组ID，以服务块划分，可自定义
// isAutoListen 为是否自动监听，默认是会监听，传入false则不监听
//
// 调用此方法后会自动发起一个监听协程
// 配置发生改变时，会自动更新传入指针地址对应的值，所以不需要另外启动监听程序
// 除非有其他监听操作，否则不需要调用 ListenConfig 方法
//
// data 必须传入指针，且如果是结构体则成员必须加上 properties tag，具体用法见 nacos_test.go TestGetPropertiesConfig方法
func GetPropertiesConfig(dataId, group string, data interface{}, isAutoListen ...bool) error {
	return GetPropertiesConfigByNamespace(dataId, group, DefaultNamespace, data, isAutoListen...)
}

// 发布配置
// dataId 相当于 etcd的key
// group 是分组ID，以服务块划分，可自定义
// content 为要发布的配置数据，传入一个指针变量(且不需要在外序列化) 必须为指针
func PublishConfig(dataId, group string, content interface{}) error {
	return PublishConfigByNamespace(dataId, group, DefaultNamespace, content)
}

// 监听配置，特殊情况进行调用
// 返回源数据可供其他操作，但会自动改变 传入的data变量底层值
// PS: 这里返回的callbackData是源数据
func ListenConfig(dataId, group string, fn func(eventType EventType, data string)) error {
	return ListenConfigByNamespace(dataId, group, DefaultNamespace, fn)
}

// 取消监听
func CancelListen(dataId, group string) error {
	return CancelListenByNamespace(dataId, group, DefaultNamespace)
}

// 搜索配置
// search  require search=accurate--精确搜索  search=blur--模糊搜索
// example:
// data, err := SearchConfig("blur", "*language_*", "i18n", 1, 10)
func SearchConfig(searchType, dataId, group string, pageNo, pageSize int) (
	data *model.ConfigPage, err error) {
	data, err = SearchConfigByNamespace(searchType, dataId, group, DefaultNamespace, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 关闭客户端
func Close() {
	CloseByNamespace(DefaultNamespace)
}
