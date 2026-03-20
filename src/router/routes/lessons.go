package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var lessonRoutes = []Routes{
	{
		URI:      "/lessons",
		Method:   http.MethodPost,
		Function: controllers.CreateLesson,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons",
		Method:   http.MethodGet,
		Function: controllers.FetchAllLessons,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons/{lessonID}",
		Method:   http.MethodGet,
		Function: controllers.FetchLesson,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons/{lessonID}",
		Method:   http.MethodPut,
		Function: controllers.UpdateLesson,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons/{lessonID}",
		Method:   http.MethodDelete,
		Function: controllers.DeleteLesson,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons/{lessonID}/status",
		Method:   http.MethodPatch,
		Function: controllers.UpdateLessonStatus,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons/{lessonID}/students",
		Method:   http.MethodGet,
		Function: controllers.FetchLessonStudents,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons/{lessonID}/students/{studentID}",
		Method:   http.MethodPost,
		Function: controllers.AddStudentToLesson,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/lessons/{lessonID}/students/{studentID}",
		Method:   http.MethodDelete,
		Function: controllers.RemoveStudentFromLesson,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/classes/{classID}/lessons",
		Method:   http.MethodGet,
		Function: controllers.FetchLessonsByClass,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
}
