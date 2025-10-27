package service

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func (s *Service) GenerateQRCodeBase64(id, url string) (string, error) {
	//https://example.com/verify?id=123 use this route as preview is inside for institution for verification use above
	qrBytes, err := qrcode.Encode("https://example.com/verify?id="+id, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)

	return qrBase64, nil

}
