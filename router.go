package main

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
}

// NewRouter creates a new router instance with the provided middleware stack assigned
func NewRouter(middleware ...Middleware) Router {
	return Router{
		mux:        http.DefaultServeMux,
		middleware: middleware,
	}
}

// All registers a handler on all request methods on the provided uri
func (r Router) HandleFunc(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc(path, r.wrap(handler, middleware...))
}

func (r Router) Handle(path string, handler http.Handler, middleware ...Middleware) {
	r.mux.HandleFunc(path, r.wrap(handler.ServeHTTP, middleware...))
}

// wrap handler with middleware
func (r Router) wrap(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	stack := append(r.middleware, middleware...)

	for i := range stack {
		if stack[len(stack)-1-i] == nil {
			continue
		}

		handler = stack[len(stack)-1-i](handler)
	}

	return handler
}
