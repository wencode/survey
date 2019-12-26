package survey

import (
	"testing"
)

func TestGen0to10(t *testing.T) {
	SetOutoutDir("./")
	a := NewAdder("foo", "adder")
	b := NewMaxer("foo", "maxer")
	c := NewMiner("foo", "miner")

	for i := 0; i <= 10; i++ {
		a.Put(i)
		b.Put(i)
		c.Put(i)
	}

	if r := a.Get(); r != 55 {
		t.Errorf("0-10 adder is %d error", r)
	}
	if r := b.Get(); r != 10 {
		t.Errorf("0-10 maxer is %d error", r)
	}
	if r := c.Get(); r != 0 {
		t.Errorf("0-10 miner is %d error", r)
	}

	Quit()
}

func TestMeanFunc(t *testing.T) {
	mean := &IntMean{}
	v := []int{1, 3, 5, 7, 9, 10}
	out := []int{1, 2, 3, 4, 5, 5}
	for i, n := range v {
		mean.Put(n)
		if a := mean.Get(); a != out[i] {
			t.Errorf("+%d mean %d != %d", n, a, out[i])
		}
	}
}
