package dto

// PageQuery holds common pagination query parameters shared by all list endpoints.
type PageQuery struct {
	Page     int `form:"page"      binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

func (q *PageQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

func (q *PageQuery) Offset() int {
	q.Normalize()
	return (q.Page - 1) * q.PageSize
}
