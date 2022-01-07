package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

//Pagination -> struct for Pagination
type Pagination struct {
	Page       int
	Sort       string
	PageSize   int
	Offset     int
	All        bool
	Keyword    string
	Category   string
	IsPublic   string
	Area       string
	CategoryId int
}

//BuildPagination -> builds the pagination
func BuildPagination(c *gin.Context) Pagination {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")
	sort := c.Query("sort")
	keyword := c.Query("keyword")
	category := c.Query("category")
	categoryId := c.Query("category_id")
	area := c.Query("area")
	isPublic := c.Query("is_public")

	var all bool
	if pageSizeStr == "Infinity" {
		all = true
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	cid, err := strconv.Atoi(categoryId)
	if err != nil || cid <= 0 {
		cid = 0
	}

	return Pagination{
		Page:       page,
		Sort:       sort,
		PageSize:   pageSize,
		Offset:     (page - 1) * pageSize,
		All:        all,
		Keyword:    keyword,
		Category:   category,
		IsPublic:   isPublic,
		Area:       area,
		CategoryId: cid,
	}
}

//FacilityPagination -> struct for FacilityPagination
type FacilityPagination struct {
	Keyword      string
	IsSubjected  string
	AcquiredOnly string
	UserId       int64
	Area         string
	PageSize     int
}

//BuildFacilityPagination -> builds the FacilityPagination
func BuildFacilityPagination(c *gin.Context) FacilityPagination {
	keyword := c.Query("keyword")
	isSubjected := c.Query("is_subjected")
	acquiredOnly := c.Query("acquired")
	area := c.Query("area")
	pageSizeStr := c.Query("pageSize")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	return FacilityPagination{
		Keyword:      keyword,
		IsSubjected:  isSubjected,
		AcquiredOnly: acquiredOnly,
		Area:         area,
		PageSize:     pageSize,
	}
}
