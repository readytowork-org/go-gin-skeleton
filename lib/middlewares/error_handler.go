package middlewares

import (
	"boilerplate-api/lib/api_errors"
	"boilerplate-api/lib/config"

	"github.com/gin-gonic/gin"
)

// ErrorHandler converts errors pushed via c.Error() into the canonical
// Envelope response. Handlers can either call api_errors.RespondError(c, err)
// directly (and return) or simply push c.Error(err) and let this middleware
// flush at the end. Already-written responses are left untouched.
func ErrorHandler(logger config.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}

		last := c.Errors.Last()
		logger.Error("request error: ", last.Err)
		api_errors.RespondError(c, last.Err)
	}
}
