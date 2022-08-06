package stream

var _ TerminalSink = &findSink{}

type findSink struct {
	value interface{}
}

// End data stream
func (f *findSink) End(out chan interface{}) {
	for value := range out {
		f.value = value
	}
}

func newFindSink() TerminalSink {
	return &findSink{}
}

// Flow data stream
func (f *findSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		out <- value
		break
	}
}

// Get result
func (f *findSink) Get() interface{} {
	return f.value
}
