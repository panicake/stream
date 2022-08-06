package stream

import "sync"

var _ Sink = &CountingSink{}

// CountingSink counting sick
type CountingSink struct {
	count *int64
	sync.Mutex
}

// newCountingSink new counting sink
func newCountingSink(count *int64) Sink {
	if count == nil {
		count = new(int64)
	}
	*count = 0
	return &CountingSink{
		count: count,
	}
}

// Flow flow stream
func (c *CountingSink) Flow(in chan interface{}, out chan interface{}) {
	for range in {
		c.Lock()
		*c.count++
		c.Unlock()
		out <- in
	}

}
