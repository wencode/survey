package survey

import (
	"sync/atomic"
	"unsafe"
)

type periodism interface {
	tick()
	setName(name string)
	setDumper(d Dumper)
}

type periodNode struct {
	period periodism
	next   unsafe.Pointer
}

func (node *periodNode) append(period periodism) {
	newnode := unsafe.Pointer(&periodNode{
		period: period,
	})
	cur := node
	for {
		if atomic.CompareAndSwapPointer(&cur.next, nil, newnode) {
			return
		}
		cur = (*periodNode)(atomic.LoadPointer(&cur.next))
	}
}

func (node *periodNode) traverse(fn func(periodism)) {
	cur := node
	for {
		if cur.period != nil {
			fn(cur.period)
		}
		next := atomic.LoadPointer(&cur.next)
		if next == nil {
			break
		}
		cur = (*periodNode)(next)
	}
}
