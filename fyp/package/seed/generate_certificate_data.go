package seed

import (
	"project/internals/domain/entity"
	"project/package/utils/common"
	"time"
)

var (
	studentNamesSlice = []string{
		"aayam pokharel",
		"ram sharma",
		"shyam thapa",
		"hari prasad",
		"gita rai",
		"sita thapa",
		"laxmi shrestha",
		"krishna adhikari",
		"anita karki",
		"binod dhakal",
		"suman gurung",
		"Kaido muscularman",
		"Hancock luffy",
		"Whitebeard Strong",
		"Roger san",
		"Nami Bolman",
		"Usopp nakebhai",
		"Roronoa Zoro",
		"Sanji san",
		"Chopper san",
		"Robin san",
		"Monkey D. Luffy",
		"Monkey D. Garp",
		"Monkey D. Dragon",
		"Portgas D. Ace",
		"Gol D. Roger",
		"Edward Newgate",
		"Trafalgar D. Water Law",
		"Donquixote Doflamingo",
		"Charlotte Linlin",
		"Silvers Rayleigh",
		"Red Haired Shanks",
		"Buggy the Clown",
		"Vinsmoke Sanji",
		"Jinbe Aquaman",
		"Sabo Revolutionary",
		"Aokiji Admiral",
		"Kuzan lightwala",
		"Akainu firewala",
		"Blackbeard Pirate",
		"Dracule Mihawk",
		"Bartholomew Kuma",
	}
	universitySlice = []string{
		"TU",
		"KU",
		"PU",
		"Los Aneles University",
		"Harvard University",
		"Stanford University",
		"MIT",
		"Cambridge University",
		"Oxford University",
		"Princeton University",
		"Yale University",
		"Columbia University",
		"University of Chicago",
		"University of California, Berkeley",
		"California Institute of Technology (Caltech)",
		"University of Pennsylvania",
		"University of Michigan",
		"University of Toronto",
		"University of Washington",
		"University of Edinburgh",
		"University of Tokyo",
		"National University of Singapore (NUS)",
		"University of Melbourne",
		"University of Sydney",
		"University of British Columbia (UBC)",
		"University of Alberta",
	}

	degreesSlice = []string{
		"Bachelor's",
		"Master's",
		"PhD",
		"Doctorate",
		"Associate's",
		"Certificate",
		"Diploma",
		"Professional Degree",
		"Vocational Training",
		"Continuing Education",
	}

	collegesSlice = []string{
		"St. Xavier's College",
		"Trinity College",
		"St. Mary's College",
		"St. John's College",
		"King's College",
		"Ambition College",
		"Liberty College",
		"Patan Colllege",
		"Pulchowk College",
		"Thapathali College",
		"Kalyanpur College",
		"Dhulikhel College",
		"BP College",
		"Koteshwor College",
		"Khwopa College",
		"Samriddhi College",
		"Chandragiri College",
		"Kantipur College",
		"Kathmandu College",
		"Chitwan College",
	}

	divisionSlice = []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"distinction",
	}

	majorsSlice = []string{
		"Computer Science",
		"Electrical Engineering",
		"Mechanical Engineering",
		"Civil Engineering",
		"Business Administration",
		"Medicine",
		"Law",
		"Architecture",
		"Biotechnology",
		"Data Science",
		"Artificial Intelligence",
		"Cybersecurity",
		"Finance",
		"Marketing",
		"Psychology",
		"Physics",
		"Chemistry",
		"Mathematics",
		"Economics",
		"Political Science",
	}

	certificateTypesSlice = []string{
		"DEGREE",
		"DIPLOMA",
		"TRANSCRIPT",
		"CERTIFICATE",
		"PROFESSIONAL_LICENSE",
	}

	cgpaSlice = []string{
		"3.8",
		"3.9",
		"3.5",
		"3.7",
		"3.6",
		"3.4",
		"3.2",
		"3.0",
		"3.1",
		"3.3",
		"2.8",
		"2.9",
		"2.7",
		"2.5",
		"2.6",
	}
)

func GenerateRandomCertificateData() entity.CertificateData {
	studentRandomNum := generateRandomNumber(len(studentNamesSlice))
	degreeRandomNum := generateRandomNumber(len(degreesSlice))
	universityRandomNum := generateRandomNumber(len(universitySlice))
	collegeRandomNum := generateRandomNumber(len(collegesSlice))
	divisionRandomNum := generateRandomNumber(len(divisionSlice))
	majorRandomNum := generateRandomNumber(len(majorsSlice))
	certTypeRandomNum := generateRandomNumber(len(certificateTypesSlice))
	cgpaRandomNum := generateRandomNumber(len(cgpaSlice))

	// Generate dates with realistic academic timelines
	enrollmentDate := generateRandomPastDate(4, 6)    // 4-6 years ago
	completionDate := enrollmentDate.AddDate(4, 0, 0) // Typically 4 years after enrollment
	issueDate := completionDate.AddDate(0, 3, 0)      // Issued 3 months after completion

	// Generate certificate data first to create hash
	certData := entity.CertificateData{
		CertificateID:   common.GenerateUUID(8),
		StudentID:       "STU" + common.GenerateUUID(6),
		StudentName:     studentNamesSlice[studentRandomNum],
		UniversityName:  universitySlice[universityRandomNum],
		Degree:          degreesSlice[degreeRandomNum],
		College:         collegesSlice[collegeRandomNum],
		Major:           majorsSlice[majorRandomNum],
		GPA:             cgpaSlice[cgpaRandomNum],
		Percentage:      0.0,
		Division:        divisionSlice[divisionRandomNum],
		EnrollmentDate:  enrollmentDate,
		CompletionDate:  completionDate,
		IssueDate:       issueDate,
		CertificateType: certificateTypesSlice[certTypeRandomNum],
		CreatedAt:       time.Now(),
		IssuerPublicKey: "PUBKEY_" + common.GenerateUUID(8),
	}
	if certData.GPA == "" {
		certData.Percentage = 53.33
	}

	// Generate data hash based on the certificate content
	certData.DataHash = generateCertificateHash(certData)

	return certData
}

func generateRandomPastDate(minYears, maxYears int) time.Time {
	yearsAgo := generateRandomNumber(maxYears-minYears+1) + minYears
	daysAgo := generateRandomNumber(365)
	hoursAgo := generateRandomNumber(24)
	minutesAgo := generateRandomNumber(60)

	pastDate := time.Now().AddDate(-yearsAgo, 0, -daysAgo)
	pastDate = pastDate.Add(-time.Duration(hoursAgo) * time.Hour)
	pastDate = pastDate.Add(-time.Duration(minutesAgo) * time.Minute)

	return pastDate
}

func generateCertificateHash(cert entity.CertificateData) string {
	// Create a hash from the certificate data for integrity verification
	data := cert.CertificateID + cert.StudentID + cert.StudentName +
		cert.UniversityName + cert.Degree + cert.College + cert.Major +
		cert.GPA + cert.Division + cert.EnrollmentDate.Format("2006-01-02") +
		cert.CompletionDate.Format("2006-01-02") + cert.IssueDate.Format("2006-01-02") +
		cert.CertificateType + cert.IssuerPublicKey

	hash, er := common.HashData(data)
	if er != nil {
		return ""
	}
	return hash
}
