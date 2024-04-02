# Go Rate Limiter

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/Je33/bestratelimiter)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/Je33/bestratelimiter/pipeline.yml?style=flat-square)](https://github.com/Je33/bestratelimiter/actions/workflows/pipeline.yml)


This package provides the `best` rate limiter in Go (Golang)
servers and distributed workloads. It's specifically designed with
simplicity and extensibility.


## Usage

1. Create a `Limiter` instance. This example uses an in-memory store:

    ```golang
    limiter, err := bestratelimiter.New(bestratelimiter.Config{
        LimiterConfig: bestratelimiter.LimiterConfig{
            // Period for limit
            Period:   time.Minute,
            // Limit itself
            Limit:    100,
            // Duration between takes
            Duration: 20 * time.Millisecond,
            // Timeout for waiting of next take
            Timeout:  time.Second,
        },
        StoreConfig: bestratelimiter.StoreConfig{
            // Storage type, support `memory` and `redis`
            Type:          store.TypeMemory,
            // In case of using in distributed systems with `redis`
			// URI: "redis://<uri>"
            // Time for periodically purging of old takes
            PurgeDuration: time.Second,
        },
    })
    if err != nil {
      log.Fatal(err)
    }
    ```

2. Try to get the limit by calling `Take()` on the `limiter`:

    ```golang
    // key is the unique value upon which you want to rate limit, like an IP or
    // MAC address.
    key := "127.0.0.1"
    duration, err := limiter.Take(key)

    // If taking falls with error we can get duration for next take
    if err != nil {
      return fmt.Errorf("rate limited: retry in %d", duration)
    }
   
   // Or we can Wait for next take
   err := limiter.Wait(key)
   
   // if error has returned to Wait it means Timeout of waiting reached
   if err != nil {
      return fmt.Error("take timeout reached")
   }
   ```

## Why _another_ Go rate limiter?

Just because. There are a lot of good rate limit libraries, but why not develop one more? 
Javascript frameworks are released every day, check here https://dayssincelastjavascriptframework.com/, 
and this is absolutely true. Why should we, gophers, lag behind?


### Speed and performance

How fast is it? It is not fast at all, and that's good, because always need space to grow. 
You can run the benchmarks yourself, but here's a few sample
benchmarks with the same and unique keys. I added commas to the output for clarity,
but you can run the benchmarks via `make benchmarks`:

```text
BenchmarkLimiter_Take
BenchmarkLimiter_Take/bench_limiter_with_unique_keys
BenchmarkLimiter_Take/bench_limiter_with_unique_keys-12                  2667558               429.3 ns/op           261 B/op          4 allocs/op
BenchmarkLimiter_Take/bench_limiter_with_unique_keys-12                  3013546               423.3 ns/op           252 B/op          4 allocs/op
BenchmarkLimiter_Take/bench_limiter_with_unique_keys-12                  2682770               447.7 ns/op           260 B/op          4 allocs/op
BenchmarkLimiter_Take/bench_limiter_with_unique_keys-12                  2855839               436.3 ns/op           256 B/op          4 allocs/op
BenchmarkLimiter_Take/bench_limiter_with_the_same_keys
BenchmarkLimiter_Take/bench_limiter_with_the_same_keys-12                7587766               153.5 ns/op            64 B/op          1 allocs/op
BenchmarkLimiter_Take/bench_limiter_with_the_same_keys-12                7804942               152.1 ns/op            64 B/op          1 allocs/op
BenchmarkLimiter_Take/bench_limiter_with_the_same_keys-12                7855465               151.7 ns/op            64 B/op          1 allocs/op
BenchmarkLimiter_Take/bench_limiter_with_the_same_keys-12                7885702               151.4 ns/op            64 B/op          1 allocs/op
```

There's likely still optimizations to be had, pull requests are welcome!


### Ecosystem

Many of the existing packages in the ecosystem take dependencies on other
packages. I'm an advocate of very thin libraries, and I don't think a rate
limiter should be pulling external packages. That's why **bestratelimiter uses only the
Go standard library**.


### Flexible and extensible

Most of the existing rate limiting libraries make a strong assumption that rate
limiting is only for HTTP services. Baked in that assumption are more
assumptions like rate limiting by "IP address" or are limited to a resolution of
"per second". The Best Rate Limiter can also be used to rate limit literally anything. 
It rate limits on a user-defined arbitrary string key.


### Stores

#### Memory

Memory is the fastest store, but only works on a single container/virtual
machine since there's no way to share the state.
[Learn more](https://pkg.go.dev/github.com/Je33/bestratelimiter/store/memory).

#### Redis

Limiter can be used with Redis, but with a performance cost.
[Learn more](https://pkg.go.dev/github.com/Je33/bestratelimiter/store/redis).
