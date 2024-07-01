package collection

type Function[T any, R any] func(value T) R
type Consumer[T any] func(value T)
type Predicate[T any] func(value T) bool
type Supplier[T any] func() T
type Runner func()
