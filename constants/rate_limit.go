package constants

import "time"

const (
	RateLimitPeriod   = 1 * time.Minute
	RateLimitRequests = int64(10)
)