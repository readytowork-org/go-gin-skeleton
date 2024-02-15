package infrastructure

import (
	"github.com/go-redis/redis"
)

// Redis
type Redis struct {
	RedisClient redis.Client
}

func NewRedis(logger Logger, env Env) Redis {

	var client = redis.NewClient(&redis.Options{
		// Container name + port since we are using docker
		Addr: env.RedisAddress,
		// Password: env.RedisPassword,
	})

	if client == nil {
		logger.Zap.Fatalf("Cannot run redis")
	}

	return Redis{
		RedisClient: *client,
	}
}
