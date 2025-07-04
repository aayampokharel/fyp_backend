package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func calculateHash(anyData interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(anyData)
	if err != nil {
		return nil, fmt.Errorf("json conversion for hash")
	}
	hashByte := sha256.Sum256(jsonData)

	return hashByte[:], nil

}

func CalculateHashHex(data interface{}) (string, error) {
	hash, err := calculateHash(data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}
