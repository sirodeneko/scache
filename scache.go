package scache

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Cache 定义 缓存的接口
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	GetWithF(ctx context.Context, key string, ttl time.Duration, f GetF) (interface{}, error)
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
	maxKeys int
	cache   *lru //并发不安全，需上锁
	opt     *option

	sync.RWMutex
	stat Stats
}

// noEvictionTTL 默认过期时间（不过期）
const noEvictionTTL = time.Hour * 24 * 365 * 10

// NewCache 实例化一个缓存器
// 默认的MaxKeys为0，不限制缓存的数量
// 默认的过期时间为10年
// 默认不自动扫描过期的key
func NewCache(options ...OptionF) (Cache, error) {
	o := &option{TTL: noEvictionTTL}
	for _, f := range options {
		f(o)
	}
	if o.MaxKeys < 0 {
		return nil, errors.New("the maxKeys is too low")
	}

	s := &cacheImpl{
		opt:   o,
		cache: newLRU(o.MaxKeys),
	}
	// 绑定cacheImpl
	s.cache.cacheImpl = s

	if o.AutoClean {
		go s.purgePeriodically()
	}

	return s, nil
}

// purgePeriodically 自动清理
func (s *cacheImpl) purgePeriodically() {
	heartbeat := time.NewTicker(s.opt.CleanIntervalTime)
	defer heartbeat.Stop()

	for range heartbeat.C {
		s.DeleteExpired()
	}
}

// Set 放入key,value并设置过期时间，如果ttl为0，则采用cacheImpl定义的默认ttl
func (s *cacheImpl) Set(key string, value interface{}, ttl time.Duration) {
	s.Lock()
	defer s.Unlock()

	if ttl == 0 {
		ttl = s.opt.TTL
	}
	s.cache.add(key, value, ttl)
	s.stat.Added++
}

// Get 返回一个key的值，如果这个key没有过期
func (s *cacheImpl) Get(key string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()

	reply, ok := s.cache.get(key)
	if ok {
		s.stat.Hits++
	} else {
		s.stat.Misses++
	}

	return reply, ok
}

// GetWithF 返回一个key的值，如果这个可以没有不存在或者过期了，将调用函数去获取
func (s *cacheImpl) GetWithF(ctx context.Context, key string, ttl time.Duration, f GetF) (interface{}, error) {
	s.RLock()
	reply, ok := s.cache.get(key)
	s.RUnlock()

	if !ok {
		reply, err := f(ctx)
		if err != nil {
			s.stat.Misses++
			return nil, err
		}
		s.stat.Added++
		s.stat.Misses++
		s.Lock()
		s.cache.add(key, reply, ttl)
		s.Unlock()
		return reply, nil
	}
	if ok {
		s.stat.Hits++
	} else {
		s.stat.Misses++
	}
	return reply, nil
}

// Peek 返回一个key的值， 但是并不更新其“最近使用”的状态
func (s *cacheImpl) Peek(key string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()

	reply, ok := s.cache.peek(key)
	if ok {
		s.stat.Hits++
	} else {
		s.stat.Misses++
	}

	return reply, ok
}

// Keys 返回所有的key，如果发现key过期了，会进行删除，故这也是一次扫描操作过期key的操作
func (s *cacheImpl) Keys() []string {
	s.Lock()
	defer s.Unlock()
	return s.cache.keys()
}

// Len 返回key的数量，包括过期了的
func (s *cacheImpl) Len() int {
	return s.cache.len()
}

// Invalidate 使一个key无效，即删除它
func (s *cacheImpl) Invalidate(key string) {
	s.Lock()
	defer s.Unlock()
	s.cache.del(key)
}

// InvalidateFn 扫描所有的key,如果函数返回值为true,删除这个key
func (s *cacheImpl) InvalidateFn(fn func(key string) bool) {
	s.Lock()
	defer s.Unlock()
	for key := range s.cache.mp {
		if fn(key) {
			s.cache.del(key)
		}
	}
}

// RemoveOldest 删除最老的key
func (s *cacheImpl) RemoveOldest() {
	s.Lock()
	defer s.Unlock()
	s.cache.tailDel()
}

// DeleteExpired 删除过期的key
func (s *cacheImpl) DeleteExpired() {
	s.Lock()
	defer s.Unlock()
	for key, value := range s.cache.mp {
		if time.Now().After(value.expiresAt) {
			s.cache.del(key)
		}
	}
}

// Purge 清除所有的key
func (s *cacheImpl) Purge() {
	s.Lock()
	defer s.Unlock()

	for k, v := range s.cache.mp {
		delete(s.cache.mp, k)
		s.stat.Evicted++
		if s.opt.onEvicted != nil {
			s.opt.onEvicted(k, v.value)
		}
	}
	s.cache = newLRU(s.maxKeys)
}

// Stat 返回状态
func (s *cacheImpl) Stat() Stats {
	return s.stat
}
