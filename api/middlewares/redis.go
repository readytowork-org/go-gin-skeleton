package middlewares

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/responses"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RedisMiddleware struct {
	logger      infrastructure.Logger
	env         infrastructure.Env
	redisClient infrastructure.Redis
}

func NewRedisMiddleware(
	logger infrastructure.Logger,
	env infrastructure.Env,
	redisClient infrastructure.Redis,
) RedisMiddleware {
	return RedisMiddleware{
		logger:      logger,
		env:         env,
		redisClient: redisClient,
	}
}

// Verify Redis Cache
func (m RedisMiddleware) VerifyRedisCache() gin.HandlerFunc {
	return func(c *gin.Context) {

		endpoint := c.Request.URL

		cachedKey := endpoint.String()

		val, err := m.redisClient.RedisClient.Get(cachedKey).Bytes()
		if err != nil {
			c.Next()
			return
		}
		responseBody := map[string]interface{}{}

		json.Unmarshal(val, &responseBody)

		responses.InterfaceJson(c, http.StatusOK, responseBody)
		// Abort other chained middlewares since we already get the response here.
		c.Abort()
	}

}
