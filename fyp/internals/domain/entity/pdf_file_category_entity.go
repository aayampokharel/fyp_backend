package entity

import "time"

type PDFFileCategoryEntity struct {
	CategoryID           string
	InstitutionID        string
	InstitutionFacultyID string
	CategoryName         string
	CreatedAt            time.Time
}
