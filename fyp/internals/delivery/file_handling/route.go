package filehandling

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.RouteWrapper {
	var prefix = "/certificate"

	var routes []common.RouteWrapper = []common.RouteWrapper{
		// /certificate/preview?id=123
		{
			Mux:             mux,
			Prefix:          prefix,
			Route:           "/preview",
			Method:          enum.METHODGET,
			RequestDataType: nil,
			URLQueries:      GetHTMLRequestQuery,
			ResponseType:    enum.HTML,
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.HandleGetHTMLFile(i.(map[string]string))
			},
		},
	}
	return routes
}
