package bestratelimiter

import (
	"github.com/Je33/bestratelimiter/store"
	"time"
)

// Config is the Global config struct for bestratelimiter instance
type Config struct {
	LimiterConfig
	StoreConfig
}

// LimiterConfig is the config struct for Limiter
type LimiterConfig struct {
	// Period for limit
	Period time.Duration
	// Limit itself
	Limit int
	// Duration between takes
	Duration time.Duration
	// Timeout for waiting of next take
	Timeout time.Duration
}

// StoreConfig is the config struct for Store
type StoreConfig struct {
	// Storage type, support `memory` and `redis`
	Type store.Type
	// URI for store
	URI string
	// Duration for purge old records
	PurgeDuration time.Duration
}
