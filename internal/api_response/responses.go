package api_response

import (
	"boilerplate-api/internal/api_errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Message struct {
	Msg string `json:"message" validate:"required"`
} // @name Message

type Error struct {
	Error api_errors.ErrResponse[[]api_errors.ErrorContext] `json:"error" validate:"required"`
} // @name ApiError

type ErrorMsg struct {
	Error string `json:"error" validate:"required"`
} // @name ErrorMessage

type Data[T any] struct {
	Data T `json:"data" validate:"required"`
} // @name Data

type DataCount[T any] struct {
	Data  []T   `json:"data" validate:"required"`
	Count int64 `json:"count" validate:"required"`
} // @name DataCount

// JSON : json response function
func JSON[T interface{}](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, Data[T]{Data: data})
}

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, errData api_errors.ErrResponse[[]api_errors.ErrorContext]) {
	c.JSON(statusCode, Error{Error: errData})
}

// ErrorMessage : json error response function
func ErrorMessage(c *gin.Context, statusCode int, errData string) {
	c.JSON(statusCode, ErrorMsg{Error: errData})
}

// SuccessJSON : json error response function
func SuccessJSON(c *gin.Context, data string) {
	c.JSON(http.StatusOK, Message{Msg: data})
}

// JSONCount : json response function
func JSONCount[T interface{}](c *gin.Context, statusCode int, data []T, count int64) {
	c.JSON(statusCode, DataCount[T]{Data: data, Count: count})
}
