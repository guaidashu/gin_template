/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package enum

const (
	AutoWhere        = "and"     // 自动 and 拼接
	AutoOr           = "or"      // 自动 or 拼接
	AutoCustomHandle = "custom"  // 自定义函数 处理
	AutoLike         = "like"    // 自动 like 模糊匹配
	FieldIsHandle    = "sql"     // 操作类型
	FieldHandle      = "sqlCond" // 操作具体语句
)
