package routes

import (
	"github.com/brunoob35/TreeHouse-API/src/controllers"
	"net/http"
)

var loginRoutes = Routes{
	URI:      "/login",
	Method:   http.MethodPost,
	Function: controllers.Login,
	Auth:     false,
}
