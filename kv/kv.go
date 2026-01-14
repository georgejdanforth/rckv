package kv

import (
	"errors"
)

var (
	KeyNotFound = errors.New("Key not found")
)

type KVStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}
