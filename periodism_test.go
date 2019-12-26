package survey

import (
	"fmt"
	"testing"
)

type testPeriod int

func (i testPeriod) tick() {
	fmt.Printf("tick %d\n", int(i))
}

func (i testPeriod) setName(name string) {}
func (i testPeriod) setDumper(e Dumper)  {}

func TestPeriodNodeAppendTranvers(t *testing.T) {
	head := &periodNode{testPeriod(0), nil}
	for i := 1; i < 10; i++ {
		head.append(testPeriod(i))
	}

	i := 0
	head.traverse(func(period periodism) {
		period.tick()
		if int(period.(testPeriod)) != i {
			t.Errorf("transver node %d error", i)
		}
		i++
	})
}
