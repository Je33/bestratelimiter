package bestratelimiter

import (
	"github.com/Je33/bestratelimiter/store"
	"time"
)

type Config struct {
	LimiterConfig
	StoreConfig
}

type LimiterConfig struct {
	Period   time.Duration
	Limit    int
	Duration time.Duration
	Timeout  time.Duration
}

type StoreConfig struct {
	Type          store.Type
	URI           string
	PurgeDuration time.Duration
}
