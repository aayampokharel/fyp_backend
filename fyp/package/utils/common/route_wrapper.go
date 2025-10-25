package common

import (
	"encoding/json"
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	err "project/package/errors"
	"reflect"
)

type RouteWrapper struct {
	Mux             *http.ServeMux
	Prefix          string
	Route           string
	Method          enum.HTTPMETHOD
	RequestDataType interface{}
	InnerFunc       func(interface{}) entity.Response
}

func NewRouteWrapper(routeInfos ...RouteWrapper) {
	for index := range routeInfos {
		routeInfo := routeInfos[index]
		routeInfo.Mux.HandleFunc(routeInfo.Prefix+routeInfo.Route, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
			(w).Header().Set("Content-Type", "application/json")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != routeInfo.Method.ToString() {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			reqType := reflect.TypeOf(routeInfo.RequestDataType)
			reqValue := reflect.New(reqType).Interface()

			if er := json.NewDecoder(r.Body).Decode(reqValue); er != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(HandleErrorResponse(400, err.ErrDecodingJSONString, er))
				return
			}

			returnFinalResponse := routeInfo.InnerFunc(reflect.ValueOf(reqValue).Elem().Interface())

			(w).WriteHeader(returnFinalResponse.Code)
			json.NewEncoder(w).Encode(returnFinalResponse)

		})

	}

}
