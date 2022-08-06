package stream

import (
	"math/rand"
	"sort"
	"sync"
	"testing"
)

func TestStreamFilter(t *testing.T) {
	var expected []int
	num := 500
	for i := 0; i < num; i++ {
		if i%3 == 0 {
			expected = append(expected, i)
		}
	}
	testFunc := func() {
		threads := rand.Intn(50)
		var result []int
		var mu sync.Mutex
		NewStream(func(in chan interface{}) {
			for i := 0; i < num; i++ {
				in <- i
			}
		}).Parallel(threads).Filter(func(i interface{}) bool {
			return i.(int)%3 == 0
		}).ForEach(func(i interface{}) {
			mu.Lock()
			result = append(result, i.(int))
			mu.Unlock()
		})
		sort.Ints(result)
		for index, value := range expected {
			if result[index] != value {
				t.Fatalf("expected value, %d, result value: %d", value, result[index])
			}
		}
	}
	for i := 0; i < 10; i++ {
		testFunc()
	}

}

func TestStreamCountAndSum(t *testing.T) {
	expectedSum := 0
	expectedCount := 0
	for i := 0; i < 100; i++ {
		if i%3 == 0 {
			expectedSum += i
			expectedCount += 1
		}
	}
	baseSteam := NewStream(func(in chan interface{}) {
		for i := 0; i < 100; i++ {
			in <- i
		}
	}).Filter(func(i interface{}) bool {
		return i.(int)%3 == 0
	})
	count := baseSteam.Count()
	sum := baseSteam.Reduce(func(i1, i2 interface{}) interface{} {
		return i1.(int) + i2.(int)
	}).(int)
	if count != int64(expectedCount) {
		t.Fatalf("expected value, %d, result value: %d", expectedCount, count)
	}
	if sum != expectedSum {
		t.Fatalf("expected value, %d, result value: %d", expectedSum, sum)
	}
}

func TestStreamDistinct(t *testing.T) {
	num := 500
	var expected []int
	var list []int
	for i := 0; i < num; i++ {
		expected = append(expected, i)
		len := rand.Intn(10) + 1
		for j := 0; j < len; j++ {
			list = append(list, i)
		}
	}
	var result []int
	var mu sync.Mutex
	NewStream(func(in chan interface{}) {
		for _, value := range list {
			in <- value
		}
	}).Parallel(50).Distinct().ForEach(func(i interface{}) {
		mu.Lock()
		result = append(result, i.(int))
		mu.Unlock()
	})
	sort.Ints(result)
	for index, value := range expected {
		if result[index] != value {
			t.Fatalf("expected value, %d, result value: %d", value, result[index])
		}
	}
}
