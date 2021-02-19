package types

type Collector interface {
	Receive(value interface{})
	Data() interface{}
}



