package stream

var _ TerminalSink = &countSink{}

type countSink struct {
	count int64
}

// End data stream
func (c *countSink) End(out chan interface{}) {
	for value := range out {
		c.count += value.(int64)
	}
}

// Get result
func (c *countSink) Get() interface{} {
	return c.count
}

func newCountSink() TerminalSink {
	return &countSink{}
}

// Flow data stream
func (c *countSink) Flow(in chan interface{}, out chan interface{}) {
	count := int64(0)
	for range in {
		count++
	}
	out <- count
}
