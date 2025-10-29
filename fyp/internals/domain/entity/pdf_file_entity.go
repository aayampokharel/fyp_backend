package entity

import "time"

type PDFFileEntity struct {
	FileID     string     `json:"file_id"`
	CategoryID string     `json:"category_id"`
	PDFData    []byte     `json:"pdf_data"`
	FileName   string     `json:"file_name"`
	UploadedAt *time.Time `json:"uploaded_at"`
}
