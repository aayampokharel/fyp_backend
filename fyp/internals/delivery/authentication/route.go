package authentication

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.RouteWrapper {
	var prefix = "/auth"

	var routes []common.RouteWrapper = []common.RouteWrapper{
		{
			Mux:             mux,
			Prefix:          prefix,
			Route:           "/new-institution",
			Method:          enum.METHODPOST,
			RequestDataType: CreateInstitutionRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleCreateNewInstitution(i.(CreateInstitutionRequest))
			},
		},

		{
			Mux:             mux,
			Prefix:          prefix,
			Route:           "/new-user",
			Method:          enum.METHODPOST,
			RequestDataType: CreateUserAccountRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleCreateNewUserAccount(i.(CreateUserAccountRequest))
			},
		},
	}

	return routes
}
