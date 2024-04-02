package store

import (
	"github.com/Je33/bestratelimiter/model"
	"github.com/Je33/bestratelimiter/store/memory"
	"github.com/Je33/bestratelimiter/store/redis"
	"time"
)

type Type string

const (
	TypeMemory Type = "memory"
	TypeRedis  Type = "redis"
)

func (s Type) IsValid() bool {
	switch s {
	case TypeMemory, TypeRedis:
		return true
	}

	return false
}

type Config struct {
	Type          Type
	URI           string
	PurgeDuration time.Duration
}

type client interface {
	Add(key string, lim *model.Limit) error
	Set(key string, lim *model.Limit) error
	Get(key string) (*model.Limit, error)
}

type Store struct {
	client client
}

func New(config Config) (*Store, error) {
	var c client
	var err error
	switch config.Type {
	case TypeMemory:
		c, err = memory.New(memory.Config{
			PurgeDuration: config.PurgeDuration,
		})
		if err != nil {
			return nil, err
		}
	case TypeRedis:
		c, err = redis.New(redis.Config{
			URI:           config.URI,
			PurgeDuration: config.PurgeDuration,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, model.ErrStoreTypeInvalid
	}
	return &Store{
		client: c,
	}, nil
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) Add(key string, lim *model.Limit) error {
	return s.client.Add(key, lim)
}

func (s *Store) Set(key string, lim *model.Limit) error {
	return s.client.Set(key, lim)
}

func (s *Store) Get(key string) (*model.Limit, error) {
	return s.client.Get(key)
}
