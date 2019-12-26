package survey

import (
	"os"
	"reflect"
	"unsafe"

	"github.com/wencode/hack/mmap"
)

const (
	key_length   = 31
	value_length = 31
	line_length  = 64
)

type RecordValue []byte

func (rv RecordValue) DumpInt(v int) {
	u := v
	if v < 0 {
		u = -v
	}
	i := len(rv) - 1
	for i >= 1 {
		rv[i] = byte('0' + u%10)
		u /= 10
		i--
		if u == 0 {
			break
		}
	}
	if v < 0 {
		rv[i] = byte('-')
	}
}

func (rv RecordValue) DumpString(str string) {
	var (
		src []byte
		bh  = (*reflect.SliceHeader)(unsafe.Pointer(&src))
		sh  = (*reflect.StringHeader)(unsafe.Pointer(&str))
	)
	bh.Data, bh.Len = sh.Data, sh.Len
	bh.Cap = bh.Len
	copy([]byte(rv), src)
}

type RecordFile struct {
	*mmap.MapFile

	curmm mmap.MapBuf
	index int
}

func OpenRecordFile(filename string, recordNum int) (*RecordFile, error) {
	length := os.Getpagesize()
	for recordLen := recordNum * line_length; length < recordLen; length += os.Getpagesize() {
	}
	mf, err := mmap.Open(filename,
		mmap.WithLength(length),
		mmap.WithTruncate(),
		mmap.WithWrite())
	if err != nil {
		return nil, err
	}

	mm := mf.Buffer()
initChar:
	for i := 0; i < len(mm); i++ {
		for j := 0; j < line_length-1; j++ {
			mm[i] = byte(' ')
			i++
			if i >= len(mm) {
				break initChar
			}
		}
		mm[i] = byte('\n')
	}
	return &RecordFile{
		MapFile: mf,

		curmm: mf.Buffer(),
		index: 0,
	}, nil
}

func (rf *RecordFile) AddRecordValue(key string) RecordValue {
	valuebuf := rf.addLine(key)
	return RecordValue(valuebuf)
}

func (rf *RecordFile) addLine(key string) []byte {
	newindex := rf.index + line_length
	if newindex > len(rf.curmm) {
		rf.requestMapBuf()
		//todo
		return nil
		//if newindex > len(rf.curmm) {
		//	return nil
		//}
	}
	linebuf := rf.curmm[rf.index:newindex:newindex]
	rf.index = newindex
	keybuf := []byte(key)
	if len(keybuf) > key_length {
		keybuf = keybuf[:key_length]
	}
	n := copy(linebuf, keybuf)
	for i := n; i < key_length; i++ {
		linebuf[i] = byte(' ')
	}
	linebuf[key_length] = byte(':')
	linebuf[line_length-1] = byte('\n')
	return linebuf[key_length : line_length-1 : line_length-1]
}

func (rf *RecordFile) requestMapBuf() {

}
