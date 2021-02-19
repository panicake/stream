package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init() {
	var _ types.TerminalSink = &ReduceSink{}
}

type ReduceSink struct {
	state    interface{}
	empty    bool
	operator types.BiFunction
}

func NewReduceSink(operator types.BiFunction) *ReduceSink {
	return &ReduceSink{ operator: operator, empty: true}
}

func (r *ReduceSink)reduce(in chan interface{}) interface{} {
	empty := true
	var state interface{}
	for value := range in {
		if empty || state == nil {
			empty = false
			state = value
		} else {
			state = r.operator.Apply(state, value)
		}
	}
	return state
}

func (r *ReduceSink) Flow(in chan interface{}, out chan interface{}) {
	state := r.reduce(in)
	out <- state
}

func (r *ReduceSink) End(out chan interface{}) {
	r.state = r.reduce(out)
}

func (r *ReduceSink) Get() types.T {
	return r.state
}



