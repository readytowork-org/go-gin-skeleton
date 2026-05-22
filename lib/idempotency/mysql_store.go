package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// idempotencyKey mirrors the idempotency_keys table. Kept private to the
// package so callers can't accidentally write rows around the Store API.
type idempotencyKey struct {
	Route     string    `gorm:"column:route;primaryKey;size:255"`
	Key       string    `gorm:"column:key;primaryKey;size:255"`
	Status    int       `gorm:"column:status;not null"`
	Headers   string    `gorm:"column:headers;type:text"`
	Body      []byte    `gorm:"column:body;type:blob"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null;index"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (idempotencyKey) TableName() string { return "idempotency_keys" }

// MySQLStore persists idempotency records in MySQL via GORM.
// Expired rows are filtered at read time; a separate sweeper (or a TTL job)
// should periodically delete rows where expires_at < NOW() to keep the table
// from growing without bound.
type MySQLStore struct {
	db *gorm.DB
}

// NewMySQLStore returns a Store backed by the given *gorm.DB. The caller is
// expected to have created the idempotency_keys table via migrations.
func NewMySQLStore(db *gorm.DB) *MySQLStore {
	return &MySQLStore{db: db}
}

func (s *MySQLStore) Get(ctx context.Context, route, key string) (Record, error) {
	var row idempotencyKey
	err := s.db.WithContext(ctx).
		Where("route = ? AND `key` = ? AND expires_at > ?", route, key, time.Now()).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Record{}, ErrNotFound
		}
		return Record{}, err
	}

	rec := Record{Status: row.Status, Body: row.Body}
	if row.Headers != "" {
		_ = json.Unmarshal([]byte(row.Headers), &rec.Headers)
	}
	return rec, nil
}

func (s *MySQLStore) Set(ctx context.Context, route, key string, r Record, ttl time.Duration) error {
	headersJSON, err := json.Marshal(r.Headers)
	if err != nil {
		return err
	}
	expires := time.Now().Add(ttl)
	if ttl <= 0 {
		expires = time.Now().Add(24 * time.Hour)
	}
	row := idempotencyKey{
		Route:     route,
		Key:       key,
		Status:    r.Status,
		Headers:   string(headersJSON),
		Body:      r.Body,
		ExpiresAt: expires,
	}
	// Insert-or-ignore semantics: if a concurrent writer beat us to this
	// (route, key), keep their record so retries replay the original response.
	return s.db.WithContext(ctx).
		Where("route = ? AND `key` = ?", route, key).
		Attrs(row).
		FirstOrCreate(&row).Error
}
