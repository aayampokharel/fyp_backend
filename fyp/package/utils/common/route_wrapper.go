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
	URLQueries      map[string]string
	InnerFunc       func(interface{}) entity.Response
	ResponseType    enum.RESPONSETYPE
}

type SSERouteWrapper struct {
	Mux       *http.ServeMux
	Prefix    string
	Route     string
	Method    enum.HTTPMETHOD
	InnerFunc func(w http.ResponseWriter, r *http.Request)
}

func NewSSERouteWrapper(sseRouteWrapper SSERouteWrapper) {
	sseRouteWrapper.Mux.HandleFunc(sseRouteWrapper.Prefix+sseRouteWrapper.Route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		(w).Header().Set("Content-Type", "text/event-stream")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if sseRouteWrapper.Method != enum.METHODGET || r.Method != sseRouteWrapper.Method.ToString() {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		sseRouteWrapper.InnerFunc(w, r)
	})
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
			var returnFinalResponse entity.Response

			if routeInfo.RequestDataType == nil && routeInfo.Method == enum.METHODGET {
				// var queryParams map[string]string
				for key, _ := range routeInfo.URLQueries {
					routeInfo.URLQueries[key] = r.URL.Query().Get(key)
				}

				returnFinalResponse = routeInfo.InnerFunc(routeInfo.URLQueries)

			} else {
				reqType := reflect.TypeOf(routeInfo.RequestDataType)
				reqValue := reflect.New(reqType).Interface()

				if er := json.NewDecoder(r.Body).Decode(reqValue); er != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(HandleErrorResponse(400, err.ErrDecodingJSONString, er))
					return
				}

				returnFinalResponse = routeInfo.InnerFunc(reflect.ValueOf(reqValue).Elem().Interface())
			}
			if routeInfo.ResponseType == enum.HTML {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if returnFinalResponse.Data != nil {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(returnFinalResponse.Data.(string)))
					return
				} else if returnFinalResponse.Data == nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(returnFinalResponse.Message))
					return
				}

			}
			(w).WriteHeader(returnFinalResponse.Code)
			json.NewEncoder(w).Encode(returnFinalResponse)

		})

	}

}
