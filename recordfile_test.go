package survey

import (
	"fmt"
	"testing"
)

func TestIntChange(t *testing.T) {
	rf, err := OpenRecordFile("bar1.txt", 100)
	if err != nil {
		t.Fatalf("open file error: %v", err)
	}
	defer rf.Close()

	values := make([]RecordValue, 0)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("intvalue%d", i)
		v := rf.AddRecordValue(key)
		if v == nil {
			t.Fatalf("add %s failed", key)
		}
		v.DumpInt(i)
		values = append(values, v)
	}

	for i, v := range values {
		v.DumpInt(-i)
	}
}

func BenchmarkSetIntRecord(b *testing.B) {
	rf, err := OpenRecordFile("bar2.txt", 1)
	if err != nil {
		b.Fatalf("open file error: %v", err)
	}
	defer rf.Close()

	v := rf.AddRecordValue("benchmark")
	for i := 0; i < b.N; i++ {
		v.DumpInt(i)
	}
}
