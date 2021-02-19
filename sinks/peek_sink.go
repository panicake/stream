package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.ISink = &peekSink{}
}

type peekSink struct {
	action types.Consumer
}

func (d *peekSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		d.action.Accept(value)
		out <- value
	}
}

func NewPeekSink(action types.Consumer) types.ISink {
	return &peekSink{
		action: action,
	}
}

