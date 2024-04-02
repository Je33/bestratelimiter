package redis

import (
	"context"
	"encoding/json"
	"github.com/Je33/bestratelimiter/model"
	"github.com/redis/go-redis/v9"
	"time"
)

// Work in progress
// TODO: Add purge
// TODO: Add tests

type Config struct {
	URI           string
	PurgeDuration time.Duration
}

type Client struct {
	db     *redis.Client
	config Config
}

func New(config Config) (*Client, error) {
	opt, err := redis.ParseURL(config.URI)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	return &Client{
		db:     client,
		config: config,
	}, nil
}

func (c *Client) Add(key string, lim *model.Limit) error {
	val, err := json.Marshal(lim)
	if err != nil {
		return err
	}
	err = c.db.Set(context.Background(), key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Set(key string, lim *model.Limit) error {
	val, err := json.Marshal(lim)
	if err != nil {
		return err
	}
	err = c.db.Set(context.Background(), key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(key string) (*model.Limit, error) {
	val, err := c.db.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	dest := new(model.Limit)
	err = json.Unmarshal([]byte(val), &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}
