package stream

// Sink data processor in stream
type Sink interface {
	// Flow data processing from in channel and output to out channel
	Flow(in chan interface{}, out chan interface{})
}

// SinkFunc function for sink
type SinkFunc func(in, out chan interface{})

// Flow data stream
func (s SinkFunc) Flow(in, out chan interface{}) {
	s(in, out)
}

// TerminalSink data processor at the end of the stream, usually collecting result
type TerminalSink interface {
	Sink
	// End collect result from out channel
	End(out chan interface{})
	// Get result
	Get() interface{}
}

// Comparator comparing interface
type Comparator interface {
	Compare(a interface{}, b interface{}) bool
}

// ComparatorFunc function implements Comparator
type ComparatorFunc func(a interface{}, b interface{}) bool

// Compare compare function
func (c ComparatorFunc) Compare(a interface{}, b interface{}) bool {
	return c(a, b)
}

// Consumer consume value
type Consumer func(interface{})

// ReduceFunc reduce two value into one value
type ReduceFunc func(interface{}, interface{}) interface{}

// Predicate returns true if value is needed
type Predicate = func(interface{}) bool

// Mapper map one value to another value
type Mapper func(interface{}) interface{}

// Collector collect result
type Collector interface {
	Receive(value interface{})
	Get() interface{}
}
