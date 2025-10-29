package entity

import "project/package/enum"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FileResponse struct {
	Code     int               `json:"code"`
	FileName string            `json:"file_name"`
	Message  string            `json:"message"`
	Data     []byte            `json:"data"`
	FileType enum.RESPONSETYPE `json:"file_type"`
}
