package common

import (
	"encoding/json"
	"fmt"
)

func PrintPrettyJSON(v any) {
	jsonBytes, err := json.MarshalIndent(v, "", "  ") // two-space indentation
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
