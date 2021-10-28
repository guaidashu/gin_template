/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/09/2021
 * @Desc: 小程序支付所需工具
 */

package miniprogram

import (
	"strings"

	"github.com/shopspring/decimal"
)

func TruncatedText(data string, length int) string {
	data = FilterTheSpecialSymbol(data)
	if len([]rune(data)) > length {
		return string([]rune(data)[:length-1])
	}
	return data
}

// 过滤特殊符号
func FilterTheSpecialSymbol(data string) string {
	// 定义转换规则
	specialSymbol := func(r rune) rune {
		if r == '`' || r == '[' || r == '~' || r == '!' || r == '@' || r == '#' || r == '$' ||
			r == '^' || r == '&' || r == '*' || r == '~' || r == '(' || r == ')' || r == '=' ||
			r == '~' || r == '|' || r == '{' || r == '}' || r == '~' || r == ':' || r == ';' ||
			r == '\'' || r == ',' || r == '\\' || r == '[' || r == ']' || r == '.' || r == '<' ||
			r == '>' || r == '/' || r == '?' || r == '~' || r == '！' || r == '@' || r == '#' ||
			r == '￥' || r == '…' || r == '&' || r == '*' || r == '（' || r == '）' || r == '—' ||
			r == '|' || r == '{' || r == '}' || r == '【' || r == '】' || r == '‘' || r == '；' ||
			r == '：' || r == '”' || r == '“' || r == '\'' || r == '"' || r == '。' || r == '，' ||
			r == '、' || r == '？' || r == '%' || r == '+' || r == '_' || r == ']' || r == '"' || r == '&' {
			return ' '
		}
		return r
	}
	data = strings.Map(specialSymbol, data)
	return strings.Replace(data, "\n", " ", -1)
}

// 微信金额浮点转字符串
func WechatMoneyFeeToString(moneyFee float64) string {
	aDecimal := decimal.NewFromFloat(moneyFee)
	bDecimal := decimal.NewFromFloat(100)
	return aDecimal.Mul(bDecimal).Truncate(0).String()
}
