package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.TerminalSink = &CountingSink{}
}


type CountingSink struct {
	count int64
}

func (c *CountingSink) End(out chan interface{}) {
	for value := range out {
		c.count += value.(int64)
	}
}

func (c *CountingSink) Get() interface{} {
	return c.count
}

func NewCountingSink() types.TerminalSink {
	return &CountingSink{}
}

func (c *CountingSink) Flow(in chan interface{}, out chan interface{}) {
	count := int64(0)
	for range in {
		count++
	}
	out <- count
}



