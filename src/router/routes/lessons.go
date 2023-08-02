package routes

import (
	"github.com/brunoob35/TreeHouse-API/src/controllers"
	"net/http"
)

var lessonsRoutes = []Routes{
	{
		URI:      "/lessons",
		Method:   http.MethodPost,
		Function: controllers.CreateLessonWStudent,
		Auth:     true,
	},
	{
		URI:      "/lessons",
		Method:   http.MethodGet,
		Function: controllers.FetchAllLessons,
		Auth:     true,
	},
	{
		URI:      "/lessons/teacher/{techerID}",
		Method:   http.MethodGet,
		Function: controllers.FetchLessonByTecherID,
		Auth:     true,
	},
	{
		URI:      "/lessons/class/{classID}",
		Method:   http.MethodGet,
		Function: controllers.FetchLessonByClassID,
		Auth:     true,
	},
	//{
	//	URI:      "/classes/{classID}",
	//	Method:   http.MethodPut,
	//	Function: controllers.UpdateClass,
	//	Auth:     true,
	//},
	//{
	//	URI:      "/lessons/{studentID}",
	//	Method:   http.MethodPut,
	//	Function: controllers.SetStudentToClass,
	//	Auth:     true,
	//},
	//{
	//	URI:      "/lessons/{classID}",
	//	Method:   http.MethodDelete,
	//	Function: controllers.RemoveStudentFromClass,
	//	Auth:     true,
	//},
}
