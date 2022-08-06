package stream

var _ Sink = &peekSink{}

type peekSink struct {
	action Consumer
}

// Flow data stream
func (d *peekSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		d.action(value)
		out <- value
	}
}

func newPeekSink(action Consumer) Sink {
	return &peekSink{
		action: action,
	}
}
