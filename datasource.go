package stream

import (
	"fmt"
	"reflect"
)

// DataSource data source function
type DataSource func(chan interface{})

// NewSliceDataSource new data source form slice
func NewSliceDataSource(array interface{}) DataSource {
	kind := reflect.TypeOf(array).Kind()
	if kind == reflect.Slice {
		source := func(in chan interface{}) {
			arrayValue := reflect.ValueOf(array)
			for i := 0; i < arrayValue.Len(); i++ {
				in <- arrayValue.Index(i).Interface()
			}
		}
		return source
	}
	panic(fmt.Sprintf("kind must be slice, real kind: %s", kind))
}
