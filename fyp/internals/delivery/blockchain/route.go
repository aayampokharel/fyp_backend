package delivery

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

// mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	blocks, err := module.Controller.InsertNewCertificateData()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(blocks)
// })

func RegisterRoutes(mux *http.ServeMux, module *Module) []common.RouteWrapper {
	var prefix = "/blockchain"

	var routes []common.RouteWrapper = []common.RouteWrapper{
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
	}

	return routes
}
