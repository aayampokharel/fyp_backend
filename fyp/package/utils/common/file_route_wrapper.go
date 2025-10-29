package common

import (
	"encoding/json"
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	err "project/package/errors"
	"reflect"
	"strconv"
)

type FileRouteWrapper struct {
	Mux    *http.ServeMux
	Prefix string
	Route  string
	Header http.Header

	Method           enum.HTTPMETHOD
	RequestDataType  interface{}
	GetDeleteHandler func(map[string]string) entity.FileResponse
	URLQueries       []string
	PostHandler      func(interface{}) entity.FileResponse
}

// if routeInfo.ResponseType == enum.HTML {
// 				w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 				if returnFinalFileResponse.Data != nil {
// 					w.WriteHeader(http.StatusOK)
// 					w.Write([]byte(returnFinalFileResponse.Data.(string)))
// 					return
// 				} else if returnFinalFileResponse.Data == nil {
// 					w.WriteHeader(http.StatusBadRequest)
// 					w.Write([]byte(returnFinalFileResponse.Message))
// 					return
// 				}

// 			} else if routeInfo.ResponseType == enum.PDFORZIP {
// 				if returnFinalFileResponse.Data != nil {
// 					w.WriteHeader(http.StatusOK)
// 					w.Write([]byte(returnFinalFileResponse.Data.(string)))
// 					return
// 				} else if returnFinalFileResponse.Data == nil {
// 					w.WriteHeader(http.StatusBadRequest)
// 					w.Write([]byte(returnFinalFileResponse.Message))
// 					return
// 				}
// 			}

func setFileHeaders(w http.ResponseWriter, fileData []byte, filename string, fileType enum.RESPONSETYPE) {
	switch fileType {
	case enum.PDF:
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	case enum.ZIP:
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	case enum.HTML:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
}

func NewFileRouteWrapper(routeInfos ...FileRouteWrapper) {
	for index := range routeInfos {
		routeInfo := routeInfos[index]
		routeInfo.Mux.HandleFunc(routeInfo.Prefix+routeInfo.Route, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
			// (w).Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != routeInfo.Method.ToString() {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			var returnFinalFileResponse entity.FileResponse

			if routeInfo.RequestDataType == nil && (routeInfo.Method == enum.METHODGET || routeInfo.Method == enum.METHODDELETE) {
				var queryParams map[string]string = make(map[string]string, 0)
				for _, val := range routeInfo.URLQueries {
					queryParams[val] = r.URL.Query().Get(val)
				}

				returnFinalFileResponse = routeInfo.GetDeleteHandler(queryParams)

			} else {
				reqType := reflect.TypeOf(routeInfo.RequestDataType)
				reqValue := reflect.New(reqType).Interface()

				if er := json.NewDecoder(r.Body).Decode(reqValue); er != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(HandleErrorResponse(400, err.ErrDecodingJSONString, er))
					return
				}

				returnFinalFileResponse = routeInfo.PostHandler(reflect.ValueOf(reqValue).Elem().Interface())
			}

			setFileHeaders(w, returnFinalFileResponse.Data, returnFinalFileResponse.FileName, returnFinalFileResponse.FileType)
			(w).WriteHeader(returnFinalFileResponse.Code)
			if returnFinalFileResponse.Data == nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(entity.Response{
					Code:    returnFinalFileResponse.Code,
					Message: returnFinalFileResponse.Message,
					Data:    nil,
				})
				return
			}
			w.Write(returnFinalFileResponse.Data)

		})

	}

}
