package utils

import (
	"fmt"
	"log"
)

const (
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

// ReturnError creates a new error that combines a custom message with the original error
func ReturnError(customErrString string, err error) error {
	return fmt.Errorf("ERROR %s: %w", customErrString, err)
}

// LogError prints the error message in red color to the console
func LogError(err error) {
	log.Printf("%sERROR: %v%s", colorRed, err, colorReset)
}

// LogErrorWithContext prints a custom error message along with the error in red color
func LogErrorWithContext(context string, err error) {
	log.Printf("%sERROR [%s]: %v%s", colorRed, context, err, colorReset)

}
