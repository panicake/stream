package stream

import "sync"

var _ Sink = &sliceSink{}

type sliceSink struct {
	skip  int64
	limit int64
	sync.Mutex
}

const MaxUint = ^uint64(0)
const MinUint = 0
const MaxInt = int64(MaxUint >> 1)
const MinInt = -MaxInt - 1

func newSliceSink(skip, limit int64) Sink {
	sink := &sliceSink{
		skip: skip,
	}
	if limit >= 0 {
		sink.limit = limit
	} else {
		sink.limit = MaxInt
	}
	return sink
}

func (s *sliceSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		s.Lock()
		if s.skip == 0 {
			if s.limit > 0 {
				s.limit--
				out <- value
			}
		} else {
			s.skip--
		}
		s.Unlock()
	}
}
