package _interface

import "gin_template/app/enum"

type (
	ServiceParam struct{}
)

// 第三方组件初始化接口定义
type ComponentsInit interface {
	Init(params *ServiceParam) error    // 初始化方法
	ComponentName() enum.BootModuleType // 组件名
	Close() error                       // 关闭服务
}
