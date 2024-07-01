package repository

import (
	"goproject/src/collection"
	"goproject/src/domain"
	"goproject/src/optional"
	"goproject/src/persistence"
)

type Repository[T domain.Entity] struct {
	Manager *persistence.Manager[T]
}

func (r *Repository[T]) Find(id int64) *optional.Optional[T] {
	return r.Manager.Find(id)
}

func (r *Repository[T]) FindAll() *collection.Collection[T] {
	return r.Manager.FindAll()
}

func (r *Repository[T]) Count() int64 {
	return r.Manager.Count()
}

func (r *Repository[T]) Exists(criteria string) bool {
	return r.Manager.Exists(criteria)
}
