package types

// TerminalOp terminal operation at the end of the stream
type TerminalOp interface {
	Evaluate(helper StreamHelper) TerminalSink
}
