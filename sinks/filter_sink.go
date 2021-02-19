package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.ISink = &filterSink{}
}

type filterSink struct {
	predicate  types.Predicate
}

func NewFilterSink(predicate  types.Predicate) types.ISink {
	return &filterSink{
		predicate: predicate,
	}
}

func (f filterSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		if f.predicate(value) {
			out <- value
		}
	}
}