package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/middlewares"
	"github.com/gorilla/mux"
)

type Routes struct {
	URI         string
	Method      string
	Function    func(http.ResponseWriter, *http.Request)
	Auth        bool
	Permissions []authentication.Permission
	RequireAll  bool
}

//	How to declarre a new route
//	New public route
//	var loginRoutes = []Routes{
//		{
//			URI:      "/login",
//			Method:   http.MethodPost,
//			Function: controllers.Login,
//			Auth:     false,
//		},
//	}
//	Route with specific permissions
//	{
//	URI:      "/aulas",
//	Method:   http.MethodGet,
//	Function: controllers.ListarAulas,
//	Auth:     true,
//	Permissions: []authentication.Permission{
//		authentication.PermGestao,
//		authentication.PermProfessor,
//		},
//	}

// Config adds all routes to the router.
//
// Route protection flow:
//   - public routes: no authentication, no authorization
//   - authenticated routes: Authentication middleware only
//   - authenticated routes with permissions:
//   - Authentication + Authorize
//   - or Authentication + AuthorizeAll
func Config(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, loginRoutes...)
	// routes = append(routes, studentsRoutes...)
	// routes = append(routes, classesRoutes...)
	// routes = append(routes, lessonsRoutes...)

	for _, route := range routes {
		handler := route.Function

		// Authorization only makes sense for authenticated routes.
		if route.Auth {
			handler = middlewares.Authentication(handler)

			if len(route.Permissions) > 0 {
				if route.RequireAll {
					handler = middlewares.AuthorizeAll(route.Permissions...)(handler)
				} else {
					handler = middlewares.Authorize(route.Permissions...)(handler)
				}
			}

			handler = middlewares.Logger(handler)
		} else {
			handler = middlewares.Logger(handler)
		}

		r.HandleFunc(route.URI, handler).Methods(route.Method, http.MethodOptions)
	}

	return r
}
