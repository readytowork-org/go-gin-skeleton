package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	service services.BlogService
	logger  infrastructure.Logger
}

func NewBlogController(
	service services.BlogService,
	logger infrastructure.Logger,
) BlogController {
	return BlogController{
		service: service,
		logger:  logger,
	}
}

func (cc BlogController) GetAllBlogs(c *gin.Context) {
	pagination := utils.BuildPagination(c)

	blogs, count, err := cc.service.GetAllBlogs(pagination)
	if err != nil {
		cc.logger.Zap.Error("Error Finding blogs: ", err.Error())
		err := errors.InternalError.Wrap(err, "failed to get blogs")
		responses.HandleError(c, err)
		return
	}
	responses.JSONCount(c, http.StatusOK, blogs, count)
}
