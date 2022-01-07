package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"

	"github.com/gin-gonic/gin"
)

// UtilityService -> struct
type UtilityService struct {
	env    infrastructure.Env
	logger infrastructure.Logger
	bucket StorageBucketService
}

type Location struct {
	models.Coordinates
}
type Geometry struct {
	Location Location `json:"location"`
}

type Results struct {
	Geometry Geometry `json:"geometry"`
}

type GeoCodeRes struct {
	Results []Results `json:"results"`
	Status  string    `json:"status"`
}

// NewUtilityService -> creates a new Utility service
func NewUtilityService(env infrastructure.Env, logger infrastructure.Logger, bucket StorageBucketService) UtilityService {
	return UtilityService{
		env:    env,
		logger: logger,
		bucket: bucket,
	}
}

func (s UtilityService) GetCoordinateFromAddress(address string) (Location, error) {
	URL := "https://maps.googleapis.com/maps/api/geocode/json?key=" + s.env.GeoCodeApiKey + "&address=" + url.QueryEscape(address)
	s.logger.Zap.Info("URL :: ", URL)

	resp, err := http.Get(URL)
	if err != nil {
		return Location{}, err
	}

	defer resp.Body.Close()
	var apiResp GeoCodeRes

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return Location{}, err
	}

	return apiResp.Results[0].Geometry.Location, nil
}

func (uc UtilityService) SingleThumbnailImage(ctx *gin.Context, imgUrl string) (string, error) {

	imageName := imgUrl[strings.LastIndex(imgUrl, "/")+1:]
	fileExtension := strings.ToLower(imageName[strings.LastIndex(imageName, ".")+1:])
	thumbnailFileName := "images/thumbnail/" + imageName
	fileType := "image/" + fileExtension

	imgByte, err := utils.URLToBinary(imgUrl)
	if err != nil {
		return "", err
	}

	thumbnail, err := utils.CreateThumbnailFromByteString(imgByte, fileType, 200, 0)
	if err != nil {
		uc.logger.Zap.Error("Error Failed create thumbnail", err.Error())
		return "", err
	}

	uploadThumbnailUrl, err := uc.bucket.UploadThumbnailFile(ctx.Request.Context(), thumbnail, thumbnailFileName, fileExtension)
	if err != nil {
		uc.logger.Zap.Error("Error Failed to upload File::", err.Error())
		return "", err
	}

	return uploadThumbnailUrl, nil
}
