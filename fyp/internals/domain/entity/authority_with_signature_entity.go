package entity

import (
	"encoding/json"
	"log"
)

type AuthorityWithSignatureEntity struct {
	AuthorityName        string `json:"authority_name"`
	SignatureImageBase64 string `json:"signature_image_base64"`
}

func (a AuthorityWithSignatureEntity) FromString(jsonString string) ([]AuthorityWithSignatureEntity, error) {
	var authorities []AuthorityWithSignatureEntity
	er := json.Unmarshal([]byte(jsonString), &authorities)
	if er != nil {

		log.Fatal("Failed to parse authority JSON:", er)
		return nil, er
	}
	return authorities, nil
}
