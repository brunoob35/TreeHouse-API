package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var userRoutes = []Routes{
	{
		URI:      "/users",
		Method:   http.MethodPost,
		Function: controllers.CreateUser,
		Auth:     false,
	},
	{
		URI:      "/users/gestor",
		Method:   http.MethodPost,
		Function: controllers.CreateGestor,
		Auth:     false,
	},
	{
		URI:      "/users",
		Method:   http.MethodGet,
		Function: controllers.FetchUsers,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/users/active",
		Method:   http.MethodGet,
		Function: controllers.FetchActiveUsers,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
		},
	},
	{
		URI:      "/users/{userID}",
		Method:   http.MethodGet,
		Function: controllers.FetchUser,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/users/{userID}",
		Method:   http.MethodPut,
		Function: controllers.UpdateUser,
		Auth:     true,
	},
	{
		URI:      "/users/{userID}",
		Method:   http.MethodDelete,
		Function: controllers.DeleteUser,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
}
