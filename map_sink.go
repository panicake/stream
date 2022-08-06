package stream

var _ Sink = &mapSink{}

type mapSink struct {
	mapper Mapper
}

func newMapSink(mapper Mapper) Sink {
	return &mapSink{
		mapper: mapper,
	}
}

// Flow data stream
func (f mapSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		out <- f.mapper(value)
	}
}
