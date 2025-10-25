package common

import (
	"net/http"
	"project/internals/domain/entity"
)

func HandleErrorResponse(statusCode int, message string, er error) entity.Response {
	var errorString string

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
	return errorResponse
}

func HandleSuccessResponse(data interface{}) entity.Response {

	successResponse := entity.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
	// json.NewEncoder(w).Encode(successResponse)
	return successResponse
}
