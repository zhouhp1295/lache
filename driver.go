package lache

import (
	"time"
)

type Driver interface {
	Get(key string) (result any, ok bool)
	GetT(key string, result any) (ok bool)
	Set(key string, value any, expiration time.Duration) (ok bool)
	Delete(key string) (ok bool)
}
