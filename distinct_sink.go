package stream

import (
	mapset "github.com/deckarep/golang-set"
)

var _ Sink = &distinctSink{}

type distinctSink struct {
	set mapset.Set
}

// Flow data stream
func (d *distinctSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		if d.set.Add(value) {
			out <- value
		}
	}
}

func newDistinctSink() Sink {
	return &distinctSink{
		set: mapset.NewSet(),
	}
}
