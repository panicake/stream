package types

type Comparator interface {
	Compare(a interface{}, b interface{}) int
}

type ComparatorFunc func(a interface{}, b interface{}) int

func(c ComparatorFunc)Compare(a interface{}, b interface{}) int {
	return c(a, b)
}


