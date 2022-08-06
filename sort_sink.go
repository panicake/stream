package stream

import (
	"sort"
)

var _ Sink = &sortSink{}

// sortSink for sorting stream
type sortSink struct {
	comparator Comparator
}

func newSortSink(comparator Comparator) Sink {
	return &sortSink{
		comparator: comparator,
	}
}

// Flow data stream
func (c *sortSink) Flow(in chan interface{}, out chan interface{}) {
	array := make([]interface{}, 0)
	for value := range in {
		array = append(array, value)
	}
	sort.Slice(array, func(i, j int) bool {
		return c.comparator.Compare(array[i], array[j])
	})
	for _, value := range array {
		out <- value
	}
}
