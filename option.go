package scache

import "time"

type OptionF func(*option)

type option struct {
	MaxKeys           int
	AutoClean         bool
	CleanIntervalTime time.Duration // 默认扫描间隔,因为进行清理的时候会堵塞，故建议自动清理间隔为TTL/2
	TTL               time.Duration
	onEvicted         func(key string, value interface{})
}

// 设定缓存的数量
func MaxKeys(maxKeys int) OptionF {
	return func(opt *option) {
		opt.MaxKeys = maxKeys
	}
}

// 开启自动清理,定时间隔gap
func AutoClean(gap time.Duration) OptionF {
	return func(opt *option) {
		opt.AutoClean = true
		opt.CleanIntervalTime = gap
	}
}

// 设定默认过期时间
func TTL(ttl time.Duration) OptionF {
	return func(opt *option) {
		opt.TTL = ttl
	}
}

// 设定删除时执行的函数
func OnEvicted(fn func(key string, value interface{})) OptionF {
	return func(opt *option) {
		opt.onEvicted = fn
	}
}
