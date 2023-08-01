package routes

import (
	"github.com/brunoob35/TreeHouse-API/src/controllers"
	"net/http"
)

var iotRoutes = []Routes{
	{
		URI:      "/classes",
		Method:   http.MethodPost,
		Function: controllers.CreateClass,
		Auth:     true,
	},
	{
		URI:      "/classes",
		Method:   http.MethodGet,
		Function: controllers.FetchAllClasses,
		Auth:     true,
	},
	{
		URI:      "/classes/active",
		Method:   http.MethodGet,
		Function: controllers.FetchActiveClasses,
		Auth:     true,
	},
	{
		URI:      "/classes/{classID}",
		Method:   http.MethodGet,
		Function: controllers.FetchClassByID,
		Auth:     true,
	},
	//{
	//	URI:      "/classes/{classID}",
	//	Method:   http.MethodPut,
	//	Function: controllers.UpdateClass,
	//	Auth:     true,
	//},
	{
		URI:      "/classes/{classID}",
		Method:   http.MethodPut,
		Function: controllers.SetStudentToClass,
		Auth:     true,
	},
	{
		URI:      "/classes/{classID}",
		Method:   http.MethodDelete,
		Function: controllers.RemoveStudentFromClass,
		Auth:     true,
	},
}
