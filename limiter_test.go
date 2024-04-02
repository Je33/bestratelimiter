package bestratelimiter

import (
	"fmt"
	"github.com/Je33/bestratelimiter/model"
	"github.com/Je33/bestratelimiter/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	t.Parallel()

	t.Run("test limiter by count", func(t *testing.T) {
		t.Parallel()

		limiter, err := New(Config{
			LimiterConfig: LimiterConfig{
				Period:   time.Second,
				Limit:    10,
				Duration: 0,
				Timeout:  time.Second,
			},
			StoreConfig: StoreConfig{
				Type:          store.TypeMemory,
				PurgeDuration: 0,
			},
		})

		assert.NoError(t, err)

		for i := 1; i <= 10; i++ {
			dur, err := limiter.Take("some-key-1")
			assert.NoError(t, err)
			assert.Equal(t, time.Duration(0), dur)
			time.Sleep(10 * time.Millisecond)
		}

		dur, err := limiter.Take("some-key-1")
		assert.ErrorIs(t, err, model.ErrRateLimit)
		assert.Greater(t, dur, time.Duration(0))
	})

	t.Run("test limiter by duration", func(t *testing.T) {
		t.Parallel()

		limiter, err := New(Config{
			LimiterConfig: LimiterConfig{
				Period:   time.Second,
				Limit:    10,
				Duration: 30 * time.Millisecond,
				Timeout:  time.Second,
			},
			StoreConfig: StoreConfig{
				Type:          store.TypeMemory,
				PurgeDuration: 0,
			},
		})

		assert.NoError(t, err)

		for i := 1; i <= 2; i++ {
			dur, err := limiter.Take("some-key-2")
			assert.NoError(t, err)
			assert.Equal(t, time.Duration(0), dur)
			if i == 1 {
				time.Sleep(30 * time.Millisecond)
			}
		}

		dur, err := limiter.Take("some-key-2")
		assert.ErrorIs(t, err, model.ErrRateLimit)
		assert.Greater(t, dur, time.Duration(0))
	})

	t.Run("test limiter wait", func(t *testing.T) {
		t.Parallel()

		limiter, err := New(Config{
			LimiterConfig: LimiterConfig{
				Period:   time.Second,
				Limit:    10,
				Duration: 30 * time.Millisecond,
				Timeout:  time.Second,
			},
			StoreConfig: StoreConfig{
				Type:          store.TypeMemory,
				PurgeDuration: 0,
			},
		})

		assert.NoError(t, err)
		now := time.Now()

		for i := 1; i <= 10; i++ {
			err = limiter.Wait("some-key-3")
			assert.NoError(t, err)
			assert.Greater(t, time.Now().Sub(now), time.Duration(0))
			now = time.Now()
		}
	})
}

func BenchmarkLimiter_Take(b *testing.B) {

	b.Run("bench limiter with unique keys", func(b *testing.B) {
		limiter, err := New(Config{
			LimiterConfig: LimiterConfig{
				Period:   time.Minute,
				Limit:    b.N,
				Duration: 0,
				Timeout:  time.Second,
			},
			StoreConfig: StoreConfig{
				Type:          store.TypeMemory,
				PurgeDuration: time.Second,
			},
		})

		assert.NoError(b, err)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err = limiter.Take(fmt.Sprintf("some-key-%d", i))
			assert.NoError(b, err)
		}

		b.StopTimer()
	})

	b.Run("bench limiter with the same keys", func(b *testing.B) {
		limiter, err := New(Config{
			LimiterConfig: LimiterConfig{
				Period:   time.Minute,
				Limit:    b.N,
				Duration: 0,
				Timeout:  time.Second,
			},
			StoreConfig: StoreConfig{
				Type:          store.TypeMemory,
				PurgeDuration: time.Second,
			},
		})

		assert.NoError(b, err)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err = limiter.Take("some-key-4")
			assert.NoError(b, err)
		}

		b.StopTimer()
	})
}
