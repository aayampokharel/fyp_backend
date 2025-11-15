package helper

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"project/constants"
	"project/internals/data/config"
	"project/package/utils/common"
)

func VerifyDigitalSignature(env *config.Env, currentNode string, dataHash []byte, signatureString string) bool {
	fmt.Printf("=== DEBUG VerifyDigitalSignature ===\n")
	fmt.Printf("Current node: %s\n", currentNode)
	fmt.Printf("DataHash length: %d\n", len(dataHash))
	fmt.Printf("Signature length: %d\n", len(signatureString))

	// 1) Resolve public key for leader
	var publicKeyValueFromEnv string
	switch currentNode {
	case "9500":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode1)
	case "9501":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode2)
	case "9502":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode3)
	case "9503":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode4)
	default:
		fmt.Printf("❌ Unknown leader node: %s\n", currentNode)
		return false
	}

	publicKey, er := common.ConvertToRsaPublicKey(publicKeyValueFromEnv)
	if er != nil {
		fmt.Printf("❌ Public key conversion failed: %v\n", er)
		return false
	}
	fmt.Printf("✅ Public key parsed successfully\n")

	// 2) Decode signature from base64
	signatureBytes, er := base64.StdEncoding.DecodeString(signatureString)
	if er != nil {
		fmt.Printf("❌ Signature base64 decode failed: %v\n", er)
		return false
	}
	fmt.Printf("✅ Signature decoded, length: %d bytes\n", len(signatureBytes))

	// 3) Interpret dataHash:
	// If dataHash is already 32 bytes, assume it's SHA256 digest
	// Otherwise, compute SHA256 from the raw bytes
	var hashed []byte
	if len(dataHash) == 32 {
		hashed = dataHash
		fmt.Printf("✅ DataHash is already 32-byte SHA256 digest\n")
	} else {
		sum := sha256.Sum256(dataHash)
		hashed = sum[:]
		fmt.Printf("ℹ️ DataHash computed SHA256: %x\n", hashed)
	}

	// 4) Verify signature: pass the *digest* bytes directly (do NOT hash again)
	er = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, signatureBytes)
	if er != nil {
		fmt.Printf("❌ RSA signature verification failed: %v\n", er)
		return false
	}

	fmt.Printf("✅ Signature verification SUCCESSFUL!\n")
	return true
}
