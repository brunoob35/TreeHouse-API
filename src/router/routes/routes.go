package routes

import (
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
		r.HandleFunc(route.URI, route.Function).Methods(route.Method)
	}

	return r
}
