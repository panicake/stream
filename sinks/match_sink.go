package sinks

import (
	"github.com/lovermaker/stream/types"
)

func init()  {
	var _ types.TerminalSink = &MatchSink{}
}

type MatchKind string

const (
	MatchAll MatchKind = "All"
	MatchAny MatchKind = "Any"
	MathNone MatchKind = "None"
)
type MatchSink struct {
	predicate types.Predicate
	kind MatchKind
	value bool
}

func (f *MatchSink) End(out chan interface{}) {
	f.value = f.match(out)
}

func NewMatchSink(predicate types.Predicate, kind MatchKind) types.TerminalSink {
	return &MatchSink{
		predicate: predicate,
		kind: kind,
	}
}

func (f *MatchSink)match(in chan interface{}) bool {
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

func (f *MatchSink) Flow(in chan interface{}, out chan interface{}) {
	out <- f.match(in)
}

func (f *MatchSink)Get() interface{}  {
	return f.value
}



