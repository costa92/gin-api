package model

type PageRequest struct {
	Page     int64 `json:"page" query:"page" form:"page"`
	PageSize int64 `json:"page_size" query:"page_size" form:"page_size"`
}

func (p *PageRequest) GetOffsetPageRequest() int64 {
	if p == nil {
		return 0
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	switch {
	case p.PageSize > 500:
		p.PageSize = 500
	case p.PageSize <= 0:
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}
