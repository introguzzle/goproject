package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"goproject/src/domain"
	"goproject/src/http/response"
	"goproject/src/persistence"
	"log"
	"net/http"
	"strconv"
)

type Controller[T domain.Entity] struct {
	Manager *persistence.Manager[T]
}

func NewController[T domain.Entity](
	manager *persistence.Manager[T],
) *Controller[T] {
	return &Controller[T]{
		Manager: manager,
	}
}

func (c *Controller[T]) IntegerKey(
	request *http.Request,
	key string,
) (int64, error) {
	vars := mux.Vars(request)
	value, err := strconv.Atoi(vars[key])
	return int64(value), err
}

func (c *Controller[T]) GetAll(request *http.Request) response.Response {
	models := c.Manager.FindAll()
	return response.WrapData(models, http.StatusOK)
}

func (c *Controller[T]) Get(request *http.Request) response.Response {
	id, err := c.IntegerKey(request, "id")
	if err != nil {
		return response.WrapError(err, http.StatusBadRequest)
	}

	optional := c.Manager.Find(id)
	if optional.IsEmpty() {
		return response.WrapErrorString("Not found", http.StatusNotFound)
	}

	return response.WrapData(optional.Get(), http.StatusOK)
}

func (c *Controller[T]) Create(request *http.Request) response.Response {
	return c.Update(c.Manager.Factory.Create(), request)
}

func (c *Controller[T]) Put(request *http.Request) response.Response {
	id, err := c.IntegerKey(request, "id")
	if err != nil {
		return response.BadRequest()
	}

	if c.Manager.Find(id).IsEmpty() {
		return response.NotFound()
	}

	entity := c.Manager.Factory.Create()
	entity.SetId(id)

	return c.Update(entity, request)
}

func (c *Controller[T]) Patch(request *http.Request) response.Response {
	id, err := c.IntegerKey(request, "id")

	if err != nil {
		return response.BadRequest()
	}

	if c.Manager.Find(id).IsEmpty() {
		return response.NotFound()
	}

	entity := c.Manager.Factory.Create()
	entity.SetId(id)

	return c.Update(entity, request)
}

func (c *Controller[T]) Update(entity T, request *http.Request) response.Response {
	err := json.NewDecoder(request.Body).Decode(&entity)
	if err != nil {
		return response.WrapErrorString(err.Error(), http.StatusBadRequest)
	}

	c.Manager.Persist(entity)
	err = c.Manager.Flush()

	log.Printf("c.Manager.Flush returned: %v \n\n", err)

	if err != nil {
		return response.Internal()
	}

	return response.WrapData(entity, http.StatusOK)
}

func (c *Controller[T]) Delete(request *http.Request) response.Response {
	id, err := c.IntegerKey(request, "id")
	if err != nil {
		return response.WrapError(err, http.StatusBadRequest)
	}

	if !c.Manager.Exists(fmt.Sprintf("id = %d", id)) {
		return response.WrapErrorString("Not found", http.StatusNotFound)
	}

	c.Manager.Delete(c.Manager.Find(id).Get())
	return response.Ok()
}
