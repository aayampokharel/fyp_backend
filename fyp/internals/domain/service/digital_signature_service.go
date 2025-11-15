package service

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"project/constants"
	"project/internals/data/config"
	err "project/package/errors"
	"project/package/utils/common"
	logger "project/package/utils/pkg"
	"strings"

	"go.uber.org/zap"
)

type DigitalSignatureService struct {
	logger *zap.SugaredLogger
	env    *config.Env
}

func NewDigitalSignature() *DigitalSignatureService {
	envConfig, er := config.NewEnv()
	if er != nil {
		return nil
	}
	return &DigitalSignatureService{logger: logger.Logger, env: envConfig}
}

func (d *DigitalSignatureService) SignMessageWithHash(hashedData []byte) (string, error) {
	privateKeyFlag := common.GetPrivatekey()
	privateKey, _, er := privateOrPublicKeyFromBase64DER(*privateKeyFlag, true)
	if er != nil {
		return "", er
	}
	if privateKey == nil {
		return "", err.ErrEmptyPrivateKey
	}

	// Sign the pre-hashed data directly
	signatureByte, er := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashedData)
	if er != nil {
		return "", er
	}
	return base64.StdEncoding.EncodeToString(signatureByte), nil
}

func privateOrPublicKeyFromBase64DER(keyString string, isTypePrivateKey bool) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	fmt.Printf("=== DEBUG KEY PARSING ===\n")
	fmt.Printf("Key length: %d\n", len(keyString))
	//fmt.Printf("Key: %s\n", keyString)
	fmt.Printf("First 50 chars: %s\n", safeSubstring(keyString, 0, 50))

	cleanKey := strings.TrimSpace(keyString)

	if isTypePrivateKey {
		privKey, pubKey, err := parsePrivateKey(cleanKey)
		if err != nil || privKey == nil {
			return nil, nil, fmt.Errorf("invalid private key: %v", err)
		}
		return privKey, pubKey, nil
	} else {
		_, pubKey, err := parsePublicKey(cleanKey)
		if err != nil || pubKey == nil {
			return nil, nil, fmt.Errorf("invalid public key: %v", err)
		}
		return nil, pubKey, nil
	}
}

// parsePrivateKey handles private key parsing from different formats
func parsePrivateKey(cleanKey string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if strings.Contains(cleanKey, "BEGIN RSA PRIVATE KEY") {
		block, _ := pem.Decode([]byte(cleanKey))
		if block == nil {
			return nil, nil, fmt.Errorf("failed to parse PEM block")
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse private key from PEM: %v", err)
		}
		fmt.Printf("✅ Parsed private key from PEM format\n")
		return privateKey, &privateKey.PublicKey, nil
	} else {
		keyBytes, err := base64.StdEncoding.DecodeString(cleanKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode base64 private key: %v", err)
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse private key from base64: %v", err)
		}
		fmt.Printf("✅ Parsed private key from raw base64 format\n")
		return privateKey, &privateKey.PublicKey, nil
	}
}

// parsePublicKey handles public key parsing from different formats
func parsePublicKey(cleanKey string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var derBytes []byte
	var err error

	// Check if it's PEM format
	if strings.Contains(cleanKey, "BEGIN") {
		block, _ := pem.Decode([]byte(cleanKey))
		if block == nil {
			return nil, nil, fmt.Errorf("failed to parse PEM block for public key")
		}
		derBytes = block.Bytes
		fmt.Printf("✅ Extracted public key from PEM format\n")
	} else {
		// Parse as raw base64
		derBytes, err = base64.StdEncoding.DecodeString(cleanKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode base64 public key: %v", err)
		}
		fmt.Printf("✅ Decoded public key from raw base64 format\n")
	}

	// Try PKCS#1 format first
	publicKey, err := x509.ParsePKCS1PublicKey(derBytes)
	if err == nil {
		fmt.Printf("✅ Parsed public key as PKCS#1 format\n")
		return nil, publicKey, nil
	}

	// Fall back to PKIX format
	pubKey, err := x509.ParsePKIXPublicKey(derBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key (tried PKCS#1 and PKIX): %v", err)
	}

	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("not an RSA public key, got: %T", pubKey)
	}

	fmt.Printf("✅ Parsed public key as PKIX format\n")
	return nil, rsaPublicKey, nil
}

// Helper function for safe substring
func safeSubstring(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		return ""
	}
	return s[start:end]
}

func (d *DigitalSignatureService) VerifySignature(message interface{}, signature string, nodeNumber int) error {
	_, hashedByte, er := common.HashData(message)
	if er != nil {
		d.logger.Errorw("[service] Error: VerifySignature::", zap.Error(er))
		return er
	}
	var publicKeyString string
	switch nodeNumber {
	case 9500:
		publicKeyString = d.env.GetValueForKey(constants.PublicKeyNode1)
	case 9501:
		publicKeyString = d.env.GetValueForKey(constants.PublicKeyNode2)
	case 9502:
		publicKeyString = d.env.GetValueForKey(constants.PublicKeyNode3)
	}
	_, publicKey, er := privateOrPublicKeyFromBase64DER(publicKeyString, false)
	if er != nil {
		return er
	}
	if publicKey == nil {
		return err.ErrEmptyPublicKey
	}
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %v", err)
	}
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashedByte, sigBytes, nil)
}
