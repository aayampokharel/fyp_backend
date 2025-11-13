package common

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	err "project/package/errors"
)

func HashData(data interface{}) (string, []byte, error) {
	bytes, error := json.Marshal(data)
	if error != nil {
		return "", nil, err.ErrMarshaling
	}

	hash := sha256.Sum256(bytes)
	return fmt.Sprintf("%x", hash), hash[:], nil
}
