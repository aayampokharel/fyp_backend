package common

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	err "project/package/errors"
)

func HashData(data interface{}) (string, error) {
	bytes, error := json.Marshal(data)
	if error != nil {
		return "", err.ErrMarshaling
	}

	hash := sha256.Sum256(bytes)
	return fmt.Sprintf("%x", hash), nil
}
