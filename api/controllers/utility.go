package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"boilerplate-api/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UtilityController struct {
	logger   infrastructure.Logger
	env      infrastructure.Env
	bucket   services.StorageBucketService
	s3Bucket services.S3BucketService
}

func NewUtilityController(logger infrastructure.Logger,
	env infrastructure.Env,
	bucket services.StorageBucketService,
	s3Bucket services.S3BucketService,
) UtilityController {
	return UtilityController{
		logger:   logger,
		env:      env,
		bucket:   bucket,
		s3Bucket: s3Bucket,
	}
}

// Response -> response for the util scope
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    string      `json:"data"`
	Path    string      `json:"path"`
	Value   interface{} `json:"attributes"`
}

const storageURL string = "https://storage.googleapis.com/"

// FileUploadHandler -> handles file upload
func (uc UtilityController) FileUploadHandler(ctx *gin.Context) {
	file, uploadFile, err := ctx.Request.FormFile("file")
	if err != nil {
		uc.logger.Zap.Error("Error Get File from request :: ", err.Error())
		responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to get file form request")
		return
	}

	fileExtension := filepath.Ext(uploadFile.Filename)
	fileName := utils.GenerateRandomFileName() + fileExtension
	originalFileName := "images/original/" + fileName
	thumbnailFileName := "images/thumbnail/" + fileName

	// File type
	file1, _, _ := ctx.Request.FormFile("file")
	fileHeader := make([]byte, 512)
	if _, err := file1.Read(fileHeader); err != nil {
		uc.logger.Zap.Error("Error File Read upload File::", err.Error())
		responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to read upload  File")
		return
	}
	fileType := http.DetectContentType(fileHeader)
	if fileType == "image/png" || fileType == "image/jpg" || fileType == "image/jpeg" || fileType == "image/gif" {
		uploadedOriginalURL, err := uc.bucket.UploadFile(ctx.Request.Context(), file, originalFileName)
		if err != nil {
			uc.logger.Zap.Error("Error Failed to upload File::", err.Error())
			responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to upload File")
			return
		}

		//uploadedthumbnail
		thumbnail, err := utils.CreateThumbnail(file, fileType, 200, 0)
		if err != nil {
			uc.logger.Zap.Error("Error Failed create thumbnail", err.Error())
			responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to upload File")
			return
		}
		uploadThumbnailUrl, err := uc.bucket.UploadThumbnailFile(ctx.Request.Context(), thumbnail, thumbnailFileName, fileExtension)
		if err != nil {
			uc.logger.Zap.Error("Error Failed to upload File::", err.Error())
			responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to upload File")
			return
		}

		response := &Response{
			Success: true,
			Message: "Uploaded Successfully",
			Data:    storageURL + uc.env.StorageBucketName + "/" + uploadedOriginalURL,
			Path:    uploadedOriginalURL,
			Value: map[string]string{
				"original_image_url":   storageURL + uc.env.StorageBucketName + "/" + uploadedOriginalURL,
				"original_image_path":  uploadedOriginalURL,
				"thumbnail_image_url":  storageURL + uc.env.StorageBucketName + "/" + uploadThumbnailUrl,
				"thumbnail_image_path": uploadThumbnailUrl,
			}}
		ctx.JSON(http.StatusOK, response)
		return
	}

	originalFileName = "files/" + fileName
	uploadedFileURL, err := uc.bucket.UploadFile(ctx.Request.Context(), file, originalFileName)
	if err != nil {
		uc.logger.Zap.Error("Error Failed to upload File::", err.Error())
		responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to upload file ")
		return
	}
	response := &Response{
		Success: true,
		Message: "Uploaded Successfully",
		Data:    storageURL + uc.env.StorageBucketName + "/" + uploadedFileURL,
		Path:    uploadedFileURL,
	}
	ctx.JSON(http.StatusOK, response)
}

// Input model
type Input struct {
	Path *string `form:"path" json:"path" binding:"required"`
}

// FileUploadS3Handler handles aws s3 file upload
func (cc UtilityController) FileUploadS3Handler(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		cc.logger.Zap.Error("Error Get File from request: ", err.Error())
		responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to get file form request")
		return
	}
	var input Input
	err = ctx.ShouldBind(&input)
	if err != nil {
		cc.logger.Zap.Error("Error Failed to bind input:: ", err.Error())
		responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to Bind")
		return
	}

	fileExtension := filepath.Ext(fileHeader.Filename)
	fileName := utils.GenerateRandomFileName() + fileExtension
	originalFileNamePath := *input.Path + "/" + fileName

	uploadedFileURL, err := cc.s3Bucket.UploadtoS3(file, fileHeader, originalFileNamePath)
	if err != nil {
		cc.logger.Zap.Error("Error Failed to upload File:: ", err.Error())
		responses.ErrorJSON(ctx, http.StatusBadRequest, "Failed to upload file")
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
