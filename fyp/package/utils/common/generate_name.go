package common

import (
	"fmt"
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

func GenerateFileNameWithExtension(baseName string, noOfDigitsForUUID int, extension string) string {
	name := fmt.Sprintf("%s_%s.%s", baseName, GenerateUUID(noOfDigitsForUUID), extension)
	if len(name) > 50 {
		name = strings.ReplaceAll(name[:50], ".", "_")
		ext := fmt.Sprintf(".%s", extension)
		name = name + ext
	}
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
