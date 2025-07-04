package digitalsignature

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	//read pem file ,
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	//parse it to get the key
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	//again parse the key itself
	rawKeyData, _ := pem.Decode(pemBlock.Bytes)
	if err != nil {

	}
	//again decode for removing asn.1 DER
	key, err := x509.ParsePKCS8PrivateKey(rawKeyData.Bytes)
	if err != nil {
		return nil, err
	}
	//(for safety) : type assert this
	rsaPrivKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return rsaPrivKey, nil
}
