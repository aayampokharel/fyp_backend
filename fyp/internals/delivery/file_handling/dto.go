package filehandling

type GetRequestQueryType []string

const CategoryId = "category_id"
const CategoryName = "category_name"
const FileID = "file_id"
const IsDownloadAll = "is_download_all"

var GetHTMLRequestQuery = GetRequestQueryType{"id"}
var GetPDFFileInListQuery = GetRequestQueryType{CategoryId, CategoryName, FileID, IsDownloadAll}

type GetImageFileRequestDto struct {
	ImageBase64 string `json:"image_base64"`
	ImageName   string `json:"image_name"`
}
