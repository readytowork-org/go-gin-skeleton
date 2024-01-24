package url_query

type UserPagination struct {
	Pagination
	Keyword2 string `form:"keyword2"`
}
