package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	err "project/package/errors"
	"strconv"
	"strings"
)

func CheckMapKeysReturnValues(providedMap map[string]string, keys []string) (map[string]string, error) {
	for _, key := range keys {
		providedMapValue, exists := providedMap[key]
		if !exists {
			return nil, err.ErrWithMoreInfo(nil, "key doesnot exist")
		}
		if providedMapValue == "" {
			return nil, err.ErrWithMoreInfo(nil, key+"is required")
		}

	}
	return providedMap, nil
}

func ConvertToBool(boolStr string) (bool, error) {
	if boolStr == "" {
		return false, err.ErrEmptyString
	}
	val, er := strconv.ParseBool(boolStr)
	if er != nil {
		return false, err.ErrCannotConvertToBool
	}
	return val, nil
}
func ConvertDigestHexToBytes(hexDigest string) ([]byte, error) {
	cleanDigest := strings.TrimPrefix(hexDigest, "0x")

	digestBytes, err := hex.DecodeString(cleanDigest)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex digest: %v", err)
	}

	return digestBytes, nil
}

func ConvertToInt(intStr string) (int, error) {
	if intStr == "" {
		return 0, err.ErrEmptyString
	}
	val, er := strconv.Atoi(intStr)
	if er != nil {
		return 0, err.ErrCannotConvertToInt
	}
	return int(val), nil
}

func ConvertToFloat(floatStr string) (float64, error) {
	if floatStr == "" {
		return 0.0, err.ErrEmptyString
	}
	val, er := strconv.ParseFloat(floatStr, 64)
	if er != nil {
		return 0.0, err.ErrCannotConvertToFloat
	}
	return val, nil
}

func ConvertToRsaPublicKey(publicKey string) (*rsa.PublicKey, error) {
	fmt.Printf("=== DEBUG ConvertToRsaPublicKey ===\n")
	//fmt.Printf("Original publicKey: '%s'\n", publicKey)

	cleanKey := strings.Split(publicKey, " ")[0]
	cleanKey = strings.TrimSpace(cleanKey)
	fmt.Printf("After cleaning: '%s'\n", cleanKey)

	// ✅ Use RSA PUBLIC KEY headers for PKCS#1 format
	pemKey := "-----BEGIN RSA PUBLIC KEY-----\n" +
		cleanKey + "\n" +
		"-----END RSA PUBLIC KEY-----"

	//fmt.Printf("PEM key:\n%s\n", pemKey)

	block, rest := pem.Decode([]byte(pemKey))
	fmt.Printf("PEM decode - block: %v, rest: %v\n", block != nil, string(rest))

	if block == nil {
		fmt.Printf("❌ PEM decode failed with RSA PUBLIC KEY headers\n")
		return nil, err.ErrConvertToRsa
	}

	fmt.Printf("PEM Block Type: %s\n", block.Type)
	fmt.Printf("PEM Block Bytes length: %d\n", len(block.Bytes))

	// ✅ Use ParsePKCS1PublicKey for PKCS#1 format
	rsaPub, er := x509.ParsePKCS1PublicKey(block.Bytes)
	if er != nil {
		fmt.Printf("❌ ParsePKCS1PublicKey failed: %v\n", er)
		return nil, err.ErrConvertToRsa
	}

	fmt.Printf("✅ Success with ParsePKCS1PublicKey! Key size: %d bits\n", rsaPub.Size()*8)
	return rsaPub, nil
}
