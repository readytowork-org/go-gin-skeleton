package responses

import (
	"boilerplate-api/config/errors"
	"os"

	"github.com/gin-gonic/gin"
)

// JSON : json response function
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
}

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"error": data})
}

// SuccessJSON : json error response function
func SuccessJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"msg": data})
}

// JSONCount : json response function
func JSONCount(c *gin.Context, statusCode int, data interface{}, count int64) {
	c.JSON(statusCode, gin.H{"data": data, "count": count})
}

type errResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Errors  interface{} `json:"errors"`
}

// HandleError func
func HandleError(c *gin.Context, err error) {
	errorType := errors.GetErrorType(err)
	status := errors.GetStatusCode(errorType)

	errorContext := errors.GetErrorContext(err)
	customMessage := errors.GetCustomMessage(err)
	response := &errResponse{}
	if os.Getenv("ENV") != "production" {
		response.Error = err.Error()
	}
	response.Message = customMessage

	if customMessage == "" {
		response.Message = "An error has occurred. Please try again later."
	}
	if status == 500 {
		response.Message = "An error has occurred. Please try again later."
	}
	if errorContext != nil {
		response.Errors = errorContext
	}
	c.JSON(status, gin.H{"error": &response})
}
