package entity

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int `json:"page" form:"page"`         // 页码，从1开始
	PageSize int `json:"page_size" form:"page_size"` // 每页大小
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页大小
	Total      int64       `json:"total"`       // 总记录数
	TotalPages int         `json:"total_pages"` // 总页数
	Data       interface{} `json:"data"`        // 数据列表
}

// NewPaginationRequest 创建分页请求，设置默认值
func NewPaginationRequest(page, pageSize int) *PaginationRequest {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}
	
	return &PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(page, pageSize int, total int64, data interface{}) *PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}
}