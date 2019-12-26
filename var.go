package survey

import (
	"golang.org/x/sys/cpu"
	"sync/atomic"
)

const (
	VarType_Current = iota
	VarType_Adder
	VarType_Maxer
	VarType_Miner
	VarType_IntMean
	VarType_Window
	VarType_UnitWindow
	VarType_Latency
)

type Var interface {
	Name() string
	Put(v int)
	Get() int
	Flush()
}

type Dumper interface {
	DumpInt(n int)
	DumpString(str string)
}

type VarParam struct {
	NoDump     bool
	windowSize int
	src        Var
}

type VarOption func(*VarParam)

func WithNoDump() VarOption {
	return func(param *VarParam) {
		param.NoDump = true
	}
}

func NewCurrent(family, name string, opts ...VarOption) *Current {
	v := &Current{}
	if !expose(v, family, name, parseOption(opts)) {
		return nil
	}
	return v
}

func NewAdder(family, name string, opts ...VarOption) *Adder {
	v := &Adder{}
	if !expose(v, family, name, parseOption(opts)) {
		return nil
	}
	return v
}

func NewMaxer(family, name string, opts ...VarOption) *Maxer {
	v := &Maxer{}
	if !expose(v, family, name, parseOption(opts)) {
		return nil
	}
	return v
}

func NewMiner(family, name string, opts ...VarOption) *Miner {
	v := &Miner{}
	if !expose(v, family, name, parseOption(opts)) {
		return nil
	}
	return v
}

func NewIntMean(family, name string, opts ...VarOption) *IntMean {
	v := &IntMean{}
	if !expose(v, family, name, parseOption(opts)) {
		return nil
	}
	return v
}

func parseOption(opts []VarOption) *VarParam {
	param := &VarParam{
		windowSize: 60,
	}
	for _, opt := range opts {
		opt(param)
	}
	return param
}

// VarOp must be satisfy:
//  - associative: 	(a Op b) Op c == a Op (b Op c)
//  - commutative: 	a Op b == b Op a
//  - identically equal: 	a Op b never change if a and b not change
type VarOp func(a, b int64) int64

type NilVar struct{}

func (_ NilVar) Name() string { return "nil" }
func (_ NilVar) Put(v int)    {}
func (_ NilVar) Get() int     { return 0 }
func (_ NilVar) Flush()       {}

type varbase struct {
	_      cpu.CacheLinePad
	value  int64
	dirty  int32
	name   string
	dumper Dumper
}

func (vb varbase) Name() string {
	return vb.name
}

func (vb varbase) Get() int {
	return int(vb.load())
}

func (vb *varbase) Flush() {
	if atomic.CompareAndSwapInt32(&vb.dirty, 1, 0) {
		if vb.dumper != nil {
			v := int(vb.load())
			vb.dumper.DumpInt(v)
		}
	}
}

func (vb *varbase) tick() { vb.Flush() }

func (vb varbase) load() int64 {
	return atomic.LoadInt64(&vb.value)
}

func (vb *varbase) store(a int64) {
	atomic.StoreInt64(&vb.value, a)
	atomic.StoreInt32(&vb.dirty, 1)
}

func (vb *varbase) exchange(expected, newvalue int64) bool {
	if atomic.CompareAndSwapInt64(&vb.value, expected, newvalue) {
		atomic.StoreInt32(&vb.dirty, 1)
		return true
	}
	return false
}

func (vb *varbase) modify(op VarOp, b int64) {
	var old = vb.load()
	var new = op(old, b)
	for !vb.exchange(old, new) {
		old = vb.load()
		new = op(old, b)
	}
	atomic.StoreInt32(&vb.dirty, 1)
}

func (vb *varbase) setName(name string) {
	vb.name = name
}

func (vb *varbase) setDumper(e Dumper) {
	vb.dumper = e
}

type Current struct {
	varbase
}

func (c *Current) Put(v int) {
	c.store(int64(v))
}

type Adder struct {
	varbase
}

func (a *Adder) Put(v int) { a.modify(add, int64(v)) }
func add(a, b int64) int64 { return a + b }

type Maxer struct {
	varbase
}

func (m *Maxer) Put(v int) { m.modify(max, int64(v)) }
func max(a, b int64) int64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

type Miner struct {
	varbase
}

func (m *Miner) Put(v int) { m.modify(min, int64(v)) }
func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

type IntMean struct {
	varbase
}

func (m IntMean) Get() int {
	total, count := m.totalCount()
	return int(total) / int(count)
}

func (m IntMean) GetFloat() float64 {
	total, count := m.totalCount()
	return float64(total) / float64(count)
}

func (m IntMean) totalCount() (total, count int64) {
	v := m.load()
	total = v & 0xFFFFFFFF
	count = v >> 32
	if count <= 0 {
		count = 1
	}
	return
}

func (m *IntMean) Put(v int) {
	b := int64(1<<32) | (int64(v) & 0xFFFFFFFF)
	m.modify(add, b)
}
