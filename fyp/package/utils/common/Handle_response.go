package common

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
)

func HandleErrorResponse(statusCode int, message string, er error) entity.Response {
	var errorString string

	if er != nil {
		errorString = er.Error()
	} else {
		errorString = "ERROR:"
	}
	errorResponse := entity.Response{
		Code:    statusCode,
		Message: errorString + message,
		Data:    nil,
	}
	return errorResponse
}

func HandleFileErrorResponse(statusCode int, message string, er error) entity.FileResponse {
	var errorString string

	if er != nil {
		errorString = er.Error()
	}
	return entity.FileResponse{
		Code:    statusCode,
		Message: "ERROR:" + message + " :: " + errorString,
		Data:    nil,
	}
}

func HandleFileSuccessResponse(fileType enum.RESPONSETYPE, fileName string, data []byte) entity.FileResponse {
	return entity.FileResponse{
		FileName: fileName,
		FileType: fileType,
		Code:     http.StatusOK,
		Message:  "SUCCESS",
		Data:     data,
	}
}

func HandleSuccessResponse(data interface{}) entity.Response {

	successResponse := entity.Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    data,
	}
	// json.NewEncoder(w).Encode(successResponse)
	return successResponse
}
