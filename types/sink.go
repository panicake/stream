package types

// ISink data processor in stream
type ISink interface {
	// Accept data processing from in channel and output to out channel
	Flow(in chan interface{}, out chan interface{})
}

// TerminalSink data processor at the end of the stream, usually collecting result
type TerminalSink interface {
	ISink
	// End collect result from out channel
	End(out chan interface{})
	// Get get result
	Get() interface{}
}
