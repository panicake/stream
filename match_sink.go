package stream

func init() {
	var _ TerminalSink = &matchSink{}
}

type matchKind string

const (
	MatchAll matchKind = "All"
	MatchAny matchKind = "Any"
	MathNone matchKind = "None"
)

type matchSink struct {
	predicate Predicate
	kind      matchKind
	value     bool
}

// End stream
func (f *matchSink) End(out chan interface{}) {
	f.value = f.match(out)
}

func newMatchSink(predicate Predicate, kind matchKind) TerminalSink {
	return &matchSink{
		predicate: predicate,
		kind:      kind,
	}
}

func (f *matchSink) match(in chan interface{}) bool {
	switch f.kind {
	case MatchAll:
		for value := range in {
			if !f.predicate(value) {
				return false
			}
		}
		return true
	case MatchAny:
		for value := range in {
			if f.predicate(value) {
				return true
			}
		}
		return false
	case MathNone:
		for value := range in {
			if f.predicate(value) {
				return false
			}
		}
		return true
	}
	return false
}

// Flow data stream
func (f *matchSink) Flow(in chan interface{}, out chan interface{}) {
	out <- f.match(in)
}

// Get result
func (f *matchSink) Get() interface{} {
	return f.value
}
