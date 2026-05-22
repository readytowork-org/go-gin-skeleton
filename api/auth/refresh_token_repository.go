package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"boilerplate-api/lib/config"

	"gorm.io/gorm"
)

// RefreshToken mirrors the refresh_tokens table.
type RefreshToken struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint32     `gorm:"column:user_id;not null;index"`
	TokenHash string     `gorm:"column:token_hash;type:varchar(255);not null;unique"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// TableName matches the migration explicitly so renames here don't silently
// reroute writes to a different table.
func (RefreshToken) TableName() string { return "refresh_tokens" }

// ErrRefreshTokenNotFound is returned when a token can't be matched (lookup
// fail). Treat this and ErrRefreshTokenRevoked identically in handlers — both
// surface as a generic 401 to avoid token-existence enumeration.
var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked or expired")
)

// RefreshTokenRepository persists and rotates refresh tokens.
type RefreshTokenRepository interface {
	Store(userID uint32, rawToken string, expiresAt time.Time) error
	FindActive(rawToken string) (RefreshToken, error)
	Revoke(id uint64) error
	RevokeAllForUser(userID uint32) error
}

type refreshTokenRepo struct {
	db config.Database
}

// NewRefreshTokenRepository creates the GORM-backed repo.
func NewRefreshTokenRepository(db config.Database) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

// hashToken returns a deterministic SHA-256 hex digest. We don't use bcrypt
// here because lookup-by-token must be O(1) — bcrypt's per-row salt forces a
// full table scan. SHA-256 of a high-entropy random token is acceptable: the
// token itself carries the secret, the hash is a fingerprint.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (r *refreshTokenRepo) Store(userID uint32, rawToken string, expiresAt time.Time) error {
	row := RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(rawToken),
		ExpiresAt: expiresAt,
	}
	return r.db.DB.Create(&row).Error
}

func (r *refreshTokenRepo) FindActive(rawToken string) (RefreshToken, error) {
	var rt RefreshToken
	err := r.db.DB.
		Where("token_hash = ?", hashToken(rawToken)).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rt, ErrRefreshTokenNotFound
		}
		return rt, err
	}
	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return rt, ErrRefreshTokenRevoked
	}
	return rt, nil
}

func (r *refreshTokenRepo) Revoke(id uint64) error {
	now := time.Now()
	return r.db.DB.Model(&RefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", &now).Error
}

func (r *refreshTokenRepo) RevokeAllForUser(userID uint32) error {
	now := time.Now()
	return r.db.DB.Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
}
