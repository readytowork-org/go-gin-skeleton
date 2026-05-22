package middlewares

import (
	"crypto/rand"
	"encoding/hex"

	"boilerplate-api/lib/constants"

	"github.com/gin-gonic/gin"
)

// RequestIDHeader is the canonical header used for request correlation.
const RequestIDHeader = "X-Request-ID"

// RequestID middleware ensures every request has a stable identifier in the
// gin context and response header so logs, error envelopes, and Sentry events
// can be cross-referenced.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		c.Set(constants.RequestID, id)
		c.Writer.Header().Set(RequestIDHeader, id)
		c.Next()
	}
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// extremely unlikely; fall back to empty string rather than panic
		return ""
	}
	return hex.EncodeToString(b[:])
}
