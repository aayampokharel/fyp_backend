package filehandling

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.FileRouteWrapper {
	var prefix = "/certificate"

	var routes []common.FileRouteWrapper = []common.FileRouteWrapper{
		// /certificate/preview?id=123
		{
			Mux:             mux,
			Prefix:          prefix,
			Route:           "/preview",
			Method:          enum.METHODGET,
			RequestDataType: nil,
			URLQueries:      GetHTMLRequestQuery,
			GetORDeleteHandler: func(i map[string]string) entity.FileResponse {
				return module.Controller.HandleGetHTMLFile(i)
			},
		},
		// /certificate/download?file_id=123&category_id=123&is_download_all=true&category_name=abc
		{
			Mux:             mux,
			Prefix:          prefix,
			Route:           "/download",
			Method:          enum.METHODGET,
			RequestDataType: nil,
			URLQueries:      GetPDFFileInListQuery,
			GetORDeleteHandler: func(i map[string]string) entity.FileResponse {
				return module.Controller.HandleGetPDFFileInList(i)
			},
		},
	}
	return routes
}
