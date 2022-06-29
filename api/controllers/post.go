package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostController -> struct
type PostController struct {
	logger      infrastructure.Logger
	PostService services.PostService
}

// NewPostController -> constructor
func NewPostController(
	logger infrastructure.Logger,
	PostService services.PostService,
) PostController {
	return PostController{
		logger:      logger,
		PostService: PostService,
	}
}

// CreatePost -> Create Post
func (cc PostController) CreatePost(c *gin.Context) {
	post := models.Post{}

	if err := c.ShouldBindJSON(&post); err != nil {
		cc.logger.Zap.Error("Error [CreatePost] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := cc.PostService.CreatePost(post); err != nil {
		cc.logger.Zap.Error("Error [CreatePost] [db CreatePost]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, "Failed To Create Post")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Post Created Sucessfully")
}

// GetAllPost -> Get All Post
func (cc PostController) GetAllPost(c *gin.Context) {

	pagination := utils.BuildPagination(c)
	pagination.Sort = "created_at desc"
	posts, count, err := cc.PostService.GetAllPost(pagination)

	if err != nil {
		cc.logger.Zap.Error("Error finding Post records", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find Post")
		return
	}
	responses.JSONCount(c, http.StatusOK, posts, count)

}

// GetOnePost -> Get One Post
func (cc PostController) GetOnePost(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	post, err := cc.PostService.GetOnePost(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [GetOnePost] [db GetOnePost]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find Post")
		return
	}
	responses.JSON(c, http.StatusOK, post)

}

// UpdateOnePost -> Update One Post By Id
func (cc PostController) UpdateOnePost(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	post := models.Post{}

	if err := c.ShouldBindJSON(&post); err != nil {
		cc.logger.Zap.Error("Error [UpdatePost] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, "failed to update post")
		return
	}
	post.ID = ID

	if err := cc.PostService.UpdateOnePost(post); err != nil {
		cc.logger.Zap.Error("Error [UpdatePost] [db UpdatePost]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, "failed to update post")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Post Updated Sucessfully")
}

// DeleteOnePost -> Delete One Post By Id
func (cc PostController) DeleteOnePost(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := cc.PostService.DeleteOnePost(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [DeleteOnePost] [db DeleteOnePost]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Delete Post")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Post Deleted Sucessfully")
}
