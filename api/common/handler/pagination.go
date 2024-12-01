// common/handler/pagination.go

package handler

import "math"

type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	MaxPage    int `json:"max_page"`
}

func (req *PaginationRequest) LimitOffset() (int, int) {
	return req.PageSize, req.PageSize * (req.Page - 1)
}

func (req *PaginationRequest) NewPaginatedResponse(totalCount int) *PaginatedResponse {
	maxPage := int(math.Max(1, math.Ceil(float64(totalCount)/float64(req.PageSize))))

	return &PaginatedResponse{
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		MaxPage:    maxPage,
	}
}
