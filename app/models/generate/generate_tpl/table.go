/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 13/08/2021
 * @Desc: 模板生成器 table 对应的结构体
 */

package generate_tpl

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}
