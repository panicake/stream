package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.ISink = &mapSink{}
}


type mapSink struct {
	mapper    types.Function
}

func NewMapSink(mapper types.Function) types.ISink {
	return &mapSink{
		mapper: mapper,
	}
}

func (f mapSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		out <- f.mapper(value)
	}
}


