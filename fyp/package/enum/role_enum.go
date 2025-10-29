package enum

import (
	"strings"
)

type ROLE string

const (
	ADMIN     ROLE = "ADMIN"
	INSTITUTE ROLE = "INSTITUTE"
)

func (r ROLE) String() string {
	return string(r)
}

func StringToRole(roleString string) ROLE {
	roleStringUpper := strings.ToUpper(roleString)

	switch roleStringUpper {
	case "ADMIN":
		return ADMIN
	case "INSTITUTE":
		return INSTITUTE
	default:
		return INSTITUTE
	}

}
