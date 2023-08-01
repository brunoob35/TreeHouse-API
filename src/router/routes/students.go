package routes

import (
	"github.com/brunoob35/TreeHouse-API/src/controllers"
	"net/http"
)

var studentsRoutes = []Routes{
	{
		URI:      "/students",
		Method:   http.MethodPost,
		Function: controllers.CreateStudent,
		Auth:     true,
	},
	{
		URI:      "/students",
		Method:   http.MethodGet,
		Function: controllers.FetchAllStudents,
		Auth:     true,
	},
	{
		URI:      "/students/active",
		Method:   http.MethodGet,
		Function: controllers.FetchActiveStudents,
		Auth:     true,
	},
	{
		URI:      "/students/{studentID}",
		Method:   http.MethodGet,
		Function: controllers.FetchStudentByID,
		Auth:     true,
	},
	{
		URI:      "/students/{studentID}",
		Method:   http.MethodPut,
		Function: controllers.UpdateStudent,
		Auth:     true,
	},
	//{
	//	URI:      "/students/{studentID}",
	//	Method:   http.MethodDelete,
	//	Function: controllers.DeleteStudent,
	//	Auth:     true,
	//},
}
