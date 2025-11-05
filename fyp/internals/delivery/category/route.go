package category

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.RouteWrapper {
	var prefix = "/institution"

	var routes []common.RouteWrapper = []common.RouteWrapper{
		//POST /institution/category
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/category",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: CreatePDFCategoryDto{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleCreatePDFCategory(i.(CreatePDFCategoryDto))
			},
		},
		//GET /institution/categories?institution_id=12345&institution_faculty_id=12345
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/categories",
			Method:                  enum.METHODGET,
			RequestDataTypeInstance: nil,
			URLQueries:              GetPDFCategoryRequestDtoQuery,
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleGetPDFCategoriesList(i.(map[string]string))
			},
		},
	}

	return routes
}
