package scache

import (
	"time"
)

const (
	defaultCap = 10000
	defaultTTL = 1 * time.Second
)

var defaultSCache, _ = NewCache(MaxKeys(defaultCap), TTL(defaultTTL), AutoClean(defaultTTL/2))

func Set(key string, data interface{}) {
	defaultSCache.Set(key, data, 0)
}

func Get(key string) (interface{}, bool) {
	return defaultSCache.Get(key)
}

func Del(key string) {
	defaultSCache.Invalidate(key)
}
