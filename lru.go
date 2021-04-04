package scache

import (
	"sync/atomic"
	"time"
)

type lru struct {
	cacheImpl *cacheImpl
	maxCap    int32
	length    int32
	mp        map[string]*linkNode
	head      *linkNode
	tail      *linkNode
}

func newLRU(cap int) *lru {
	h, t := initHeadTail()
	return &lru{
		maxCap: int32(cap),
		length: 0,
		mp:     make(map[string]*linkNode),
		head:   h,
		tail:   t,
	}
}

func (l *lru) len() int {
	return int(l.length)
}

func (l *lru) add(key string, value interface{}, ttl time.Duration) {
	if n, ok := l.mp[key]; ok {
		n.value = value
		n.expiresAt = time.Now().Add(ttl)
		l.update(n)
	} else {
		newNode := &linkNode{key: key, value: value, expiresAt: time.Now().Add(ttl)}
		if l.length < l.maxCap || l.maxCap == 0 {
			l.headInsert(newNode)
			atomic.AddInt32(&l.length, 1)
		} else {
			l.replace(l.tail.pre, newNode)
		}
		l.mp[key] = newNode
	}
}

func (l *lru) get(key string) (interface{}, bool) {
	if n, ok := l.mp[key]; ok {
		// 惰性过期检查
		if time.Now().After(n.expiresAt) {
			l.del(n.key)
			return nil, false
		}

		l.update(n)
		return n.value, ok
	}
	return nil, false
}

func (l *lru) peek(key string) (interface{}, bool) {
	if n, ok := l.mp[key]; ok {
		// 惰性过期检查
		if time.Now().After(n.expiresAt) {
			l.del(n.key)
			return nil, false
		}
		return n.value, ok
	}
	return nil, false
}

func (l *lru) keys() []string {
	keys := make([]string, 0, l.length)
	t := time.Now()
	for ent := l.head.next; ent != l.tail; ent = ent.next {
		if t.After(ent.expiresAt) {
			l.del(ent.key)
			continue
		}
		keys = append(keys, ent.key)
	}
	return keys
}

func (l *lru) del(key string) {

	if n, ok := l.mp[key]; ok {
		n.exit()
		delete(l.mp, key)
		l.cacheImpl.stat.Evicted++
		atomic.AddInt32(&l.length, -1)

		if l.cacheImpl.opt.onEvicted != nil {
			l.cacheImpl.opt.onEvicted(n.key, n.value)
		}
	}
}

func (l *lru) tailDel() {
	tail := l.tail.pre
	// 避免删除头指针
	if tail != l.head {
		l.del(tail.key)
	}
}

func (l *lru) tailDelIfExpired() {
	tail := l.tail.pre
	// 避免删除头指针
	if tail != l.head && time.Now().After(tail.expiresAt) {
		l.del(tail.key)
	}
}

func initHeadTail() (*linkNode, *linkNode) {
	head, tail := &linkNode{}, &linkNode{}
	head.next = tail
	tail.pre = head
	return head, tail
}

func (l *lru) update(n *linkNode) {
	n.exit()
	l.headInsert(n)
}

func (l *lru) replace(dead, active *linkNode) {
	dead.exit()
	delete(l.mp, dead.key)
	l.headInsert(active)
}

func (l *lru) headInsert(n *linkNode) {
	first := l.head.next
	n.bridging(l.head, first)
	first.pre = n
	l.head.next = n
}
