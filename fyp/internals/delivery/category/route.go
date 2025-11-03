package category

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.RouteWrapper {
	var prefix = "/category"

	var routes []common.RouteWrapper = []common.RouteWrapper{
		//POST /category
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: CreatePDFCategoryDto{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleCreatePDFCategory(i.(CreatePDFCategoryDto))
			},
		},
	}

	return routes
}
