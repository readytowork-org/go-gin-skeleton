package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type HttpErrorType int

//these are unsigned integer constants of custom type HttpErrorType
const (
	BadRequest    = HttpErrorType(http.StatusBadRequest)
	Unauthorized  = HttpErrorType(http.StatusUnauthorized)
	Forbidden     = HttpErrorType(http.StatusForbidden)
	NotFound      = HttpErrorType(http.StatusNotFound)
	Conflict      = HttpErrorType(http.StatusConflict)
	InternalError = HttpErrorType(http.StatusInternalServerError)
	Unavailable   = HttpErrorType(http.StatusServiceUnavailable)
)

type responseError struct {
	httpErrorType HttpErrorType
	originalError error
	customMessage CustomMessage
	contextInfo   ErrorContext
	context       []ErrorContext
}

// ErrorContext struct
type ErrorContext struct {
	Field   string
	Message string
}

//CustomMessage struct
type CustomMessage struct {
	Message string
}

// Error() returns the origin message
func (error responseError) Error() string {
	return error.originalError.Error()
}

// ----------------- HTTPErrorType ----------------------

// New func :
func (errorType HttpErrorType) New(msg string) error {
	return responseError{
		httpErrorType: errorType,
		originalError: errors.New(msg),
	}
}

// Newf func :
func (errorType HttpErrorType) Newf(msg string, args ...interface{}) error {
	return responseError{
		httpErrorType: errorType,
		originalError: fmt.Errorf(msg, args...),
	}
}

// Wrap func:
func (errorType HttpErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrapf func:
func (errorType HttpErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return responseError{
		httpErrorType: errorType,
		originalError: errors.Wrapf(err, msg, args...),
	}
}

// GetErrorType func
func GetErrorType(err error) HttpErrorType {
	if responseErr, ok := err.(responseError); ok {
		return responseErr.httpErrorType
	}
	return InternalError
}

// ---------- Ends  HttpErrorType ------------

// AddErrorContext func:
func AddErrorContext(err error, field, message string) error {
	context := ErrorContext{Field: field, Message: message}
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			httpErrorType: responseErr.httpErrorType,
			originalError: responseErr.originalError,
			customMessage: responseErr.customMessage,
			context:       append(responseErr.context, context),
		}
	}

	re := &responseError{}
	re.httpErrorType = InternalError
	re.context = append(re.context, context)
	re.originalError = err
	return re
}

// AddErrorContextBlock func:
func AddErrorContextBlock(err error, errorCtx []ErrorContext) error {
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			httpErrorType: responseErr.httpErrorType,
			originalError: responseErr.originalError,
			customMessage: responseErr.customMessage,
			context:       append(responseErr.context, errorCtx...),
		}
	}

	re := &responseError{}
	re.httpErrorType = InternalError
	re.context = append(re.context, errorCtx...)
	re.originalError = err
	return re
}

// GetErrorContext func
func GetErrorContext(err error) []ErrorContext {
	if responseErr, ok := err.(responseError); ok && responseErr.context != nil {
		return responseErr.context
	}
	return nil
}

func SetCustomMessage(err error, msg string) error {
	customMessage := CustomMessage{Message: msg}
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			httpErrorType: responseErr.httpErrorType,
			originalError: responseErr.originalError,
			customMessage: customMessage,
			context:       responseErr.context,
		}
	}
	return responseError{
		httpErrorType: InternalError,
		originalError: err,
		customMessage: customMessage,
	}
}

func GetCustomMessage(err error) string {
	emptyStruct := CustomMessage{}
	if responseErr, ok := err.(responseError); ok && responseErr.customMessage != emptyStruct {
		return responseErr.customMessage.Message
	}
	return ""
}
