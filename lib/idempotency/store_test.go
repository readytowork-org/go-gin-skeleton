package idempotency

import (
	"context"
	"errors"
	"testing"
	"time"
)

// inMemoryStore is a tiny test-only store. We don't ship it (production gets
// MySQL or Redis), but it lets us validate the Store contract and any
// upstream middleware logic without a network round-trip.
type inMemoryStore struct {
	m map[string]Record
}

func newInMemoryStore() *inMemoryStore { return &inMemoryStore{m: map[string]Record{}} }

func (s *inMemoryStore) cacheKey(route, key string) string { return route + "|" + key }

func (s *inMemoryStore) Get(_ context.Context, route, key string) (Record, error) {
	if r, ok := s.m[s.cacheKey(route, key)]; ok {
		return r, nil
	}
	return Record{}, ErrNotFound
}

func (s *inMemoryStore) Set(_ context.Context, route, key string, r Record, _ time.Duration) error {
	if _, exists := s.m[s.cacheKey(route, key)]; exists {
		return nil
	}
	s.m[s.cacheKey(route, key)] = r
	return nil
}

func TestNoopStore_AlwaysMisses(t *testing.T) {
	s := Noop()
	if _, err := s.Get(context.Background(), "/r", "k"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("noop should return ErrNotFound, got %v", err)
	}
	if err := s.Set(context.Background(), "/r", "k", Record{Status: 200}, time.Minute); err != nil {
		t.Fatalf("noop Set should never error, got %v", err)
	}
}

func TestInMemoryStore_FirstWriterWins(t *testing.T) {
	s := newInMemoryStore()
	ctx := context.Background()
	rec1 := Record{Status: 201, Body: []byte(`{"id":1}`)}
	rec2 := Record{Status: 201, Body: []byte(`{"id":2}`)}

	if err := s.Set(ctx, "/users", "abc", rec1, time.Minute); err != nil {
		t.Fatalf("first Set failed: %v", err)
	}
	if err := s.Set(ctx, "/users", "abc", rec2, time.Minute); err != nil {
		t.Fatalf("second Set failed: %v", err)
	}

	got, err := s.Get(ctx, "/users", "abc")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(got.Body) != `{"id":1}` {
		t.Fatalf("expected first write to win, got %s", got.Body)
	}
}
