package operation

import (
	"github.com/lovermaker/stream/types"
)

func init() {
	var _ types.TerminalOp = &reduceOp{}
}
type reduceOp struct {
	sink types.TerminalSink
}

func NewReduceOp(sink types.TerminalSink) types.TerminalOp {
	return &reduceOp{sink: sink}
}

func (r reduceOp) Evaluate(helper types.StreamHelper) types.TerminalSink {
	return helper.Drive(r.sink)
}

