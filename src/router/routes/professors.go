package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var professorsRoutes = []Routes{
	{
		URI:      "/users/professors",
		Method:   http.MethodPost,
		Function: controllers.CreateProfessor,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/professors",
		Method:   http.MethodGet,
		Function: controllers.FetchProfessors,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/professors/all",
		Method:   http.MethodGet,
		Function: controllers.ReturnAllProfessors,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/professors/classes-count",
		Method:   http.MethodPost,
		Function: controllers.CountProfessorClasses,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/professors/{professorID}/classes/{classID}",
		Method:   http.MethodPatch,
		Function: controllers.AssignProfessorToClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/professor/me",
		Method:   http.MethodGet,
		Function: controllers.FetchCurrentProfessor,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermProfessor,
		},
	},
	{
		URI:      "/professor/classes",
		Method:   http.MethodGet,
		Function: controllers.FetchCurrentProfessorClasses,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermProfessor,
		},
	},
	{
		URI:      "/professor/lessons",
		Method:   http.MethodGet,
		Function: controllers.FetchCurrentProfessorLessons,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermProfessor,
		},
	},
	{
		URI:      "/professor/students",
		Method:   http.MethodGet,
		Function: controllers.FetchCurrentProfessorStudents,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermProfessor,
		},
	},
}
