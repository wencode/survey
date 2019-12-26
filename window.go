package survey

func WithSrc(src Var) VarOption {
	return func(param *VarParam) {
		param.src = src
	}
}

func WithWindowSize(windowSize int) VarOption {
	return func(param *VarParam) {
		param.windowSize = windowSize
	}
}

func WithWindowParam(windowSize int, src Var) VarOption {
	return func(param *VarParam) {
		param.windowSize = windowSize
		param.src = src
	}
}

func NewWindow(family, name string, opts ...VarOption) *Window {
	param := parseOption(opts)
	w := &Window{
		src:   param.src,
		queue: NewBoundedIntQueue(param.windowSize),
	}
	if !expose(w, family, name, param) {
		return nil
	}
	return w
}

func NewUnitWindow(family, name string, opts ...VarOption) *UnitWindow {
	param := parseOption(opts)
	w := &UnitWindow{
		src:   param.src,
		queue: NewBoundedIntQueue(param.windowSize),
	}
	if !expose(w, family, name, param) {
		return nil
	}
	return w
}

type Window struct {
	varbase
	src   Var
	queue *BoundedIntQueue
}

func (w *Window) tick() {
	v := w.src.Get()
	h, _ := w.queue.Push(v)
	w.store(int64(v - h))
	w.Flush()
}

func (_ *Window) Put(v int) {}

type UnitWindow struct {
	varbase
	src   Var
	queue *BoundedIntQueue
}

func (uw *UnitWindow) tick() {
	v := uw.src.Get()
	h, length := uw.queue.Push(v)
	uw.store(int64(v-h) / int64(length))
	uw.Flush()
}

func (_ *UnitWindow) Put(v int) {}
