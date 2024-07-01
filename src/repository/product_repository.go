package repository

import (
	"goproject/src/domain"
	"goproject/src/persistence"
)

func NewProductRepository() *Repository[*domain.Product] {
	manager := &persistence.Manager[*domain.Product]{}

	return manager.GetRepository()
}
