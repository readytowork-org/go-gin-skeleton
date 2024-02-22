package user

import "boilerplate-api/common/helpers"

type UserPagination struct {
	helpers.Pagination
	Keyword2 string `form:"keyword2"`
}
