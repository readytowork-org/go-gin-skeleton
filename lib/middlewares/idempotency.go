package middlewares

import (
	"bytes"
	"net/http"

	"boilerplate-api/lib/config"
	"boilerplate-api/lib/idempotency"

	"github.com/gin-gonic/gin"
)

// IdempotencyHeader is the canonical client-supplied key. Stripe popularised
// this header; clients we expect to integrate with should already understand it.
const IdempotencyHeader = "Idempotency-Key"

// IdempotencyMiddleware constructs the gin middleware.
type IdempotencyMiddleware struct {
	store  idempotency.Store
	logger config.Logger
	env    config.Env
}

func NewIdempotencyMiddleware(store idempotency.Store, logger config.Logger, env config.Env) IdempotencyMiddleware {
	return IdempotencyMiddleware{store: store, logger: logger, env: env}
}

// Handle returns the middleware. Apply it to POST routes that should be
// safe to retry without producing duplicate side effects (payments, account
// creation, etc.). The header is opt-in: if the client doesn't send
// Idempotency-Key, the request flows through normally.
func (m IdempotencyMiddleware) Handle() gin.HandlerFunc {
	ttl := idempotency.TTL(m.env)

	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}
		key := c.GetHeader(IdempotencyHeader)
		if key == "" {
			c.Next()
			return
		}
		route := c.FullPath()
		if route == "" {
			c.Next()
			return
		}

		if rec, err := m.store.Get(c.Request.Context(), route, key); err == nil {
			for k, v := range rec.Headers {
				c.Writer.Header().Set(k, v)
			}
			c.Writer.Header().Set("Idempotent-Replay", "true")
			c.Data(rec.Status, rec.Headers["Content-Type"], rec.Body)
			c.Abort()
			return
		}

		// Wrap the writer so we can observe the response body without
		// preventing it from being sent to the client.
		bw := &bodyCaptureWriter{ResponseWriter: c.Writer, buf: &bytes.Buffer{}}
		c.Writer = bw
		c.Next()

		// Only cache responses that completed successfully. 5xx errors are
		// transient by definition; caching them would defeat the whole point.
		if bw.Status() < 200 || bw.Status() >= 500 {
			return
		}

		headers := map[string]string{}
		if ct := bw.Header().Get("Content-Type"); ct != "" {
			headers["Content-Type"] = ct
		}
		rec := idempotency.Record{
			Status:  bw.Status(),
			Headers: headers,
			Body:    bw.buf.Bytes(),
		}
		if err := m.store.Set(c.Request.Context(), route, key, rec, ttl); err != nil {
			m.logger.Warn("idempotency store write failed: ", err.Error())
		}
	}
}

// bodyCaptureWriter mirrors writes to an in-memory buffer so the handler's
// response can be replayed later. It still forwards everything to the real
// ResponseWriter, so the live request is unaffected.
type bodyCaptureWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyCaptureWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
