package scache

import (
	"errors"
	"sync"
	"time"
)

// Cache 定义 缓存的接口
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Peek(key string) (interface{}, bool)
	Keys() []string
	Len() int
	Invalidate(key string)
	InvalidateFn(fn func(key string) bool)
	RemoveOldest()
	DeleteExpired()
	Purge()
	Stat() Stats
}

// Stats 记录缓存的状态，包括命中数，未命中数，添加的数量，删除的数量
type Stats struct {
	Hits    int
	Misses  int
	Added   int
	Evicted int
}

// cacheImpl  Cache 接口的实现
type cacheImpl struct {
	ttl       time.Duration
	maxKeys   int
	onEvicted func(key string, value interface{})
	cache     *lru
	opt       *option

	sync.Mutex
	stat Stats
}

// noEvictionTTL 默认过期时间（不过期）
const noEvictionTTL = time.Hour * 24 * 365 * 10

// NewCache 实例化一个缓存器
// 默认的MaxKeys为0，不限制缓存的数量
// 默认的过期时间为10年
// 默认不自动扫描过期的key
func NewCache(maxKeys int, options ...OptionF) (Cache, error) {
	if maxKeys < 0 {
		return nil, errors.New("the maxKeys is too low")
	}
	o := &option{}
	for _, f := range options {
		f(o)
	}

	s := &cacheImpl{
		opt:   o,
		cache: newLRU(maxKeys),
		ttl:   noEvictionTTL,
	}

	if o.AutoClean {
		go s.purgePeriodically()
	}

	return s, nil
}

// purgePeriodically 自动清理
// TODO
func (s *cacheImpl) purgePeriodically() {

}

// TODO
func (s *cacheImpl) Set(key string, value interface{}, ttl time.Duration) {
	panic("implement me")
}

// TODO
func (s *cacheImpl) Get(key string) (interface{}, bool) {
	panic("implement me")
}

// TODO
func (s *cacheImpl) Peek(key string) (interface{}, bool) {
	panic("implement me")
}

// TODO
func (s *cacheImpl) Keys() []string {
	panic("implement me")
}

// TODO
func (s *cacheImpl) Len() int {
	panic("implement me")
}

// TODO
func (s *cacheImpl) Invalidate(key string) {
	panic("implement me")
}

// TODO
func (s *cacheImpl) InvalidateFn(fn func(key string) bool) {
	panic("implement me")
}

// TODO
func (s *cacheImpl) RemoveOldest() {
	panic("implement me")
}

// TODO
func (s *cacheImpl) DeleteExpired() {
	panic("implement me")
}

// TODO
func (s *cacheImpl) Purge() {
	panic("implement me")
}

// TODO
func (s *cacheImpl) Stat() Stats {
	panic("implement me")
}
