package helper

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"project/constants"
	"project/internals/data/config"
	"project/package/utils/common"
)

// func VerifyDigitalSignature(env *config.Env, allPBFTNodes []string, leaderNode string, dataHash string, signatureString string) bool {
// 	fmt.Printf("=== DEBUG VerifyDigitalSignature ===\n")
// 	fmt.Printf("Leader node: %s\n", leaderNode)
// 	fmt.Printf("DataHash: %s (length: %d)\n", dataHash, len(dataHash))
// 	fmt.Printf("Signature: %s (length: %d)\n", signatureString, len(signatureString))

// 	// 1. Get public key
// 	var publicKeyValueFromEnv string
// 	switch leaderNode {
// 	case "9500":
// 		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode1)
// 	case "9501":
// 		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode2)
// 	case "9502":
// 		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode3)
// 	case "9503":
// 		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode4)
// 	}

// 	fmt.Printf("Public key from env length: %d\n", len(publicKeyValueFromEnv))

// 	publicKey, er := common.ConvertToRsaPublicKey(publicKeyValueFromEnv)
// 	if er != nil {
// 		fmt.Printf("❌ Public key conversion failed: %v\n", er)
// 		return false
// 	}
// 	fmt.Printf("✅ Public key parsed successfully\n")

// 	// 2. Decode the signature from Base64
// 	signatureBytes, er := base64.StdEncoding.DecodeString(signatureString)
// 	if er != nil {
// 		fmt.Printf("❌ Signature base64 decode failed: %v\n", er)
// 		return false
// 	}
// 	fmt.Printf("✅ Signature decoded, length: %d bytes\n", len(signatureBytes))

// 	// 3. Handle dataHash - check if it's hex encoded or raw
// 	var dataHashBytes []byte

// 	// Try to decode as hex first (common format)
// 	if hexBytes, err := hex.DecodeString(dataHash); err == nil {
// 		dataHashBytes = hexBytes
// 		fmt.Printf("✅ DataHash is hex encoded, decoded to %d bytes\n", len(dataHashBytes))
// 	} else {
// 		// If not hex, use as raw bytes
// 		dataHashBytes = []byte(dataHash)
// 		fmt.Printf("ℹ️ DataHash used as raw bytes, length: %d\n", len(dataHashBytes))
// 	}

// 	// 4. Hash the data for verification
// 	// hash := sha256.Sum256(dataHashBytes)
// 	// fmt.Printf("✅ SHA256 hash computed: %x\n", hash[:])

// 	// // 5. Verify the signature
// 	// er = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
// 	er = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, dataHashBytes, signatureBytes)

// 	if er != nil {
// 		fmt.Printf("❌ RSA signature verification failed: %v\n", er)
// 		return false
// 	}

//		fmt.Printf("✅ Signature verification SUCCESSFUL!\n")
//		return true
//	}
func VerifyDigitalSignature(env *config.Env, allPBFTNodes []string, leaderNode string, dataHash string, signatureString string) bool {
	fmt.Printf("=== DEBUG VerifyDigitalSignature ===\n")
	fmt.Printf("Leader node: %s\n", leaderNode)
	fmt.Printf("DataHash: %s (length: %d)\n", dataHash, len(dataHash))
	fmt.Printf("Signature: %s (length: %d)\n", signatureString, len(signatureString))

	// 1) Resolve public key for leader
	var publicKeyValueFromEnv string
	switch leaderNode {
	case "9500":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode1)
	case "9501":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode2)
	case "9502":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode3)
	case "9503":
		publicKeyValueFromEnv = env.GetValueForKey(constants.PublicKeyNode4)
	default:
		fmt.Printf("❌ Unknown leader node: %s\n", leaderNode)
		return false
	}

	fmt.Printf("Public key from env length: %d\n", len(publicKeyValueFromEnv))
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
	//    - If hex-decodes to 32 bytes => treat as already-SHA256 and use directly.
	//    - Else treat as raw bytes and compute SHA256(rawBytes).
	var hashed []byte
	if hexBytes, err := hex.DecodeString(dataHash); err == nil && len(hexBytes) == 32 {
		// Already a SHA256 digest encoded as hex (common case when leader sets Digest = hex.EncodeToString(hash))
		hashed = hexBytes
		fmt.Printf("✅ DataHash is hex-encoded SHA256, using it directly (%d bytes)\n", len(hashed))
	} else {
		// Not a hex-encoded 32-byte digest -> treat as raw data and compute sha256
		raw := []byte(dataHash)
		sum := sha256.Sum256(raw)
		hashed = sum[:] // 32 bytes
		fmt.Printf("ℹ️ DataHash used as raw bytes, computed SHA256: %x\n", hashed)
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
