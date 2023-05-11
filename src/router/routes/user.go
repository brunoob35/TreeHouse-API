package routes

import (
	"github.com/brunoob35/TreeHouse-API/src/controllers"
	"net/http"
)

//Routes CRUD example
var userRoutes = []Routes{
	{
		URI: "/users",
		Method: http.MethodPost,
		Function: controllers.CreateUser,
		Auth: false,
	},
	{
		URI: "/users",
		Method: http.MethodGet,
		Function: controllers.FetchUsers,
		Auth: false,
	},
	{
		URI: "/users/{userID}",
		Method: http.MethodGet,
		Function: controllers.FetchUser,
		Auth: false,
	},
	{
		URI: "/users/{userID}",
		Method: http.MethodPut,
		Function: controllers.UpdateUser,
		Auth: false,
	},
	{
		URI: "/users/{userID}",
		Method: http.MethodDelete,
		Function: controllers.DeleteUser,
		Auth: false,
	},
}
