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
)

func GenerateRandomCertificateData() entity.CertificateData {
	studentRandomNum := generateRandomNumber(len(studentNamesSlice))
	degreeRandomNum := generateRandomNumber(len(degreesSlice))
	universityRandomNum := generateRandomNumber(len(universitySlice))
	collegeRandomNum := generateRandomNumber(len(collegesSlice))
	divisionRandomNum := generateRandomNumber(len(divisionSlice))

	return entity.CertificateData{
		ID:                 common.GenerateUUID(6),
		StudentName:        studentNamesSlice[studentRandomNum],
		Degree:             degreesSlice[degreeRandomNum],
		UniversityName:     universitySlice[universityRandomNum],
		College:            collegesSlice[collegeRandomNum],
		Division:           divisionSlice[divisionRandomNum],
		CertificateDate:    time.Now(),
		PrincipalSignature: common.GenerateUUID(6),
		TuApproval:         "true",
	}
}
