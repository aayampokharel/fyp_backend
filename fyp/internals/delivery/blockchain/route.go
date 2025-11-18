package delivery

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.RouteWrapper {
	var prefix = "/blockchain"

	var routes []common.RouteWrapper = []common.RouteWrapper{
		// POST /blockchain/certificates
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/certificates",
			Method:                  enum.METHODPOST,
			RequestDataTypeInstance: CreateCertificateDataRequest{},
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.InsertNewCertificateData(i.(CreateCertificateDataRequest))
			},
		},
		// GET /blockchain/certificates/fake?random_id=123
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/certificates/fake",
			Method:                  enum.METHODGET,
			RequestDataTypeInstance: nil,
			URLQueries:              GetAllPendingInstitutionsQuery,
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.InsertFakeCertificateData(i.(map[string]string))
			},
		},
		// GET /blockchain/certificate-batch?institution_id=_____&institution_faculty_id=_____&category_id=_____
		{
			Mux:                     mux,
			Prefix:                  prefix,
			Route:                   "/certificate-batch",
			Method:                  enum.METHODGET,
			RequestDataTypeInstance: nil,
			URLQueries:              GetCertificateDataListRequestQuery,
			InnerFunc: func(i interface{}) entity.Response {
				return module.Controller.GetCertificateDataList(i.(map[string]string))
			},
		},
	}

	return routes
}
