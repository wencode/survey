package survey

import "testing"

func TestBoundedIntQueueNotFull(t *testing.T) {
	q := NewBoundedIntQueue(10)
	for i := 0; i < 10; i++ {
		q.Push(i)
		h, tail, ok := q.HeadTail()
		if !ok {
			t.Errorf("get data %d error", i)
		}
		if h != 0 || tail != i {
			t.Errorf("get head %d tail %d error", h, tail)
		}
	}
}

func TestBoundedIntQueueFull(t *testing.T) {
	q := NewBoundedIntQueue(10)
	for i := 0; i < 10; i++ {
		q.Push(i)
	}

	for i := 10; i < 100; i++ {
		q.Push(i)
		h, tail, ok := q.HeadTail()
		if !ok {
			t.Errorf("get data %d error", i)
		}
		if h != (i-9) || tail != i {
			t.Errorf("get head %d tail %d error", h, tail)
		}
	}
}
