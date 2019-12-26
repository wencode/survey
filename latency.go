package survey

func NewLatency(family, prefix string, opts ...VarOption) *Latency {
	var (
		names = [4]string{
			prefix + "_latency",
			prefix + "_max_latency",
			prefix + "_count",
			prefix + "_qps",
		}
		param = parseOption(opts)
		cur   = NewCurrent(family, prefix, WithNoDump())
	)
	if cur == nil {
		return nil
	}
	l := &Latency{
		cur:   cur,
		avg:   NewUnitWindow(family, names[0], WithSrc(cur)),
		max:   NewMaxer(family, names[1]),
		count: NewAdder(family, names[2]),
		qps:   NewCurrent(family, names[3]),
	}
	param.NoDump = true
	if !expose(l, family, prefix, param) {
		return nil
	}
	return l
}

type Latency struct {
	cur   *Current
	avg   *UnitWindow
	max   *Maxer
	count *Adder
	qps   *Current

	prefix string
}

func (l *Latency) setName(name string) {
	l.prefix = name
}

func (l *Latency) setDumper(e Dumper) {}

func (l Latency) Name() string {
	return l.prefix
}

func (l *Latency) Put(v int) {
	l.cur.Put(v)
	l.max.Put(v)
	l.count.Put(1)
}

func (l Latency) Get() int {
	return l.avg.Get()
}

func (l *Latency) Flush() {
	l.avg.Flush()
	l.max.Flush()
	l.count.Flush()
}

func (l *Latency) tick() {
	l.avg.tick()
	l.max.tick()
	l.count.tick()
	avglatency := l.avg.Get()
	if avglatency > 0 {
		l.qps.Put(1000000 / avglatency)
	} else {
		l.qps.Put(0)
	}
	l.qps.Flush()
}
