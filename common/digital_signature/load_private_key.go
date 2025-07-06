package digitalsignature

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/aayampokharel/fyp/utils"
)

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	//read pem file ,
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		utils.LogError(fmt.Errorf("failed to read PEM file"))
		return nil, err
	}
	//parse it to get the key
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		utils.LogError(fmt.Errorf("failed to decode PEM file"))
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	if pemBlock.Type != "PRIVATE KEY" && pemBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM type: %s", pemBlock.Type)
	}

	// //again parse the key itself
	// rawKeyData, _ := pem.Decode(pemBlock.Bytes)
	// if rawKeyData == nil {
	// 	utils.LogErrorWithContext("DecodePrivateKey", err)
	// 	return nil, err
	// }
	//again decode for removing asn.1 DER
	rsaPrivKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		utils.LogErrorWithContext("ParsePrivateKey", err)
		return nil, err
	}
	//(for safety) : type assert this

	return rsaPrivKey, nil
}
