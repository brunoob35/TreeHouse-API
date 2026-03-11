package router

import (
	"github.com/brunoob35/TreeHouse-API/src/middlewares"
	"github.com/brunoob35/TreeHouse-API/src/router/routes"
	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()
	r.Use(middlewares.CORS)
	return routes.Config(r)
}
