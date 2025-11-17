package entity

type AdminDashboardCountsEntity struct {
	// Institutions
	ActiveInstitutions   int `json:"active_institutions"`
	DeletedInstitutions  int `json:"deleted_institutions"`
	SignedUpInstitutions int `json:"signed_up_institutions"`

	// Faculties
	TotalFaculties int `json:"total_faculties" `

	// Users
	ActiveUsers      int `json:"active_users"`
	DeletedUsers     int `json:"deleted_users"`
	ActiveAdmins     int `json:"active_admins"`
	ActiveInstitutes int `json:"active_institutes"`

	// PDF
	TotalPDFCategories int `json:"total_pdf_categories"`
	TotalPDFFiles      int `json:"total_pdf_files"`

	// Certificates
	TotalCertificates int `json:"total_certificates"`

	// Blockchain
	TotalBlocks int `json:"total_blocks"`
}
