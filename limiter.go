package bestratelimiter

import (
	"errors"
	"github.com/Je33/bestratelimiter/model"
	"github.com/Je33/bestratelimiter/store"
	"time"
)

type Store interface {
	Add(key string, lim *model.Limit) error
	Set(key string, lim *model.Limit) error
	Get(key string) (*model.Limit, error)
}

type Limiter struct {
	store  Store
	config *LimiterConfig
}

func New(config Config) (*Limiter, error) {
	limiterStore, err := NewStore(&config.StoreConfig)
	if err != nil {
		return nil, err
	}

	return NewWithStore(limiterStore, &config.LimiterConfig), nil
}

func NewWithStore(store Store, config *LimiterConfig) *Limiter {
	return &Limiter{
		store:  store,
		config: config,
	}
}

func NewStore(config *StoreConfig) (Store, error) {
	return store.New(store.Config{
		Type: config.Type,
		URI:  config.URI,
	})
}

func (l *Limiter) Take(key string) (time.Duration, error) {
	lim, err := l.store.Get(key)

	if err != nil && !errors.Is(err, model.ErrKeyNotFound) {
		return l.GetDuration(lim), err
	}

	if lim == nil {
		lim = model.NewLimit()
		lim.Increment()
		err = l.store.Add(key, lim)
		if err != nil {
			return l.GetDuration(lim), err
		}

		return 0, nil
	}

	switch {
	case time.Now().Before(lim.GetLastAttempt().Add(l.config.Duration)):
		return l.GetDuration(lim), model.ErrRateLimit
	case lim.GetCount() >= l.config.Limit:
		if time.Now().Before(lim.GetFirstAttempt().Add(l.config.Period)) {
			return l.GetDuration(lim), model.ErrRateLimit
		}
		lim.Reset()
	default:
		lim.Increment()
	}

	err = l.store.Set(key, lim)
	if err != nil {
		return l.GetDuration(lim), err
	}

	return 0, nil
}

func (l *Limiter) Wait(key string) error {

	dur, err := l.Take(key)
	if err != nil {
		if dur > l.config.Timeout {
			return model.ErrTimeout
		}
		time.Sleep(dur)
	}

	return nil
}

func (l *Limiter) GetDuration(lim *model.Limit) time.Duration {
	dur := time.Until(lim.GetLastAttempt().Add(l.config.Duration))
	if lim.GetCount() >= l.config.Limit && time.Now().Add(dur).Before(lim.GetFirstAttempt().Add(l.config.Period)) {
		dur = time.Until(lim.GetFirstAttempt().Add(l.config.Period))
	}

	if dur < 0 {
		dur = 0
	}

	return dur
}
