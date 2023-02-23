package middlewares

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// Global store
// using in-memory store with goroutine which clears expired keys.
var store = memory.NewStore()

type RateLimitOption struct {
	period time.Duration
	limit  int64
}

type Option func(*RateLimitOption)

type RateLimitMiddleware struct {
	logger infrastructure.Logger
	env infrastructure.Env
}

func NewRateLimitMiddleware(logger infrastructure.Logger,env infrastructure.Env) RateLimitMiddleware {
	return RateLimitMiddleware{
		logger: logger,
		env: env,
	}
}

func (rl RateLimitMiddleware) HandleRateLimit(limit int64, period time.Duration, options ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() // Gets cient IP Address

		rl.logger.Zap.Info("Setting up rate limit middleware...")

		// Limit -> # of API Calls
		// Period -> in a given time frame
		// setting default values
		opt := RateLimitOption{
			period: period,
			limit:  limit,
		}

		for _, o := range options {
			o(&opt)
		}

		rate := limiter.Rate{
			Limit:  opt.limit,
			Period: opt.period,
		}

		instance := limiter.New(store, rate)

		context, err := instance.Get(c, c.FullPath()+"&&"+key)

		if err != nil {
			rl.logger.Zap.Panic(err.Error())
		}

		c.Set(constants.RateLimit, instance)

		// Setting custom headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Limit exceeded
		if context.Reached {
		err := errors.TooManyRequests.New("Too many request")
		err = errors.SetCustomMessage(err, "Rate limit has exceeded")
		responses.HandleError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}

func WithOptions(period time.Duration, limit int64) Option {
	return func(o *RateLimitOption) {
		o.period = period
		o.limit = limit
	}
}
