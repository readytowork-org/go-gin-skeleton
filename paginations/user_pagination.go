package paginations

type UserPagination struct {
	Pagination
	Keyword2 string `form:"keyword2"`
}
