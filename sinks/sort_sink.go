package sinks

import (
	"github.com/lovermaker/stream/types"
	"sort"
)

func init()  {
	var _ types.ISink = &SortSink{}
}

type SortSink struct {
	comparator types.Comparator
}


func NewSortSink(comparator types.Comparator) types.ISink {
	return &SortSink{
		comparator: comparator,
	}
}

func (c *SortSink) Flow(in chan interface{}, out chan interface{}) {
	array := make([]interface{}, 0)
	for value := range in {
		array = append(array, value)
	}
	sort.Slice(array, func(i, j int) bool {
		return c.comparator.Compare(array[i], array[j]) >= 0
	})
	for _, value := range array {
		out <- value
	}
}



