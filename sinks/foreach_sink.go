package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init() {
	var _ types.TerminalSink = &foreachSink{}
}

type foreachSink struct {
	consumer types.Consumer
}

func (f foreachSink) End(out chan interface{}) {
	for range out {}
}

func (f foreachSink) Get() interface{} {
	return nil
}

func NewForeachSink(consumer types.Consumer) types.TerminalSink {
	return &foreachSink{
		consumer: consumer,
	}
}

func (f foreachSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		f.consumer.Accept(value)
	}
}
