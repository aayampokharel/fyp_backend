package block

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"log"
	"time"

	digitalsignature "github.com/aayampokharel/fyp/common/digital_signature"
	"github.com/aayampokharel/fyp/models"
	"github.com/aayampokharel/fyp/utils"
)

// func CreateBlock(w http.ResponseWriter, r *http.Request) {

// }

func CreateBlock() {
	//create 1st instance .
	// merkel tree ma update .
	//pow + nonce calc
	// insert without error
	//! change this below as well . maile directly eeuta certificatedata ko lagi matra garne ho ahile ko lagi . not for all . take the POW () , merkelroot() inside INSERT function . do there ,
	fakeCertificateWithSignature := models.CertificateDataWithSignature{
		Signature: "",
		CertificateData: models.CertificateData{

			ID:                 "",
			StudentName:        "",
			UniversityName:     "",
			Degree:             "",
			College:            "",
			CertificateDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			Division:           "",
			PrincipalSignature: "",
			TuApproval:         "",
		},
	}

	//create a signature
	//-> hash

	hashedValue, err := HashCertificateData(fakeCertificateWithSignature.CertificateData)
	if err != nil {
		utils.LogErrorWithContext("HashBlock", err)
		return
	}

	_, _, err = digitalsignature.GeneratePrivateKey()
	if err != nil {
		utils.LogErrorWithContext("GeneratePrivateKey", err)
		return
	}

	rsaPrivateKey, err := digitalsignature.LoadPrivateKey("../../certs/private_key.pem")
	if err != nil {
		utils.LogErrorWithContext("LoadPrivateKey", err)
		return
	}
	rsaPublicKey := rsaPrivateKey.PublicKey

	//->sign using standard terms(dihital signature creatin bhayo . )
	digitalSignatureByte, er := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashedValue)
	if er != nil {
		utils.LogErrorWithContext("Digital Sign Byte Creation", er)
		return
	}
	//--->attach digital signature+data
	fakeCertificateWithSignature.Signature = base64.StdEncoding.EncodeToString(digitalSignatureByte)
	//! what is byte?
	//! what is stored theere ?
	//! how does signing work , how does encryption work here ? what is happening in decryption.
	//
	//
	//
	//=========================backend process below ========================
	signatureDecodedToBytes, err := base64.StdEncoding.DecodeString(fakeCertificateWithSignature.Signature)
	if err != nil {
		utils.LogErrorWithContext("Signature Decoding", err)
		return
	}
	rsa_verify_er := rsa.VerifyPKCS1v15(&rsaPublicKey, crypto.SHA256, hashedValue, signatureDecodedToBytes)

	if rsa_verify_er != nil {
		utils.LogErrorWithContext("Digital Signature Verification", rsa_verify_er)
		return
	}
	log.Printf("Digital Signature Verified ")

	//then parse into our  struct without signature , signature not necessary now .

	parsedBlock := ConvertintoModelsBlock(fakeCertificateWithSignature)
	er = InsertBlock(parsedBlock)
	if er != nil {
		utils.LogErrorWithContext("Insert Block", er)
		return
	}
	// pow + nonce

	//insert 1st data

	//print the data .

}
