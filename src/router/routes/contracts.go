package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var contractsRoutes = []Routes{
	{
		URI:      "/contracts/types",
		Method:   http.MethodGet,
		Function: controllers.FetchContractTypes,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/contracts/statuses",
		Method:   http.MethodGet,
		Function: controllers.FetchContractStatuses,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/contracts",
		Method:   http.MethodGet,
		Function: controllers.FetchContracts,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/contracts/{contractID}",
		Method:   http.MethodGet,
		Function: controllers.FetchContract,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/contracts",
		Method:   http.MethodPost,
		Function: controllers.CreateContract,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/contracts/{contractID}",
		Method:   http.MethodPut,
		Function: controllers.UpdateContract,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
	{
		URI:      "/contracts/{contractID}/class",
		Method:   http.MethodPost,
		Function: controllers.CreateClassFromContract,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestao,
			authentication.PermGestaoMaster,
		},
	},
}
