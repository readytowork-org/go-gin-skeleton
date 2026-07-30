package utility

import (
	"net/http"
	"path/filepath"

	"boilerplate-api/lib/api_errors"
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/utils"
	"boilerplate-api/services/aws_services"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	logger   config.Logger
	s3Bucket aws_services.S3BucketService
	service  Service
}

func NewController(
	logger config.Logger,
	s3Bucket aws_services.S3BucketService,
	service Service,
) Controller {
	return Controller{
		logger:   logger,
		s3Bucket: s3Bucket,
		service:  service,
	}
}

//	@Tags			UtilityApi
//	@Summary		handles file upload
//	@Description	handles file upload
//	@Security		Bearer
//	@Produce		application/json
//	@Param			file	formData	file		true	"Upload File"
//	@Success		200		{object}	Response	"File Uploaded Successfully"
//	@Failure		400		{object}	api_errors.Envelope
//	@Router			/api/v1/utils/file-upload [post]
//	@Id				FileUpload
func (uc Controller) FileUploadHandler(ctx *gin.Context) {
	file, uploadFile, err := ctx.Request.FormFile("file")
	if err != nil {
		api_errors.RespondError(ctx, api_errors.Wrap(err, http.StatusBadRequest, api_errors.CodeBadRequest, "Failed to get file from request"))
		return
	}

	message, response, err := uc.service.UploadImage(file, uploadFile)
	if err != nil {
		api_errors.RespondError(ctx, api_errors.Wrap(err, message.StatusCode, api_errors.CodeInternal, message.Message))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

//	@Tags			UtilityApi
//	@Summary		GetSignedUrl
//	@Description	generate signed url
//	@Security		Bearer
//	@Produce		application/json
//	@Param			image_url	query		string	false	"Image Url"
//	@Success		200			{object}	json_response.Data[string]
//	@Failure		400			{object}	json_response.Error[string]
//	@Router			/api/v1/utils/images/signed_url [get]
//	@Id				GetSignedUrl
func (uc Controller) GetSignedUrl(ctx *gin.Context) {
	imageUrl := ctx.Query("image_url")
	if imageUrl == "" {
		api_errors.RespondError(ctx, api_errors.New(http.StatusBadRequest, api_errors.CodeBadRequest, "Image Url is invalid"))
		return
	}

	signedUrl, err := uc.service.GetSignedUrl(imageUrl)
	if err != nil {
		api_errors.RespondError(ctx, api_errors.Wrap(err, http.StatusInternalServerError, api_errors.CodeInternal, "Failed to convert signed url"))
		return
	}

	ctx.Redirect(http.StatusFound, signedUrl)
}

// Input model
type Input struct {
	Path *string `form:"path" json:"path" binding:"required"`
}

// FileUploadS3Handler handles aws s3 file upload
func (uc Controller) FileUploadS3Handler(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		api_errors.RespondError(ctx, api_errors.Wrap(err, http.StatusBadRequest, api_errors.CodeBadRequest, "Failed to get file from request"))
		return
	}
	var input Input
	if err := ctx.ShouldBind(&input); err != nil {
		api_errors.RespondError(ctx, api_errors.Wrap(err, http.StatusBadRequest, api_errors.CodeBadRequest, "Failed to bind"))
		return
	}

	fileExtension := filepath.Ext(fileHeader.Filename)
	fileName := utils.GenerateRandomFileName() + fileExtension
	originalFileNamePath := *input.Path + "/" + fileName

	uploadedFileURL, err := uc.s3Bucket.UploadToS3(file, fileHeader, originalFileNamePath)
	if err != nil {
		api_errors.RespondError(ctx, api_errors.Wrap(err, http.StatusBadRequest, api_errors.CodeBadRequest, "Failed to upload file to s3 bucket"))
		return
	}

	response := &Response{
		Success: true,
		Message: "Uploaded Successfully",
		Path:    uploadedFileURL,
		Data:    uploadedFileURL,
	}
	ctx.JSON(http.StatusOK, response)
}
