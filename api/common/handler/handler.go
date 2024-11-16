// common/handler/handler.go

package handler

import (
	"github.com/gin-gonic/gin"
)

func Handler(f func(*gin.Context) (string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		msg, err := f(c)
		if err != nil {
			ResponseServerError(c, msg, err)
			return
		}

		ResponseOK(c, nil)
	}
}

func HandlerIn[T any](f func(*gin.Context, T) (string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var req T
		err := c.ShouldBindJSON(&req)
		if err != nil {
			ResponseBadRequest(c, "Invalid parameter", err)
			return
		}

		msg, err := f(c, req)
		if err != nil {
			ResponseServerError(c, msg, err)
			return
		}

		ResponseOK(c, nil)
	}
}

func HandlerOut[U any](f func(*gin.Context) (U, string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		data, msg, err := f(c)
		if err != nil {
			ResponseServerError(c, msg, err)
			return
		}

		ResponseOK(c, data)
	}
}

func HandlerInOut[T any, U any](f func(*gin.Context, T) (U, string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var req T
		err := c.ShouldBindJSON(&req)
		if err != nil {
			ResponseBadRequest(c, "Invalid parameter", err)
			return
		}

		data, msg, err := f(c, req)
		if err != nil {
			ResponseServerError(c, msg, err)
			return
		}

		ResponseOK(c, data)
	}
}

func HandlerPagination[U any](f func(*gin.Context, *PaginationRequest) (U, *PaginatedResponse, string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		reqPagination := &PaginationRequest{}
		err := c.ShouldBindQuery(reqPagination)
		if err != nil {
			ResponseBadRequest(c, "Invalid parameter for pagination", err)
			return
		}

		data, paginatedResponse, msg, err := f(c, reqPagination)
		if err != nil {
			ResponseServerError(c, msg, err)
			return
		}

		ResponseOKPagination(c, data, paginatedResponse)
	}
}

func HandlerInPagination[T any, U any](f func(*gin.Context, T, *PaginationRequest) (U, *PaginatedResponse, string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		reqPagination := &PaginationRequest{}
		err := c.ShouldBindQuery(reqPagination)
		if err != nil {
			ResponseBadRequest(c, "Invalid parameter for pagination", err)
			return
		}

		var req T
		err = c.ShouldBindQuery(&req)
		if err != nil {
			ResponseBadRequest(c, "Invalid parameter", err)
			return
		}

		data, paginatedResponse, msg, err := f(c, req, reqPagination)
		if err != nil {
			ResponseServerError(c, msg, err)
			return
		}

		ResponseOKPagination(c, data, paginatedResponse)
	}
}
