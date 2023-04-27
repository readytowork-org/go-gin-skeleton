package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination -> struct for Pagination
type Pagination struct {
	Sort     string `form:"sort"`
	Keyword  string `form:"keyword"`
	Offset   int    `form:"page"`
	All      bool
	PageSize int
}

// BuildPagination -> builds the pagination
func BuildPagination(c *gin.Context) (pagination Pagination) {
	_ = c.BindQuery(&pagination)

	pagination.PageSize = 10
	pageSizeStr := c.Query("pageSize")
	if pageSizeStr == "Infinity" {
		pagination.All = true
	}

	if !pagination.All && pageSizeStr != "" {
		pagination.PageSize, _ = strconv.Atoi(pageSizeStr)
	}

	if pagination.Offset <= 0 {
		pagination.Offset = 1
	}
	pagination.Offset = (pagination.Offset - 1) * pagination.PageSize

	return pagination
}
