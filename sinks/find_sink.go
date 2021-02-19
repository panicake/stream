package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.TerminalSink = &FindSink{}
}

type FindSink struct {
	value interface{}
}

func (f *FindSink) End(out chan interface{}) {
	for value := range out {
		f.value = value
	}
}

func NewFindSink() types.TerminalSink {
	return &FindSink{}
}

func (f *FindSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		out <- value
		break
	}
}

func (f *FindSink)Get() interface{}  {
	return f.value
}



