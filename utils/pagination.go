package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// // Pagination struct for Pagination
type Pagination struct {
	Sort     string `form:"sort"`
	Keyword  string `form:"keyword"`
	Offset   int    `form:"page"`
	All      bool
	PageSize int
}

// BuildPagination builds the pagination
func (m *Pagination) BuildPagination(c *gin.Context) {
	_ = c.BindQuery(&m)

	m.PageSize = 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr == "Infinity" {
		m.All = true
	} else if pageSizeStr != "" {
		m.PageSize, _ = strconv.Atoi(pageSizeStr)
	}

	if m.Offset <= 0 {
		m.Offset = 1
	}
	m.Offset = (m.Offset - 1) * m.PageSize
}
