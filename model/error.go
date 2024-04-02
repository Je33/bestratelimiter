package model

import "errors"

var (
	ErrKeyNotFound      = errors.New("keyNotFound")
	ErrStoreTypeInvalid = errors.New("storeTypeInvalid")
	ErrRateLimit        = errors.New("rateLimit")
	ErrTimeout          = errors.New("timeout")
)
