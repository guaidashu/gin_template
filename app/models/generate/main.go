package main

import (
	"flag"
	"fmt"
	"gin_template/app/models/generate/generate_tpl"
	"os"
)

// 根据模板文件自动生成 model文件
func main() {
	var (
		flagArr []string
	)

	flag.Parse()
	flagArr = flag.Args()

	if len(flagArr) >= 1 {
		// 进入选择
		switch flagArr[0] {
		case "new":
			// 进入 创建环节
			// 首先判断是否有文件夹名字
			if len(flagArr) > 1 {
				generate_tpl.Generate(flagArr[1])
			} else {
				fmt.Println("Please input a director name.")
				os.Exit(1)
			}
		default:
			fmt.Println("Please input a correct command")
		}
	}
}
