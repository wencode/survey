package survey

import "sync"

type BoundedIntQueue struct {
	locker sync.Mutex
	data   []int
	cursor int
}

func NewBoundedIntQueue(capacity int) *BoundedIntQueue {
	return &BoundedIntQueue{
		data: make([]int, 0, capacity),
	}
}

func (q *BoundedIntQueue) Push(value int) (head, length int) {
	q.locker.Lock()
	defer q.locker.Unlock()
	capacity := cap(q.data)
	if len(q.data) < capacity {
		q.data = append(q.data, value)
		q.cursor++
		if q.cursor == capacity {
			q.cursor = 0
		}
		return q.data[0], len(q.data)
	}
	head = q.data[q.cursor]
	q.data[q.cursor] = value
	q.cursor++
	if q.cursor == capacity {
		q.cursor = 0
	}
	return head, capacity
}

func (q *BoundedIntQueue) HeadTail() (int, int, bool) {
	q.locker.Lock()
	defer q.locker.Unlock()
	capacity := cap(q.data)
	if len := len(q.data); len < capacity {
		if len == 0 {
			return 0, 0, false
		}
		return q.data[0], q.data[len-1], true
	}
	if q.cursor == 0 {
		return q.data[0], q.data[capacity-1], true
	}
	return q.data[q.cursor], q.data[q.cursor-1], true
}
