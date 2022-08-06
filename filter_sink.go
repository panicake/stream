package stream

var _ Sink = &filterSink{}

type filterSink struct {
	predicate Predicate
}

func newFilterSink(predicate Predicate) Sink {
	return &filterSink{
		predicate: predicate,
	}
}

// Flow data stream
func (f filterSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		if f.predicate(value) {
			out <- value
		}
	}
}
