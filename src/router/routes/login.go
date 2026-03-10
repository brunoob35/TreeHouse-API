package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var loginRoutes = []Routes{
	{
		URI:      "/login",
		Method:   http.MethodPost,
		Function: controllers.Login,
		Auth:     false,
	},
}
