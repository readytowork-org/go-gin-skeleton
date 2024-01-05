package responses

import (
	"boilerplate-api/common/errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Msg string `json:"message"`
}

type Error struct {
	Error interface{} `json:"error"`
}

type Data struct {
	Data interface{} `json:"data"`
}

type DataCount struct {
	Data
	Count int64 `json:"count"`
}

// JSON : json response function
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Data{Data: data})
}

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, errData interface{}) {
	c.JSON(statusCode, Error{Error: errData})
}

// SuccessJSON : json error response function
func SuccessJSON(c *gin.Context, data string) {
	c.JSON(http.StatusOK, Message{Msg: data})
}

// JSONCount : json response function
func JSONCount(c *gin.Context, statusCode int, data interface{}, count int64) {
	c.JSON(statusCode, DataCount{Data: Data{Data: data}, Count: count})
}

func InterfaceJson(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

func UnauthorizedError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Message{Msg: "Unauthorized user"})
}

func CredentialsError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Message{Msg: "Please provide valid credentials"})
}

func InternalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, Message{Msg: "An error has occurred. Please try again later."})
}

type errResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Errors  interface{} `json:"errors"`
}

// HandleError func
func HandleError(c *gin.Context, err error) {
	status := errors.GetErrorType(err)

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
	if errorContext != nil {
		response.Errors = errorContext
	}
	ErrorJSON(c, int(status), &response)
}
