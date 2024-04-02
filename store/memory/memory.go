package memory

import (
	"github.com/Je33/bestratelimiter/model"
	"sync"
	"time"
)

type Config struct {
	PurgeDuration time.Duration
}

type Client struct {
	db     *limits
	config Config
}

type limits struct {
	limits map[string]*limit
	mu     sync.RWMutex
}

type limit struct {
	firstAttempt time.Time
	lastAttempt  time.Time
	count        int
	mu           sync.RWMutex
}

func (l *limit) set(firstAttempt time.Time, lastAttempt time.Time, count int) error {
	l.mu.Lock()
	l.lastAttempt = lastAttempt
	l.firstAttempt = firstAttempt
	l.count = count
	l.mu.Unlock()

	return nil
}

func New(config Config) (*Client, error) {
	c := &Client{
		db: &limits{
			limits: make(map[string]*limit, 4096),
		},
		config: config,
	}

	if config.PurgeDuration > 0 {
		go c.purge()
	}

	return c, nil
}

func (c *Client) Add(key string, lim *model.Limit) error {
	err := c.add(key, lim.GetFirstAttempt(), lim.GetLastAttempt(), lim.GetCount())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Set(key string, lim *model.Limit) error {
	l, err := c.get(key)
	if err != nil {
		return err
	}

	err = l.set(lim.GetFirstAttempt(), lim.GetLastAttempt(), lim.GetCount())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(key string) (*model.Limit, error) {
	l, err := c.get(key)
	if err != nil {
		return nil, err
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	brl := new(model.Limit)
	brl.SetCount(l.count)
	brl.SetFirstAttempt(l.firstAttempt)
	brl.SetLastAttempt(l.lastAttempt)

	return brl, nil
}

func (c *Client) add(key string, firstAttempt time.Time, lastAttempt time.Time, count int) error {
	l := &limit{
		firstAttempt: firstAttempt,
		lastAttempt:  lastAttempt,
		count:        count,
	}

	c.db.mu.Lock()
	c.db.limits[key] = l
	c.db.mu.Unlock()

	return nil
}

func (c *Client) get(key string) (*limit, error) {
	c.db.mu.RLock()
	defer c.db.mu.RUnlock()

	l, ok := c.db.limits[key]
	if !ok {
		return nil, model.ErrKeyNotFound
	}

	return l, nil
}

func (c *Client) purge() {
	ticker := time.NewTicker(c.config.PurgeDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		}

		c.db.mu.Lock()
		for k, l := range c.db.limits {
			if l.lastAttempt.Before(l.lastAttempt.Add(-c.config.PurgeDuration)) {
				delete(c.db.limits, k)
			}
		}
		c.db.mu.Unlock()
	}
}
