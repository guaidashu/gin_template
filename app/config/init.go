/**
  create by yy on 2019-07-29
*/

package config

import (
	"fmt"
	"gin_template/app/data_struct"
	"github.com/guaidashu/go_helper/configor"
)

var Config data_struct.Config

func init() {
	fmt.Println("开始加载开发配置文件")
	err := configor.Load(&Config, "app/config/config_dev.yml")
	if err != nil {
		fmt.Println("开发环境配置文件加载失败")
		err = configor.Load(&Config, "app/config/config_product.yml")
		if err != nil {
			fmt.Println("线上环境配置文件加载失败")
			fmt.Println("配置文件加载失败")
		}
	} else {
		fmt.Println("配置文件加载完成")
	}
}
