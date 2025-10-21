package common

import (
	"encoding/json"
	"net/http"
	"project/internals/domain/entity"
)

func HandleErrorResponse(statusCode int, message string, er error, w http.ResponseWriter) {
	var errorString string
	(w).Header().Set("Content-Type", "application/json")
	(w).WriteHeader(statusCode)

	if er != nil {
		errorString = er.Error()
	} else {
		errorString = "error:"
	}
	errorResponse := entity.Response{
		Code:    statusCode,
		Message: errorString + message,
		Data:    nil,
	}

	json.NewEncoder(w).Encode(errorResponse)
}

func HandleSuccessResponse(data interface{}, w http.ResponseWriter) {
	(w).Header().Set("Content-Type", "application/json")
	(w).WriteHeader(http.StatusOK)

	successResponse := entity.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
	json.NewEncoder(w).Encode(successResponse)
}
