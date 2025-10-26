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
			Mux:             mux,
			Prefix:          prefix,
			Route:           "/login",
			Method:          enum.METHODPOST,
			RequestDataType: AdminLoginRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleAdminLogin(i.(AdminLoginRequest))
			},
		},
	}

	return routes
}
