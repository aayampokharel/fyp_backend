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
		// /auth/new-institution
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/new-institution",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: CreateInstitutionRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleCreateNewInstitution(i.(CreateInstitutionRequest))
			},
		},

		// /auth/new-user
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/new-user",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: CreateUserAccountRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleCreateNewUserAccount(i.(CreateUserAccountRequest))
			},
		},
		// /auth/new-faculty
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/new-faculty",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: CreateFacultyRequest{},
			InnerFunc: func(i interface{}) entity.Response {

				institutionInfo, response := module.Controller.HandleCreateNewFaculty(i.(CreateFacultyRequest))
				if institutionInfo != nil {
					module.UseCase.Service.Logger.Debugln("Broadcast", *institutionInfo)
					module.SSEService.Broadcast(*institutionInfo)
				}
				return response
			},
		},
		// /auth/verify-institution/{id}
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/verify-institution",
			Method:                  enum.METHODGET,
			RequestDataTypeInstance: nil,
			URLQueries:              CheckInstitutionIsActiveQuery,
			InnerFunc: func(i interface{}) entity.Response {

				return module.Controller.HandleCheckInstitutionIsActive(i.(map[string]string))

			},
		},
		//POST /auth/institution/login
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/institution/login",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: InstitutionLoginRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleInstitutionsLogin(i.(InstitutionLoginRequest))

			},
		},
	}

	return routes
}
