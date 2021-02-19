package collectors

import (
	"github.com/lovermaker/stream/types"
	"sync"
)

type ListCollector struct {
	sync.RWMutex
	data []interface{}
	consumer types.Consumer
}

func (n *ListCollector) Receive(value interface{}) {
	n.Lock()
	n.data = append(n.data, value)
	n.Unlock()
}

func (n *ListCollector) Data() interface{} {
	return n.data
}

type getter func(interface{})interface{}

type MapCollector struct {
	sync.RWMutex
	data map[interface{}]interface{}
	keyGetter getter
	valueGetter getter
}

func (n *MapCollector) Receive(value interface{}) {
	key := n.keyGetter(value)
	val := n.valueGetter(value)
	n.Lock()
	n.data[key] = val
	n.Unlock()
}

func (n *MapCollector) Data() interface{} {
	return n.data
}

// ToList collect results of the stream to a slice
func ToList() types.Collector {
	return &ListCollector{}
}

// ToMap collect results of the stream to a map
func ToMap(keyGetter getter, valueGetter getter) types.Collector {
	return &MapCollector{
		data: make(map[interface{}]interface{}),
		keyGetter:   keyGetter,
		valueGetter: valueGetter,
	}
}




