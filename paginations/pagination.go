package paginations

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type IPagination interface {
	Build(c *gin.Context)
}

// BuildPagination -> binds the query and builds pagination
func BuildPagination[T IPagination](c *gin.Context) (m T) {
	_ = c.BindQuery(&m)

	m.Build(c)
	return m
}

// Pagination struct for Pagination
type Pagination struct {
	Sort     string `form:"sort"`
	Keyword  string `form:"keyword"`
	Offset   int    `form:"page"`
	All      bool
	PageSize int
}

// Build builds the pagination
func (m *Pagination) Build(c *gin.Context) {
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
