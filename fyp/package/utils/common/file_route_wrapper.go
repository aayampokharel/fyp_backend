package common

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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

	Method             enum.HTTPMETHOD
	RequestDataType    interface{}
	GetORDeleteHandler func(map[string]string) entity.FileResponse
	URLQueries         []string
	PostHandler        func(interface{}) entity.FileResponse
}

func setFileHeaders(w http.ResponseWriter, filename string, fileType enum.RESPONSETYPE, contentLength int) {
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
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
}

func NewFileRouteWrapper(routeInfos ...FileRouteWrapper) {
	for index := range routeInfos {
		routeInfo := routeInfos[index]
		routeInfo.Mux.HandleFunc(routeInfo.Prefix+routeInfo.Route, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != routeInfo.Method.ToString() {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var returnFinalFileResponse entity.FileResponse

			// GET / DELETE
			if routeInfo.RequestDataType == nil && (routeInfo.Method == enum.METHODGET || routeInfo.Method == enum.METHODDELETE) {
				queryParams := make(map[string]string)
				for _, val := range routeInfo.URLQueries {
					queryParams[val] = r.URL.Query().Get(val)
				}
				returnFinalFileResponse = routeInfo.GetORDeleteHandler(queryParams)

			} else { // POST / PUT with body
				reqType := reflect.TypeOf(routeInfo.RequestDataType)
				reqValue := reflect.New(reqType).Interface()
				if er := json.NewDecoder(r.Body).Decode(reqValue); er != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(HandleErrorResponse(400, err.ErrDecodingJSONString, er))
					return
				}
				returnFinalFileResponse = routeInfo.PostHandler(reflect.ValueOf(reqValue).Elem().Interface())
			}

			// If there is no data, send JSON response
			if returnFinalFileResponse.Data == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(returnFinalFileResponse.Code)
				json.NewEncoder(w).Encode(entity.Response{
					Code:    returnFinalFileResponse.Code,
					Message: returnFinalFileResponse.Message,
					Data:    nil,
				})
				return
			}

			setFileHeaders(w, returnFinalFileResponse.FileName, returnFinalFileResponse.FileType, len(returnFinalFileResponse.Data))
			w.WriteHeader(returnFinalFileResponse.Code)

			reader := bytes.NewReader(returnFinalFileResponse.Data)
			if _, err := io.Copy(w, reader); err != nil {
				log.Println("Error streaming file:", err)
			}
		})
	}
}
