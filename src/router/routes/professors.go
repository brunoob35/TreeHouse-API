package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var professorsRoutes = []Routes{
	{
		URI:      "/users/professor",
		Method:   http.MethodPost,
		Function: controllers.CreateProfessor,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/users/professors",
		Method:   http.MethodGet,
		Function: controllers.FetchProfessors,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
		},
	},
	{
		URI:      "/users/professors/all",
		Method:   http.MethodGet,
		Function: controllers.ReturnAllProfessors,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
		},
	},
}
