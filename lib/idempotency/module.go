package idempotency

import (
	"time"

	"boilerplate-api/lib/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Module wires the idempotency Store and a Redis client (when configured) into
// the fx graph. The Store is selected by env.IdempotencyStore:
//   - "mysql"  -> MySQLStore backed by the existing GORM DB.
//   - "redis"  -> RedisStore (requires REDIS_ADDR).
//   - ""       -> Noop (idempotency disabled). Default to keep behaviour
//                 unchanged for existing deployments.
var Module = fx.Module(
	"idempotency",
	fx.Provide(NewRedisClient),
	fx.Provide(NewStore),
)

// NewRedisClient returns a redis.Client when REDIS_ADDR is set, otherwise nil.
// Returning nil is deliberate: the Store factory branches on store-type, and a
// nil client is fine to pass around for MySQL/noop selections.
func NewRedisClient(env config.Env) *redis.Client {
	if env.RedisAddr == "" {
		return nil
	}
	return redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})
}

// NewStore selects an idempotency backend at startup. Invalid combinations
// (e.g. redis selected without REDIS_ADDR) fall back to Noop with a log line
// rather than panicking — env validation should catch the typo earlier; this
// is a belt-and-braces guard.
func NewStore(env config.Env, db config.Database, redisClient *redis.Client, logger config.Logger) Store {
	switch env.IdempotencyStore {
	case "mysql":
		logger.Info("idempotency: using MySQL store")
		return NewMySQLStore(db.DB)
	case "redis":
		if redisClient == nil {
			logger.Warn("idempotency: REDIS_ADDR not set, disabling idempotency")
			return Noop()
		}
		logger.Info("idempotency: using Redis store")
		return NewRedisStore(redisClient)
	default:
		logger.Info("idempotency: disabled (set IDEMPOTENCY_STORE=mysql|redis to enable)")
		return Noop()
	}
}

// TTL returns env.IdempotencyTTL with a 24h fallback.
func TTL(env config.Env) time.Duration {
	if env.IdempotencyTTL > 0 {
		return env.IdempotencyTTL
	}
	return 24 * time.Hour
}
