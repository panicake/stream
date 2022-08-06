package stream

var _ TerminalSink = &foreachSink{}

type foreachSink struct {
	consumer Consumer
}

// End data stream
func (f foreachSink) End(out chan interface{}) {
	for range out {
	}
}

// Get result
func (f foreachSink) Get() interface{} {
	return nil
}

func newForeachSink(consumer Consumer) TerminalSink {
	return &foreachSink{
		consumer: consumer,
	}
}

// Flow data stream
func (f foreachSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		f.consumer(value)
	}
}
