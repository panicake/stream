## Stream API for Go, inspired by Java 8 Streams

Provides stream API for go applications like Java 8 streams, for example:

```go
nums := make([]int, 0)
for i := 0; i < 100; i++ {
    nums = append(nums, i)
}
sum := NewStream(nums).Filter(func(t types.T) bool {
    return t.(int) % 3 == 0
}).Map(func(t types.T) types.R {
    return 2 * t.(int)
}).Reduce(func(t types.T, u types.U) types.R {
    return t.(int) + u.(int)
})
fmt.Println(sum)

```



### Installation
```go get github.com/lovermaker/stream```

### Guides

#### New Stream

There are to ways to new a stream:

* new a stream from a slice: `NewStream`

  ```go
  func main() {
      nums := make([]int, 0)
      for i := 0; i < 100; i++ {
          nums = append(nums, i)
      }
      count := NewStream(nums).Filter(func(t types.T) bool {
          return t.(int) % 3 == 0
      }).Count()
      fmt.Println(count)
  }
  ```

* new a stream from a channel: `NewStreamFromChannel`

  ```go
  func produce(students []int) chan interface{} {
  	out := make(chan interface{})
  	go func() {
  		for _, s := range students {
  			out <- s
  		}
  		close(out)
  	}()
  	return out
  }
  
  func main() {
      nums := make([]int, 0)
      for i := 0; i < 100; i++ {
          nums = append(nums, i)
      }
     count := NewStreamFromChannel(produce(data)).Parallel(5).Filter(func(t types.T) bool {
      	return t.(int) % 3 == 0
  	}).Count()
  }
  ```

  the difference between `NewStream` and `NewStreamFromChannel` is that ``NewStreamFromChannel` `does not use reflection

#### Supported Operations

| Operations | Description                                                  |
| ---------- | ------------------------------------------------------------ |
| Filter     | filters elements by predicate                                |
| Map        | applies the given function to the elements                   |
| Peek       | performs the provided action on each element                 |
| Limit      | truncate size to be no longer than maxSize in length         |
| Skip       | discards the first n elements of th stream                   |
| Distinct   | returns streams consisting of distinct elements              |
| Sorted     | sort elements in order with comparator                       |
| Max        | returns max element with comparator                          |
| Min        | returns min element with comparator                          |
| Reduce     | performs a reduction  on the elements of this streamm, using an associative accumulation function, and returns an reduced value |
| Count      | the count of element of  the stream                          |
| FindFirst  | returns the first element of this stream or nil if the stream is empty |
| AnyMatch   | returns whether any elements of this stream match the provided predicate. returns `true` if any elements of the stream match the provided predicate, otherwise `false` |
| AllMatch   | returns whether all elements of this stream match the provided predicate. returns `true` if either all elements of the stream match the provided predicate or the stream is empty, otherwise `false` |
| NonMatch   | returns whether no elements of this stream match the provided predicate. returns true if either no elements of the stream match the provided predicate or the stream is empty, otherwise `false` |
| ForEach    | performs an action for each element of this stream           |
| Collect    | performs a mutable reduction operation on the elements of this stream using a collector |

#### Examples


