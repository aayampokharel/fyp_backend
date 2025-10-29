package entity

type InstitutionFaculty struct {
	InstitutionFacultyID string `json:"institution_faculty_id"`
	InstitutionID        string `json:"institution_id"`
	FacultyName          string `json:"faculty_name"`
	// faculty_authority_with_signature JSONB NOT NULL DEFAULT '[]'
	FacultyAuthorityWithSignatures []map[string]string `json:"faculty_authority_with_signatures"`
	UniversityAffiliation          string              `json:"university_affiliation"`
	UniversityCollegeCode          string              `json:"university_college_code"`
}
