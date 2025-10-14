package entity

type InstitutionFaculty struct {
	InstitutionFacultyID      string `json:"institution_faculty_id"`
	InstitutionID             string `json:"institution_id"`
	Faculty                   string `json:"faculty"`
	FacultyHODName            string `json:"faculty_hod_name"`
	FacultyHODSignatureBase64 string `json:"faculty_hod_signature_base64"`
}
