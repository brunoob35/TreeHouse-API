package routes

import (
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/controllers"
)

var reportsRoutes = []Routes{
	{
		URI:      "/reports/audit-logs",
		Method:   http.MethodGet,
		Function: controllers.FetchAuditLogs,
		Auth:     true,
		Permissions: []authentication.Permission{
			authentication.PermGestaoMaster,
		},
	},
}
