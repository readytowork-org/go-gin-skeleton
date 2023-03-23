package constants

import "time"

const (
	// Login Rate Limit
	LoginRateLimit int64 = 10
	//Login Perios
	LoginPeriod time.Duration = 1 * time.Minute

	// Basic Rate Limit
	BasicRateLimit int64 = 20
	//Basic Rate Period
	BasicPeriod time.Duration = 1 * time.Minute
)
