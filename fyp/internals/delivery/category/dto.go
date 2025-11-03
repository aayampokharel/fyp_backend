package category

import (
	"project/internals/domain/entity"
	err "project/package/errors"
	"project/package/utils/common"
)

type CreatePDFCategoryDto struct {
	FacultyName           string `json:"faculty_name"`
	PreferredCategoryName string `json:"preferred_category_name"`
	InstitutionID         string `json:"institution_id"`
	InstitutionFacultyID  string `json:"institution_faculty_id"`
}
type CreatePDFCategoryResponseDto struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
}

func (m *CreatePDFCategoryDto) ToPdfFileCategoryEntity() (entity.PDFFileCategoryEntity, error) {

	if m.InstitutionFacultyID == "" || m.InstitutionID == "" || m.PreferredCategoryName == "" || m.FacultyName == "" {
		return entity.PDFFileCategoryEntity{}, err.ErrEmptyFields
	}
	return entity.PDFFileCategoryEntity{
		CategoryID:           common.GenerateUUID(16),
		CategoryName:         common.GeneratePDFCategoryName(m.FacultyName, m.PreferredCategoryName),
		InstitutionID:        m.InstitutionID,
		InstitutionFacultyID: m.InstitutionFacultyID,
	}, nil
}
