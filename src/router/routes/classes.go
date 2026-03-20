package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var classesRoutes = []Routes{
	{
		URI:      "/classes",
		Method:   http.MethodPost,
		Function: controllers.CreateClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes",
		Method:   http.MethodGet,
		Function: controllers.FetchAllActiveClasses,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/all",
		Method:   http.MethodGet,
		Function: controllers.FetchAllClasses,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/{classID}",
		Method:   http.MethodGet,
		Function: controllers.FetchClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/{classID}",
		Method:   http.MethodPut,
		Function: controllers.UpdateClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/{classID}",
		Method:   http.MethodDelete,
		Function: controllers.DeleteClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/{classID}/students",
		Method:   http.MethodGet,
		Function: controllers.FetchClassStudents,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/{classID}/students/{studentID}",
		Method:   http.MethodPost,
		Function: controllers.AddStudentToClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/{classID}/students/{studentID}",
		Method:   http.MethodDelete,
		Function: controllers.RemoveStudentFromClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/private",
		Method:   http.MethodPost,
		Function: controllers.CreatePrivateClassFromStudent,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
}
