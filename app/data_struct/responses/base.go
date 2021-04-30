package responses

type BasePageList struct {
	Page      int         `json:"page"`       // 当前页
	PageSize  int         `json:"size"`       // 每页条数
	Total     int64       `json:"total"`      // 总条数
	TotalPage int         `json:"total_page"` // 总页数
	List      interface{} `json:"list"`       // 数据列表
}
