package server

import (
	_ "github.com/gorilla/mux"
	"goproject/src/domain"
	"goproject/src/http/controller"
	"goproject/src/persistence"
)

func loadRoutes() {
	loadProductEndpoints()
}

func loadAuthEndpoints() {
	endpoint := &Path{
		Value: "/api/v1/auth",
	}

	c := controller.NewSecurityController[*domain.User](
		persistence.Connect(),
		persistence.Factory[*domain.User]{
			Create: func() *domain.User {
				return &domain.User{}
			},
		},
	)
}

func loadProductEndpoints() {
	endpoint := &Path{
		Value: "/api/v1/products",
	}

	c := controller.NewProductController(persistence.NewManager[*domain.Product](
		persistence.Connect(),
		persistence.Factory[*domain.Product]{
			Create: func() *domain.Product {
				return &domain.Product{}
			},
		},
	))

	Get(endpoint, c.GetAll)
	Get(endpoint.AddParameter(), c.Get)
	Post(endpoint, c.Create)
	Put(endpoint.AddParameter(), c.Put)
	Patch(endpoint.AddParameter(), c.Patch)
	Delete(endpoint.AddParameter(), c.Delete)
}
