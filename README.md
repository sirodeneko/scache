<h1 align='center'>scache</h1>
<div align=center><img src="https://github.com/sirodeneko/scache/blob/master/rideGo.jpg"/></div>
<h2 align='center'>A Cache For Go</h2>

## 📖 简介

`scahe`是一个提供单机缓存池的第三方库，淘汰机制采用LRU(Least recently used,最近最少使用),支持自动过期，过期采用惰性删除+主动扫描。

## ⚠️ 注意

- 单机缓存适用场景较为局限，必须可容忍数据不一致。

## 🚀 功能

- 提供默认单机缓存池，最大容量为10000，扫描间隔为0.5s。
- 提供自定义缓存容量，到达容量则使用LRU算法淘汰。
- 容量到达最大值执行淘汰机制（LRU）。
- 支持设定key的过期时间。
- 采用惰性删除删除过期的key,即到使用到key的时候如果过期则进行删除。

## 🧰 安装
``` powershell
go get -u github.com/sirodeneko/scache
```

## 🛠 使用

### 可使用函数
```
// Set 放入key,value并设置过期时间，如果ttl为0，则采用cacheImpl定义的默认ttl
Set(key string, value interface{}, ttl time.Duration)
// Get 返回一个key的值，如果这个key没有过期
Get(key string) (interface{}, bool)
// Get 返回一个key的值，如果这个可以没有不存在或者过期了，将调用函数去获取
GetWithF(ctx context.Context, key string, ttl time.Duration, f GetF) (interface{}, error)
// Peek 返回一个key的值， 但是并不更新其“最近使用”的状态
Peek(key string) (interface{}, bool)
// Keys 返回所有的key，如果发现key过期了，会进行删除，故这也是一次扫描操作过期key的操作
Keys() []string
// Len 返回key的数量，包括过期了的
Len() int
// Invalidate 使一个key无效，即删除它
Invalidate(key string)
// InvalidateFn 扫描所有的key,如果函数返回值为true,删除这个key
InvalidateFn(fn func(key string) bool)
// RemoveOldest 删除最老的key
RemoveOldest()
// DeleteExpired 删除过期的key
DeleteExpired()
// Purge 清除所有的key
Purge()
// Stat 返回状态
Stat() Stats
```
### 可自定义的参数
```
// 设定缓存的数量
MaxKeys(maxKeys int) OptionF
// 开启自动清理,定时间隔gap
AutoClean(gap time.Duration) OptionF
// 设定默认过期时间
TTL(ttl time.Duration) OptionF 
// 设定删除时执行的函数
OnEvicted(fn func(key string, value interface{})) OptionF

```

### 使用默认缓存池
``` 
scache.Set("a","x")
scache.Get("a","x")
scache.Del("a","x")
```

### 自定义缓存池
```
// 最大key的长度为3，默认10毫秒后过期
cache, _ := NewCache(MaxKeys(3), TTL(time.Millisecond*10))
```

### 基本使用例子
```
// make cache with short TTL and 3 max keys
cache, _ := NewCache(MaxKeys(3), TTL(time.Millisecond*10))

// set value under key1.
// with 0 ttl (last parameter) will use cache-wide setting instead (10ms).
cache.Set("key1", "val1", 0)

// get value under key1
r, ok := cache.Get("key1")

// check for OK value, because otherwise return would be nil and
// type conversion will panic
if ok {
    rstr := r.(string) // convert cached value from interface{} to real type
    fmt.Printf("value before expiration is found: %v, value: %v\n", ok, rstr)
}

time.Sleep(time.Millisecond * 11)

// get value under key1 after key expiration
r, ok = cache.Get("key1")
// don't convert to string as with ok == false value would be nil
fmt.Printf("value after expiration is found: %v, value: %v\n", ok, r)

// set value under key2, would evict old entry because it is already expired.
// ttl (last parameter) overrides cache-wide ttl.
cache.Set("key2", "val2", time.Minute*5)

// Output:
// value before expiration is found: true, value: val1
// value after expiration is found: false, value: <nil>
```

