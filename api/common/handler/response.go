// common/handler/response.go

package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResponseOKPagination(c *gin.Context, data any, paginatedResponse *PaginatedResponse) {
	if data == nil {
		c.JSON(http.StatusOK, gin.H{"result": "OK"})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "OK", "data": data, "pagination": paginatedResponse})
	}
}

func ResponseOK(c *gin.Context, data any) {
	if data == nil {
		c.JSON(http.StatusOK, gin.H{"result": "OK"})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "OK", "data": data})
	}
}

func ResponseError(c *gin.Context, status int, msg string, err error) {
	if strings.HasPrefix(msg, "UNAUTHORIZED") {
		status = http.StatusUnauthorized
	}

	if err != nil {
		if msg != "" {
			slog.WarnContext(c, err.Error())
		} else {
			msg = err.Error()
		}
	}
	slog.WarnContext(c, msg)

	c.JSON(status, gin.H{"result": "Error", "data": msg})
}

func ResponseNotFound(c *gin.Context, msg string, err error) {
	ResponseError(c, http.StatusNotFound, msg, err)
}

func ResponseUnauthorized(c *gin.Context, msg string, err error) {
	ResponseError(c, http.StatusUnauthorized, msg, err)
}

func ResponseBadRequest(c *gin.Context, msg string, err error) {
	ResponseError(c, http.StatusBadRequest, msg, err)
}

func ResponseServerError(c *gin.Context, msg string, err error) {
	ResponseError(c, http.StatusInternalServerError, msg, err)
}
