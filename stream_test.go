package go_stream

import (
	"fmt"
	"github.com/lovermaker/stream/types"
	"testing"
)

type student struct {
	ID   int
	Name string
}

// buildStudent build array data
func buildStudent(max int) []*student {
	var data []*student
	for i := 0; i < max; i++ {
		data = append(data, &student{ ID: i, Name: fmt.Sprintf("Name-%d", i)})
	}
	return data
}

// produceStudentChannel build array data
func produceStudentChannel(students []*student) chan interface{} {
	out := make(chan interface{})
	go func() {
		for _, s := range students {
			out <- s
		}
		close(out)
	}()
	return out
}

func getNum(num int) int {
	flag := false
	for i := 2; i < num; i++ {
		if num / i == 0 {
			flag = true
			break
		}
	}
	if flag {
		return num
	}
	return 0
}

func TestNewStream(t *testing.T) {
	data := buildStudent(100000)
	maxId := NewStream(data).Parallel(4).Map(func(t types.T) types.R {
		return getNum(t.(*student).ID)
	}).Max(func(a interface{}, b interface{}) int {
		return a.(int) - b.(int)
	})
	t.Log(maxId)
}

func TestNewStreamFromChannel(t *testing.T) {
	data := buildStudent(100)
	count := NewStreamFromChannel(produceStudentChannel(data)).Parallel(5).Filter(func(t types.T) bool {
		if t.(*student).ID > 15 && t.(*student).ID < 55 {
			return true
		}
		return false
	}).Map(func(t types.T) types.R {
		return t.(*student).ID
	}).Count()
	t.Log(count)
}

func find(num int) int {
	sum := 0
	for i := num; i > 0; i-- {
		if num % i == 0 {
			sum += i
		}
	}
	return sum
}

func TestNewStreamExample(t *testing.T) {
	nums := make([]int, 0)
	for i := 0; i < 100; i++ {
		nums = append(nums, i)
	}
	sum := NewStream(nums).Parallel(10).Filter(func(t types.T) bool {
		return t.(int) % 3 == 0
	}).Map(func(t types.T) types.R {
		num := t.(int)
		return find(num)
	}).Reduce(func(t types.T, u types.U) types.R {
		return t.(int) + u.(int)
	})
	t.Log(sum)
}

func BenchmarkNewStreamExample2(b *testing.B) {
	nums := make([]int, 0)
	for i := 0; i < b.N; i++ {
		nums = append(nums, i)
	}
	sum := 0
	for _, value := range nums {
		if value % 3 == 0 {
			sum += find(value) / (value + 1)
		}
	}
	b.Log(sum)
}
