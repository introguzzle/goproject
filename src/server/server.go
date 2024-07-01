package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"goproject/src/http/response"
	"goproject/src/server/middleware"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"unsafe"
)

var R = &Router{
	R: mux.NewRouter(),
}

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	PATCH  Method = "PATCH"
)

type Handler func(r *http.Request) response.Response
type Method string
type Path struct {
	Value string
}

func (p *Path) AddParameter() *Path {
	return &Path{
		Value: p.Value + "/{id:[0-9]+}",
	}
}

type Router struct {
	R           *mux.Router
	Middlewares []middleware.Middleware
}

func (r *Router) HandleFunc(p *Path, m Method, h Handler) *Router {
	r.R.HandleFunc(p.Value, wrapHandler(h)).Methods(string(m))
	fmt.Printf("Registered route: %s %s\n", m, p.Value)
	return r
}

func wrapHandler(h Handler) func(http.ResponseWriter, *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {
		resp := h(r)
		data, err := resp.Serialize()

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		status := resp.Status()

		w.WriteHeader(status)
		w.Header().Set(middleware.STATUS, strconv.FormatUint(uint64(status), 10))

		if data != nil {
			_, _ = w.Write(data)
		}
	}

	for i := len(R.Middlewares) - 1; i >= 0; i-- {
		m := R.Middlewares[i]
		f = m.Handle(f)
	}

	return f
}

func loadMiddlewares() {
	var middlewares []middleware.Middleware
	middlewares = append(middlewares, &middleware.RestMiddleware{})
	middlewares = append(middlewares, &middleware.LoggingMiddleware{})
	middlewares = append(middlewares, &middleware.AuthenticationMiddleware{})

	R.Middlewares = middlewares
}

func setupRouter() {
	loadMiddlewares()
	loadRoutes()
}

func printRoutes(router *mux.Router) {
	v := reflect.ValueOf(router).Elem()
	f := v.FieldByName("routes")
	f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()

	for i := 0; i < f.Len(); i++ {
		route := f.Index(i).Elem()
		routeValue := reflect.Indirect(route)

		pathTemplate := routeValue.FieldByName("name")
		fmt.Println(pathTemplate.String())
	}
}

func Serve(ctx context.Context, statusChan chan<- string) {
	setupRouter()
	printRoutes(R.R)

	server := &http.Server{
		Addr:    ":8080",
		Handler: R.R,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			statusChan <- fmt.Sprintf("Server error: %v", err)
		}
	}()

	statusChan <- "Server started"

	<-ctx.Done()
	log.Println("Shutting down the server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}
