package collection

type Collection[T any] struct {
	Items []T `json:"items"`
}

func NewCollection[T any](items ...T) *Collection[T] {
	c := &Collection[T]{}
	c.Items = items

	return c
}

func (c *Collection[T]) Map(mapper Function[T, any]) *Collection[any] {
	result := &Collection[any]{}
	for _, v := range c.Items {
		result.Items = append(result.Items, mapper(v))
	}

	return result
}

func (c *Collection[T]) Empty() bool {
	return c.Size() == 0
}

func (c *Collection[T]) Each(consumer Consumer[T]) {
	for _, v := range c.Items {
		consumer(v)
	}
}

func (c *Collection[T]) Filter(filter Predicate[T]) *Collection[T] {
	result := &Collection[T]{}
	for _, v := range c.Items {
		if filter(v) {
			result.Items = append(result.Items, v)
		}
	}

	return result
}

func (c *Collection[T]) Get(index int) T {
	return c.Items[index]
}

func (c *Collection[T]) Set(index int, value T) *Collection[T] {
	c.Items[index] = value
	return c
}

func (c *Collection[T]) Size() int {
	return len(c.Items)
}

func (c *Collection[T]) Append(values ...T) *Collection[T] {
	c.Items = append(c.Items, values...)
	return c
}

func (c *Collection[T]) Pop() T {
	v := c.Items[c.Size()-1]
	c.Items = c.Items[:len(c.Items)-1]
	return v
}

func (c *Collection[T]) Clear() *Collection[T] {
	c.Items = []T{}
	return c
}
