package types

type Cal[T any] interface {
	Calculate(a ...T) T
}
