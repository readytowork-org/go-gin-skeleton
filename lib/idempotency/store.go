// Package idempotency provides a pluggable cache for replaying responses
// to retried requests carrying the same Idempotency-Key header.
//
// The middleware (see Middleware in this package) records the response for
// the first request with a given (route, key) tuple and replays it for any
// subsequent identical request within the TTL.
//
// Two backends are provided: MySQL (durable, single source of truth) and
// Redis (fast, TTL-native). The active backend is selected at startup from
// env.IdempotencyStore so applications can switch without code changes.
package idempotency

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound is returned by Store.Get when no record exists for the key.
// It is a sentinel — callers MUST check with errors.Is, since backends may
// wrap it (e.g. with additional context like the key that was missed).
var ErrNotFound = errors.New("idempotency: not found")

// Record is the cached representation of a previously-handled request's
// response. It is intentionally small: status + headers + body is enough
// to faithfully replay a JSON API response.
type Record struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

// Store is the persistence contract. Implementations MUST be safe for
// concurrent use, since the middleware is invoked from many goroutines.
type Store interface {
	// Get returns the stored record for the (route, key) pair, or ErrNotFound.
	Get(ctx context.Context, route, key string) (Record, error)
	// Set stores the record with the given TTL. A zero ttl means "forever",
	// which is rarely what you want — pass env.IdempotencyTTL.
	Set(ctx context.Context, route, key string, r Record, ttl time.Duration) error
}

// noopStore is selected when idempotency is disabled. It always reports
// "not found" and silently drops writes, so the middleware behaves as a
// no-op without conditional checks inside the hot path.
type noopStore struct{}

func (noopStore) Get(context.Context, string, string) (Record, error) {
	return Record{}, ErrNotFound
}
func (noopStore) Set(context.Context, string, string, Record, time.Duration) error { return nil }

// Noop returns a Store that disables idempotency.
func Noop() Store { return noopStore{} }
