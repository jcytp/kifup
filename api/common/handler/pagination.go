// common/handler/pagination.go

package handler

import "math"

type PaginationRequest struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	MaxPage    int `json:"max_page"`
}

func NewPaginatedResponse(req *PaginationRequest, totalCount int) *PaginatedResponse {
	maxPage := int(math.Max(1, math.Ceil(float64(totalCount)/float64(req.Limit))))

	return &PaginatedResponse{
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
		MaxPage:    maxPage,
	}
}
