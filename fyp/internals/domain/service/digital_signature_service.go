package service

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"project/constants"
	"project/internals/data/config"
	err "project/package/errors"
	"project/package/utils/common"
	logger "project/package/utils/pkg"

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

func (d *DigitalSignatureService) SignMessage(message interface{}) (string, error) {
	privateKeyFlag := common.GetPrivatekey()
	_, hashedByte, er := common.HashData(message)
	if er != nil {
		d.logger.Errorw("[service] Error: SignMessage::", zap.Error(er))
		return "", er
	}

	privateKey, _, er := privateOrPublicKeyFromBase64DER(*privateKeyFlag, true)
	if er != nil {
		return "", er
	}
	if privateKey == nil {
		return "", err.ErrEmptyPrivateKey
	}
	signatureByte, er := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashedByte, nil)
	if er != nil {
		return "", er
	}
	return base64.StdEncoding.EncodeToString(signatureByte), nil
}

func privateOrPublicKeyFromBase64DER(base64String string, isTypePrivateKey bool) (*rsa.PrivateKey, *rsa.PublicKey, error) {

	derBytes, er := base64.StdEncoding.DecodeString(base64String)
	if er != nil {
		return nil, nil, fmt.Errorf("failed to decode base64: %v", er)
	}

	if isTypePrivateKey {
		privateKey, er := x509.ParsePKCS1PrivateKey(derBytes)
		if er == nil {
			return privateKey, nil, nil
		}

		priv, er := x509.ParsePKCS8PrivateKey(derBytes)
		if er != nil {
			return nil, nil, fmt.Errorf("failed to parse private key: %v", er)
		}

		rsaPrivateKey, ok := priv.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, fmt.Errorf("not an RSA private key")
		}

		return rsaPrivateKey, nil, nil
	} else {

		publicKey, er := x509.ParsePKCS1PublicKey(derBytes)
		if er == nil {
			return nil, publicKey, nil
		}

		priv, er := x509.ParsePKIXPublicKey(derBytes)
		if er != nil {
			return nil, nil, fmt.Errorf("failed to parse public key: %v", er)
		}

		rsaPublicKey, ok := priv.(*rsa.PublicKey)
		if !ok {
			return nil, nil, fmt.Errorf("not an RSA public key")
		}
		return nil, rsaPublicKey, nil
	}
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
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashedByte, []byte(signature), nil)
}
