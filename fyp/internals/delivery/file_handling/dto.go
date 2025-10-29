package filehandling

type GetRequestQueryType []string

const CategoryId = "category_id"
const CategoryName = "category_name"
const FileID = "file_id"
const IsDownloadAll = "is_download_all"

var GetHTMLRequestQuery = GetRequestQueryType{"id"}
var GetPDFFileInListQuery = GetRequestQueryType{CategoryId, CategoryName, FileID, IsDownloadAll}
