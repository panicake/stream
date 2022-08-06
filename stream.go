package stream

import (
	"sync"
)

// Stream array concurrent Stream
type Stream struct {
	head    *Stream
	next    *Stream
	in      chan interface{}
	out     chan interface{}
	sink    Sink
	threads int
}

// NewStream init data stream from channel
func NewStream(source DataSource) *Stream {
	stream := &Stream{
		threads: 1,
		sink: SinkFunc(func(in, out chan interface{}) {
			source(out)
		}),
	}
	stream.head = stream
	return stream
}

// Parallel set threads count of Stream
func (s *Stream) Parallel(threads int) *Stream {
	if threads < 1 {
		threads = 1
	}
	next := s.spawn(SinkFunc(func(in, out chan interface{}) {
		for value := range in {
			out <- value
		}
	}))
	next.threads = threads
	return next
}

// Filter returns a stream consisting of elements filtered by predicate
func (s *Stream) Filter(predicate Predicate) *Stream {
	return s.spawn(newFilterSink(predicate))
}

// Map returns a stream consisting of results of applying the given function to the elements
func (s *Stream) Map(mapper Mapper) *Stream {
	return s.spawn(newMapSink(mapper))
}

// Peek returns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream
func (s *Stream) Peek(action Consumer) *Stream {
	return s.spawn(newPeekSink(action))
}

// Limit returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length
func (s *Stream) Limit(maxSize int64) *Stream {
	return s.spawn(newSliceSink(0, maxSize))
}

// Skip returns a stream consisting of the remaining elements of this stream after
// discarding the first n elements of th stream
// if this stream contains fewer than n elements then an empty stream will be returned
func (s *Stream) Skip(n int64) *Stream {
	return s.spawn(newSliceSink(n, -1))
}

// Distinct returns a stream consisting of the distinct elements
func (s *Stream) Distinct() *Stream {
	return s.spawn(newDistinctSink())
}

// Sorted returns a stream consisting of the sorted elements with comparator
func (s *Stream) Sorted(comparator ComparatorFunc) *Stream {
	s.threads = 1
	return s.spawn(newSortSink(comparator))
}

// Counting returns a stream consisting of the count of elements
func (s *Stream) Counting(count *int64) *Stream {
	return s.spawn(newCountingSink(count))
}

// Max returns a stream consisting of the max element with comparator
func (s *Stream) Max(comparator ComparatorFunc) interface{} {
	maxOperator := ReduceFunc(func(a interface{}, b interface{}) interface{} {
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}
		if comparator.Compare(a, b) {
			return a
		}
		return b
	})
	return s.evaluate(newReduceSink(maxOperator)).Get()
}

// Min returns a stream consisting of the min element with comparator
func (s *Stream) Min(comparator ComparatorFunc) interface{} {
	maxOperator := ReduceFunc(func(a interface{}, b interface{}) interface{} {
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}
		if !comparator.Compare(a, b) {
			return a
		}
		return b
	})
	return s.evaluate(newReduceSink(maxOperator)).Get()
}

// Count returns the count of element of this stream
func (s *Stream) Count() int64 {
	return s.evaluate(newCountSink()).Get().(int64)
}

// Reduce performs a reduction  on the elements of this stream, using an associative accumulation function,
// and returns an reduced value
func (s *Stream) Reduce(reduce ReduceFunc) interface{} {
	return s.evaluate(newReduceSink(reduce)).Get()
}

// FindFirst returns the first element of this stream or nil if the stream is empty
func (s *Stream) FindFirst() interface{} {
	s.threads = 1
	return s.evaluate(newFindSink()).Get()
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// returns true if any elements of the stream match the provided predicate, otherwise false
func (s *Stream) AnyMatch(predicate Predicate) bool {
	return s.evaluate(newMatchSink(predicate, MatchAny)).Get().(bool)
}

// AllMatch returns whether all elements of this stream match the provided predicate.
// returns true if either all elements of the stream match the provided predicate or the stream is empty,
// otherwise false
func (s *Stream) AllMatch(predicate Predicate) bool {
	return s.evaluate(newMatchSink(predicate, MatchAll)).Get().(bool)
}

// NonMatch returns whether no elements of this stream match the provided predicate
// returns true if either no elements of the stream match the provided predicate or the stream is empty, otherwise false
func (s *Stream) NonMatch(predicate Predicate) bool {
	return s.evaluate(newMatchSink(predicate, MathNone)).Get().(bool)
}

// ForEach performs an action for each element of this stream
func (s *Stream) ForEach(consumer Consumer) {
	s.evaluate(newForeachSink(consumer))
}

// Collect performs a mutable reduction operation on the elements of this stream using a collector
func (s *Stream) Collect(collector Collector) interface{} {
	s.ForEach(func(t interface{}) {
		collector.Receive(t)
	})
	return collector.Get()
}

// spawn new a stream connecting to the parent stream
func (s *Stream) spawn(sink Sink) *Stream {
	next := &Stream{head: s.head, threads: s.threads, sink: sink}
	s.next = next
	return next
}

// clone a stream
func (s *Stream) clone() *Stream {
	return &Stream{
		head:    s.head,
		next:    s.next,
		threads: s.threads,
		in:      s.in,
		out:     s.out,
		sink:    s.sink,
	}
}

// evaluate return a terminal sink evaluated by terminal operation
func (s *Stream) evaluate(sink TerminalSink) TerminalSink {
	return s.spawn(nil).drive(sink)
}

// flow data stream flows
func (s *Stream) flow(sink Sink) {
	var wg sync.WaitGroup
	wg.Add(s.threads)
	for i := 0; i < s.threads; i++ {
		go func() {
			sink.Flow(s.in, s.out)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(s.out)
	}()
}

// drive data stream flows from source to end
func (s *Stream) drive(sink TerminalSink) TerminalSink {
	for p := s.head; p != nil; p = p.next {
		ins := p.clone()
		if ins.out == nil {
			ins.out = make(chan interface{})
		}
		if p.next != nil {
			p.next.in = ins.out
		}
		if ins.sink != nil {
			ins.flow(p.sink)
		}
	}
	if s.out == nil {
		s.out = make(chan interface{})
	}
	s.flow(sink)
	sink.End(s.out)
	return sink
}
