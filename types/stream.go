package types

// IStream stream operation interface
type IStream interface {
	Parallel(threads int) IStream

	Filter(predicate Predicate) IStream
	Map(function Function) IStream
	Peek(action Consumer) IStream

	Limit(n int64) IStream
	Skip(n int64) IStream
	Distinct() IStream
	Sorted(comparator ComparatorFunc) IStream

	Max(comparator ComparatorFunc) interface{}
	Min(comparator ComparatorFunc) interface{}
	Reduce(function BiFunctionFunc) interface{}
	Count() int64

	FindFirst() interface{}
	AnyMatch(Predicate Predicate) bool
	AllMatch(Predicate Predicate) bool
	NonMatch(Predicate Predicate) bool

	ForEach(consumer ConsumerFunc)
	Collect(collector Collector) interface{}
}

// StreamHelper drives stream flows from the source to the end
type StreamHelper interface {
	Drive(sink TerminalSink) TerminalSink
}
