package digitalsignature

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/aayampokharel/fyp/utils"
)

func GeneratePrivateKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048) // 2048-bit key random
	if err != nil {
		utils.LogErrorWithContext("GenerateRSAKeyPair", err)
	}
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",                       // PEM header
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey), // Convert key to ASN.1 DER format
	}

	privateKeyFile, err := os.Create("../../certs/private_key.pem")
	if err != nil {
		utils.LogErrorWithContext("StoreKey", err)
		return nil, nil, err
	}
	defer privateKeyFile.Close()

	// Write PEM block to file
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		utils.LogErrorWithContext("EncodeKey", err)
		return nil, nil, err
	}
	publicKey := privateKey.PublicKey
	return privateKey, &publicKey, nil
}
