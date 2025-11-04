package enum

import "strings"

type INSTITUTION_ACTIVE string

const (
	INSTITUTIONACCEPETED INSTITUTION_ACTIVE = "TRUE"
	INSTITUTIONREJECTED  INSTITUTION_ACTIVE = "FALSE"
	INSTITUTIONPENDING   INSTITUTION_ACTIVE = "PENDING"
)

func (i INSTITUTION_ACTIVE) String() string {
	return string(i)
}

func StringToInstitutionActive(institutionActiveString string) INSTITUTION_ACTIVE {
	institutionActiveStringUpper := strings.ToUpper(institutionActiveString)

	switch institutionActiveStringUpper {
	case "TRUE":
		return INSTITUTIONACCEPETED
	case "FALSE":
		return INSTITUTIONREJECTED
	default:
		return INSTITUTIONPENDING
	}
}
