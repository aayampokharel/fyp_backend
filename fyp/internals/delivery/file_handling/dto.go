package filehandling

type GetRequestQueryType []string

const CategoryId = "category_id"
const CategoryName = "category_name"
const FileID = "file_id"
const IsDownloadAll = "is_download_all"
const CertificateID = "certificate_id"
const CertificateHash = "certificate_hash"

var GetHTMLRequestQuery = GetRequestQueryType{CertificateHash, CertificateID}
var GetPDFFileInListQuery = GetRequestQueryType{CategoryId, CategoryName, FileID, IsDownloadAll}

type GetImageFileRequestDto struct {
	ImageBase64 string `json:"image_base64"`
	ImageName   string `json:"image_name"`
}
