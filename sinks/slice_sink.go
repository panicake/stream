package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.ISink = &sliceSink{}
}

type sliceSink struct {
	skip     int64
	limit    int64
}

const MaxUint = ^uint64(0)
const MinUint = 0
const MaxInt = int64(MaxUint >> 1)
const MinInt = -MaxInt - 1

func NewSliceSink(skip, limit int64) types.ISink  {
	sliceSink := &sliceSink{
		skip:       skip,
	}
	if limit >= 0 {
		sliceSink.limit = limit
	} else {
		sliceSink.limit = MaxInt
	}
	return sliceSink
}

func (s *sliceSink) Flow(in chan interface{}, out chan interface{}) {
	for value := range in {
		if s.skip == 0 {
			if s.limit > 0 {
				s.limit--
				out <- value
			}
		} else {
			s.skip --
		}
	}
}
