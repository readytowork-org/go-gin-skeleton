package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore persists idempotency records in Redis. Each record is stored
// under a key namespaced by route + idempotency key, with the TTL applied
// at write time so expiry is handled by Redis natively.
type RedisStore struct {
	client *redis.Client
	prefix string
}

// NewRedisStore wraps an existing *redis.Client. A namespace prefix avoids
// collisions when the same Redis instance is shared with other features.
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client, prefix: "idem:"}
}

func (s *RedisStore) cacheKey(route, key string) string {
	return s.prefix + route + ":" + key
}

func (s *RedisStore) Get(ctx context.Context, route, key string) (Record, error) {
	raw, err := s.client.Get(ctx, s.cacheKey(route, key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return Record{}, ErrNotFound
		}
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(raw, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func (s *RedisStore) Set(ctx context.Context, route, key string, r Record, ttl time.Duration) error {
	raw, err := json.Marshal(r)
	if err != nil {
		return err
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	// SetNX: do not overwrite an existing record. Matches the MySQL store's
	// "first writer wins" semantics so behaviour is consistent across backends.
	return s.client.SetNX(ctx, s.cacheKey(route, key), raw, ttl).Err()
}
