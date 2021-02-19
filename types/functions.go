package types

type T = interface{}
type U = interface{}
type R = interface{}

type Consumer interface {
	Accept(t T)
}

type ConsumerFunc func(T)

func (f ConsumerFunc) Accept(t T) {
	f(t)
}

type Function = func(T) R

type BiFunction interface {
	Apply(T, U) R
}

type BiFunctionFunc func(T, U) R

func (b BiFunctionFunc) Apply(t1 T, t2 U) R {
	return b(t1, t2)
}

type Predicate = func(T) bool
