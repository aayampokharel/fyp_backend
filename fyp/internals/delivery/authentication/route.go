package authentication

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) {
	mux.HandleFunc("/new-institution", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		module.Controller.HandleCreateNewInstitution(w, r)
	})
}
