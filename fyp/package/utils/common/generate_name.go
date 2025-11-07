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

func GeneratePDFCategoryName(faculty, preferredName string) string {
	name := preferredName + "_" + faculty
	if len(name) > 50 {
		name = name[:50]
	}
	name = name + GenerateUUID(8)
	return name
}
