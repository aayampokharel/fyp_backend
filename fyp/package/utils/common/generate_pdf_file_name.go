package common

import (
	"strconv"
	"strings"
)

func GeneratePDFFileName(studentName, facultyName string, index int) string {
	studentNamewithUnderscore := strings.Split(studentName, " ")
	facultyNameWithUnderscore := strings.Split(facultyName, " ")
	facultyName = strings.Join(facultyNameWithUnderscore, "_")
	studentName = strings.Join(studentNamewithUnderscore, "_")
	name := studentName + "_" + facultyName + "_"
	if len(name) > 50 {
		name = name[:50]
	}
	name = name + strconv.Itoa(index) + ".pdf"
	return name

}
