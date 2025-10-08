package delivery

import (
	"encoding/json"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, module *Module) {
	mux.HandleFunc("/certificates", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		blocks, err := module.Controller.InsertNewCertificateData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blocks)
	})
}
