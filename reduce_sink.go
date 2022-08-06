package stream

var _ TerminalSink = &reduceSink{}

type reduceSink struct {
	state      interface{}
	empty      bool
	reduceFunc ReduceFunc
}

func newReduceSink(operator ReduceFunc) *reduceSink {
	return &reduceSink{reduceFunc: operator, empty: true}
}

func (r *reduceSink) reduce(in chan interface{}) interface{} {
	empty := true
	var state interface{}
	for value := range in {
		if empty || state == nil {
			empty = false
			state = value
		} else {
			state = r.reduceFunc(state, value)
		}
	}
	return state
}

// Flow data stream
func (r *reduceSink) Flow(in chan interface{}, out chan interface{}) {
	out <- r.reduce(in)
}

func (r *reduceSink) End(out chan interface{}) {
	r.state = r.reduce(out)
}

func (r *reduceSink) Get() interface{} {
	return r.state
}
