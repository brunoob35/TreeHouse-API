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
	{
		URI:      "/auth/forgot-password",
		Method:   http.MethodPost,
		Function: controllers.ForgotPassword,
		Auth:     false,
	},
	{
		URI:      "/auth/reset-password",
		Method:   http.MethodPost,
		Function: controllers.ResetPassword,
		Auth:     false,
	},
}
