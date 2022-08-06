package stream

import (
	"sync"
)

// ListCollector collect list data
type ListCollector struct {
	sync.RWMutex
	data     []interface{}
	consumer Consumer
}

// Receive stream data
func (n *ListCollector) Receive(value interface{}) {
	n.Lock()
	n.data = append(n.data, value)
	n.Unlock()
}

// Get result
func (n *ListCollector) Get() interface{} {
	return n.data
}

// MapCollector collect map
type MapCollector struct {
	sync.RWMutex
	data        map[interface{}]interface{}
	keyMapper   Mapper
	valueMapper Mapper
}

// Receive data
func (n *MapCollector) Receive(value interface{}) {
	key := n.keyMapper(value)
	val := n.valueMapper(value)
	n.Lock()
	n.data[key] = val
	n.Unlock()
}

// Get map result
func (n *MapCollector) Get() interface{} {
	return n.data
}

// ToList collect results of the stream to a slice
func ToList() Collector {
	return &ListCollector{}
}

// ToMap collect results of the stream to a map
func ToMap(keyMapper Mapper, valueMapper Mapper) Collector {
	return &MapCollector{
		data:        make(map[interface{}]interface{}),
		keyMapper:   keyMapper,
		valueMapper: valueMapper,
	}
}
