package filehandling

import (
	"project/internals/domain/entity"
	"project/package/enum"
)

type GetRequestQueryType map[string]string

var GetHTMLRequestQuery = GetRequestQueryType{"id": ""}
var GetPDFFileInListQuery = GetRequestQueryType{"category_id": "", "file_id": "", "is_download_all": ""}

type ResponseWithFileTypeAndCount struct {
	FileType enum.RESPONSETYPE
	Count    int
	Data     []entity.PDFFileEntity
}
