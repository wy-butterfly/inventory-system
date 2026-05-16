package request

// PageQuery 分页查询通用参数
type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

func (q *PageQuery) GetPage() int {
	if q.Page <= 0 {
		return 1
	}
	return q.Page
}

func (q *PageQuery) GetPageSize() int {
	if q.PageSize <= 0 {
		return 20
	}
	return q.PageSize
}

func (q *PageQuery) GetOffset() int {
	return (q.GetPage() - 1) * q.GetPageSize()
}
