package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var customersRoutes = []Routes{
	{
		URI:      "/customers",
		Method:   http.MethodPost,
		Function: controllers.CreateCustomer,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/customers",
		Method:   http.MethodGet,
		Function: controllers.FetchCustomers,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/customers/{customerID}",
		Method:   http.MethodGet,
		Function: controllers.FetchCustomer,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/customers/{customerID}",
		Method:   http.MethodPut,
		Function: controllers.UpdateCustomer,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/customers/{customerID}",
		Method:   http.MethodDelete,
		Function: controllers.DeleteCustomer,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/customers/{customerID}/students",
		Method:   http.MethodGet,
		Function: controllers.FetchCustomerStudents,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/customers/{customerID}/addresses",
		Method:   http.MethodGet,
		Function: controllers.FetchCustomerAddresses,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
}
