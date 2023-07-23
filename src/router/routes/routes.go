package routes

import (
	"github.com/brunoob35/TreeHouse-API/src/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

type Routes struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
	Auth     bool
}

// Config adds all routes in the Router
func Config(r *mux.Router) *mux.Router {
	var routes []Routes
	routes = append(userRoutes, loginRoutes)

	for _, route := range routes {

		if route.Auth {
			r.HandleFunc(route.URI,
				middlewares.Logger(middlewares.Authentication(route.Function)),
			).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, route.Function).Methods(route.Method)
		}

		r.HandleFunc(route.URI, route.Function).Methods(route.Method)
	}

	return r
}
