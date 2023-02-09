package nacos_test

import (
	"fmt"
	"gin_template/app/nacos"
	"testing"
)

func TestInitNacosConfig(t *testing.T) {
	configHandler, err := nacos.NewNacosConfigHandler(&nacos.Params{
		Endpoints: []string{"127.0.0.1:8848"},
		Namespace: "851bb9ed-322a-4e8e-8dd3-4e9387d30cd8",
		Username:  "nacos",
		Password:  "nacos",
		TimeoutMs: 5000,
	})
	if err != nil {
		fmt.Println("初始化nacos客户端失败, err: ", err)
		return
	}

	ok, err := configHandler.PublishConfig("test-1", "test", "test-data-3", "")
	if err != nil {
		fmt.Println("发布失败, err: ", err)
		return
	}
	if !ok {
		fmt.Println("发布失败")
		return
	}
	fmt.Println(ok)

}
