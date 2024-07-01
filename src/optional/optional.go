package optional

import "errors"

type Optional[T any] struct {
	value  *T
	exists bool
}

func Of[T any](value T) *Optional[T] {
	return &Optional[T]{
		value:  &value,
		exists: true,
	}
}

func OfNil[T any]() *Optional[T] {
	return &Optional[T]{
		exists: false,
	}
}

func (o *Optional[T]) Get() T {
	if o.IsEmpty() {
		panic(errors.New("value is nil"))
	}

	return *o.value
}

func (o *Optional[T]) IsEmpty() bool {
	return !o.IsPresent()
}

func (o *Optional[T]) IsPresent() bool {
	return o.value != nil || o.exists
}

func (o *Optional[T]) OrElse(other T) T {
	if o.IsEmpty() {
		return other
	}

	return *o.value
}

func (o *Optional[T]) OrElsePanic() T {
	if o.IsEmpty() {
		panic(errors.New("value cannot be nil"))
	}
	return *o.value
}
