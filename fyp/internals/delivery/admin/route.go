package admin

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.RouteWrapper {
	var prefix = "/admin"

	var routes []common.RouteWrapper = []common.RouteWrapper{
		// /admin/login
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/login",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: AdminLoginRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleAdminLogin(i.(AdminLoginRequest))
			},
		},

		// /admin/pending-institutions?admin_id=12345
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/pending-institutions",
			Method:                  enum.METHODGET,
			URLQueries:              GetAllPendingInstitutionsQuery,
			RequestDataTypeInstance: nil,
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleGetPendingInstitutionList(i.(map[string]string))
			},
		},
	}

	return routes
}
