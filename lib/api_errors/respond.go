package api_errors

import (
	"errors"
	"net/http"

	"boilerplate-api/lib/constants"

	"github.com/gin-gonic/gin"
)

// Envelope is the unified shape returned for any error response.
// Handlers should not construct ad-hoc error JSON anymore — push a typed
// error via c.Error() (or call RespondError directly) and the ErrorHandler
// middleware will format the response.
type Envelope struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
} //	@name	ErrorEnvelope

// Common error codes returned to clients. Keep them stable; clients may
// branch on Code for UX (toast vs inline field error, etc.).
const (
	CodeBadRequest      = "bad_request"
	CodeValidation      = "validation_failed"
	CodeUnauthorized    = "unauthorized"
	CodeForbidden       = "forbidden"
	CodeNotFound        = "not_found"
	CodeConflict        = "conflict"
	CodeTooManyRequests = "too_many_requests"
	CodeInternal        = "internal_error"
)

// New constructs a typed error response. The returned value implements
// error, so callers can `return New(...)` or `c.Error(New(...))`.
func New(httpStatus int, code, message string) *ErrorResponse {
	return &ErrorResponse{
		Err:       errors.New(message),
		Message:   message,
		ErrorType: HttpErrorType(httpStatus),
		code:      code,
	}
}

// Wrap attaches an existing error to a typed response so callers can keep
// the original cause for logs while exposing a sanitized message.
func Wrap(err error, httpStatus int, code, message string) *ErrorResponse {
	return &ErrorResponse{
		Err:       err,
		Message:   message,
		ErrorType: HttpErrorType(httpStatus),
		code:      code,
	}
}

// WithValidation returns a 422 response carrying field-level details.
func WithValidation(details []ValidationError, message string) *ErrorResponse {
	if message == "" {
		message = "Invalid input information"
	}
	return &ErrorResponse{
		Err:            errors.New(message),
		Message:        message,
		ErrorType:      HttpErrorType(http.StatusUnprocessableEntity),
		ValidationErrs: &details,
		code:           CodeValidation,
	}
}

// RespondError writes the canonical Envelope for an error. Unknown errors
// fall through to a 500 to keep internals from leaking to clients.
func RespondError(c *gin.Context, err error) {
	traceID, _ := c.Get(constants.RequestID)
	traceStr, _ := traceID.(string)

	var er *ErrorResponse
	if errors.As(err, &er) {
		env := Envelope{
			Code:    er.code,
			Message: er.Message,
			TraceID: traceStr,
		}
		if er.ValidationErrs != nil {
			env.Details = *er.ValidationErrs
		}
		if env.Code == "" {
			env.Code = codeFromStatus(er.ErrorType.ToInt())
		}
		c.AbortWithStatusJSON(er.ErrorType.ToInt(), env)
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, Envelope{
		Code:    CodeInternal,
		Message: "Internal server error",
		TraceID: traceStr,
	})
}

func codeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return CodeBadRequest
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	case http.StatusTooManyRequests:
		return CodeTooManyRequests
	case http.StatusUnprocessableEntity:
		return CodeValidation
	default:
		return CodeInternal
	}
}
