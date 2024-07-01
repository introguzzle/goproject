package controller

import (
	"goproject/src/domain"
	"goproject/src/persistence"
)

func NewProductController(
	manager *persistence.Manager[*domain.Product],
) *Controller[*domain.Product] {
	return &Controller[*domain.Product]{
		Manager: manager,
	}
}
