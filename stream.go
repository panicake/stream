package go_stream

import (
	"github.com/lovermaker/stream/operation"
	"github.com/lovermaker/stream/sinks"
	"github.com/lovermaker/stream/types"
	"reflect"
	"sync"
)

func init() {
	var _ types.IStream = &Stream{}
	var _ types.StreamHelper = &Stream{}
}

type sinkGetter func() types.ISink

// Stream array concurrent types.IStream
type Stream struct {
	threads      int
	sourceStream *Stream
	nextSteam    *Stream
	in           chan interface{}
	out          chan interface{}
	OnSink       sinkGetter
}

// NewStream init a data stream from a slice
func NewStream(array interface{}) types.IStream {
	stream := &Stream{
		threads: 1,
	}
	kind := reflect.TypeOf(array).Kind()
	if kind == reflect.Slice {
		stream.out = make(chan interface{})
		go func() {
			arrayValue := reflect.ValueOf(array)
			for i := 0; i < arrayValue.Len(); i++ {
				stream.out <- arrayValue.Index(i).Interface()
			}
			close(stream.out)
		}()
	}
	stream.sourceStream = stream
	return stream
}

// NewStreamFromChannel init data stream from channel
func NewStreamFromChannel(in chan interface{}) types.IStream {
	stream := &Stream{
		threads: 1,
	}
	stream.out = make(chan interface{})
	go func() {
		for value := range in {
			stream.out <- value
		}
		close(stream.out)
	}()
	stream.sourceStream = stream
	return stream
}

// Filter returns a stream consisting of elements filtered by predicate
func (s *Stream) Filter(predicate types.Predicate) types.IStream {
	stream := s.spawn(func() types.ISink {
		return sinks.NewFilterSink(predicate)
	})
	return stream
}

// Map returns a stream consisting of results of applying the given function to the elements
func (s *Stream) Map(function types.Function) types.IStream {
	stream := s.spawn(func() types.ISink {
		return sinks.NewMapSink(function)
	})
	return stream
}

// Peek returns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream
func (s *Stream) Peek(action types.Consumer) types.IStream {
	stream := s.spawn(func() types.ISink {
		return sinks.NewPeekSink(action)
	})
	return stream
}

// Limit returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length
func (s *Stream) Limit(maxSize int64) types.IStream {
	stream := s.spawn(func() types.ISink {
		return sinks.NewSliceSink(0, maxSize)
	})
	return stream
}

// Skip returns a stream consisting of the remaining elements of this stream after
// discarding the first n elements of th stream
// if this stream contains fewer than n elements then an empty stream will be returned
func (s *Stream) Skip(n int64) types.IStream {
	stream := s.spawn(func() types.ISink {
		return sinks.NewSliceSink(n, -1)
	})
	return stream
}

// Distinct returns a stream consisting of the distinct elements
func (s *Stream) Distinct() types.IStream {
	stream := s.spawn(func() types.ISink {
		return sinks.NewDistinctSink()
	})
	return stream
}

// Sorted returns a stream consisting of the sorted elements with comparator
func (s *Stream) Sorted(comparator types.ComparatorFunc) types.IStream {
	s.threads = 1
	stream := s.spawn(func() types.ISink {
		return sinks.NewSortSink(comparator)
	})
	return stream
}

// Max returns a stream consisting of the max element with comparator
func (s *Stream) Max(comparator types.ComparatorFunc) interface{} {
	maxOperator := types.BiFunctionFunc(func(a interface{}, b interface{}) interface{} {
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}
		if comparator.Compare(a, b) >= 0 {
			return a
		}
		return b
	})
	sink := sinks.NewReduceSink(maxOperator)
	return s.evaluate(operation.NewReduceOp(sink)).(*sinks.ReduceSink).Get()
}

// Min returns a stream consisting of the min element with comparator
func (s *Stream) Min(comparator types.ComparatorFunc) interface{} {
	maxOperator := types.BiFunctionFunc(func(a interface{}, b interface{}) interface{} {
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}
		if comparator.Compare(a, b) <= 0 {
			return a
		}
		return b
	})
	sink := sinks.NewReduceSink(maxOperator)
	return s.evaluate(operation.NewReduceOp(sink)).(*sinks.ReduceSink).Get()
}

// Count returns the count of element of this stream
func (s *Stream) Count() int64 {
	return s.evaluate(operation.NewReduceOp(sinks.NewCountingSink())).Get().(int64)
}

// Reduce performs a reduction  on the elements of this stream, using an associative accumulation function,
// and returns an reduced value
func (s *Stream) Reduce(function types.BiFunctionFunc) interface{} {
	sink := sinks.NewReduceSink(function)
	return s.evaluate(operation.NewReduceOp(sink)).Get()
}

// FindFirst returns the first element of this stream or nil if the stream is empty
func (s *Stream) FindFirst() interface{} {
	s.threads = 1
	return s.evaluate(operation.NewReduceOp(sinks.NewFindSink())).Get()
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// returns true if any elements of the stream match the provided predicate, otherwise false
func (s *Stream) AnyMatch(predicate types.Predicate) bool {
	return s.evaluate(operation.NewReduceOp(sinks.NewMatchSink(predicate, sinks.MatchAny))).Get().(bool)
}
// AllMatch returns whether all elements of this stream match the provided predicate.
// returns true if either all elements of the stream match the provided predicate or the stream is empty, otherwise false
func (s *Stream) AllMatch(predicate types.Predicate) bool {
	return s.evaluate(operation.NewReduceOp(sinks.NewMatchSink(predicate, sinks.MatchAll))).Get().(bool)
}

// NonMatch returns whether no elements of this stream match the provided predicate
// returns true if either no elements of the stream match the provided predicate or the stream is empty, otherwise false
func (s *Stream) NonMatch(predicate types.Predicate) bool {
	return s.evaluate(operation.NewReduceOp(sinks.NewMatchSink(predicate, sinks.MathNone))).Get().(bool)
}

// ForEach performs an action for each element of this stream
func (s *Stream) ForEach(action types.ConsumerFunc) {
	s.evaluate(operation.NewReduceOp(sinks.NewForeachSink(action)))
}

// Collect performs a mutable reduction operation on the elements of this stream using a collector
func (s *Stream) Collect(collector types.Collector) interface{} {
	s.ForEach(func(t types.T) {
		collector.Receive(t)
	})
	return collector.Data()
}

// evaluate return a terminal sink evaluated by terminal operation
func (s *Stream) evaluate(op types.TerminalOp) types.TerminalSink {
	return op.Evaluate(s.spawn(nil))
}

// spawn new a stream connecting to the parent stream
func (s *Stream) spawn(getter sinkGetter) *Stream {
	if s.out == nil {
		s.out = make(chan interface{})
	}
	newSteam := &Stream{sourceStream: s.sourceStream, in: s.out, threads: s.threads, OnSink: getter}
	s.nextSteam = newSteam
	return newSteam
}

// Parallel set threads count of types.IStream
func (s *Stream) Parallel(threads int) types.IStream {
	s.threads = threads
	return s
}

// flow data stream flows
func (s *Stream) flow(sink types.ISink) {
	if s.out == nil {
		s.out = make(chan interface{})
	}
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

// Drive drives data stream flows from source to end
func (s *Stream) Drive(sink types.TerminalSink) types.TerminalSink {
	for p := s.sourceStream.nextSteam; p != nil; p = p.nextSteam {
		if p.OnSink != nil {
			p.flow(p.OnSink())
		}
	}
	s.flow(sink)
	sink.End(s.out)
	return sink
}
