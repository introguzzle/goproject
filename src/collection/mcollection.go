package collection

import (
	"goproject/src/optional"
	"sync"
)

type SyncCollection[T any] struct {
	Items []T
	mutex sync.Mutex
}

func NewSyncCollection[T any](items ...T) *SyncCollection[T] {
	c := &SyncCollection[T]{}
	c.Items = items
	c.mutex = sync.Mutex{}

	return c
}

func (c *SyncCollection[T]) Map(mapper Function[T, any]) *SyncCollection[any] {
	result := &SyncCollection[any]{}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, v := range c.Items {
		result.Items = append(result.Items, mapper(v))
	}

	return result
}

func (c *SyncCollection[T]) Each(consumer Consumer[T]) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, v := range c.Items {
		consumer(v)
	}
}

func (c *SyncCollection[T]) Empty() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Size() == 0
}

func (c *SyncCollection[T]) Filter(filter Predicate[T]) *SyncCollection[T] {
	result := &SyncCollection[T]{}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, v := range c.Items {
		if filter(v) {
			result.Items = append(result.Items, v)
		}
	}

	return result
}

func (c *SyncCollection[T]) Get(index int) *optional.Optional[T] {
	if (index < 0) || (index > c.Size()) {
		return optional.OfNil[T]()
	}

	return optional.Of(c.Items[index])
}

func (c *SyncCollection[T]) Set(index int, value T) *SyncCollection[T] {
	c.Items[index] = value
	return c
}

func (c *SyncCollection[T]) Size() int {
	return len(c.Items)
}

func (c *SyncCollection[T]) Append(values ...T) *SyncCollection[T] {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Items = append(c.Items, values...)
	return c
}

func (c *SyncCollection[T]) Pop() T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v := c.Items[c.Size()-1]
	c.Items = c.Items[:len(c.Items)-1]
	return v
}

func (c *SyncCollection[T]) Clear() *SyncCollection[T] {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Items = []T{}
	return c
}
