package sinks

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.ISink = &distinctSink{}
}

type distinctSink struct {
	set mapset.Set
}

func (d *distinctSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		if d.set.Add(value) {
			out <- value
		}
	}
}

func NewDistinctSink() types.ISink {
	return &distinctSink{
		set: mapset.NewSet(),
	}
}
