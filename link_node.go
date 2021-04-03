package scache

import "time"

type linkNode struct {
	key       string
	value     interface{}
	expiresAt time.Time
	pre       *linkNode
	next      *linkNode
}

func (ln *linkNode) exit() {
	ln.pre.next = ln.next
	ln.next.pre = ln.pre
}

// 桥接
func (ln *linkNode) bridging(pf, nx *linkNode) {
	ln.pre = pf
	ln.next = nx
}
