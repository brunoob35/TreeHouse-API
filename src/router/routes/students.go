package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var studentsRoutes = []Routes{
	{
		URI:      "/students",
		Method:   http.MethodPost,
		Function: controllers.CreateStudent,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermProfessor,
		},
	},
	{
		URI:      "/students",
		Method:   http.MethodGet,
		Function: controllers.FetchStudents,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermProfessor,
		},
	},
	{
		URI:      "/students/{studentID}",
		Method:   http.MethodGet,
		Function: controllers.FetchStudent,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermProfessor,
		},
	},
	{
		URI:      "/students/{studentID}",
		Method:   http.MethodPut,
		Function: controllers.UpdateStudent,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermProfessor,
		},
	},
	{
		URI:      "/students/{studentID}",
		Method:   http.MethodDelete,
		Function: controllers.DeleteStudent,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermProfessor,
		},
	},
}
