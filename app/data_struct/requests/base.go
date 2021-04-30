package requests

type Base struct {
	OrderField string `json:"order_field"` // 排序字段
	OrderType  uint8  `json:"order_type"`  // 排序类型 0升序 1降序
	Page       int    `json:"page"`        // 当前页
	PageSize   int    `json:"size"`        // 每页条数
}
