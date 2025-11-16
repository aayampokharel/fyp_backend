package entity

type AdminDashboardCountsEntity struct {
	// Institutions
	ActiveInstitutions   int
	DeletedInstitutions  int
	SignedUpInstitutions int

	// Faculties
	TotalFaculties int

	// Users
	ActiveUsers      int
	DeletedUsers     int
	ActiveAdmins     int
	ActiveInstitutes int

	// PDF
	TotalPDFCategories int
	TotalPDFFiles      int

	// Certificates
	TotalCertificates int

	// Blockchain
	TotalBlocks int
}
