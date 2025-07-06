package block

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"log"

	"github.com/aayampokharel/fyp/common"
	digitalsignature "github.com/aayampokharel/fyp/common/digital_signature"
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
	// fakeData := models.BlockWithSignature{
	// 	Signature: "",
	// 	BlockData: models.Block{
	// 		Header: models.Header{
	// 			BlockNumber:  0,
	// 			TimeStamp:    time.Now(),
	// 			PreviousHash: "",
	// 			Nonce:        "",
	// 			CurrentHash:  "",
	// 			MerkleRoot:   "",
	// 		},
	// 		CertificateData: [4]models.CertificateData{
	// 			{
	// 				ID:                 "",
	// 				StudentName:        "",
	// 				UniversityName:     "",
	// 				Degree:             "",
	// 				College:            "",
	// 				CertificateDate:    time.Now(),
	// 				Division:           "",
	// 				PrincipalSignature: "",
	// 				TuApproval:         "",
	// 			},
	// 		},
	// 	},
	// }
	//create a signature
	//-> hash

	hashedValue, err := HashBlock(fakeData.BlockData)
	if err != nil {
		utils.LogErrorWithContext("HashBlock", err)
		return
	}

	_, _, err = digitalsignature.GeneratePrivateKey()
	if err != nil {
		utils.LogErrorWithContext("GeneratePrivateKey", err)
		return
	}

	rsaPrivateKey, err := digitalsignature.LoadPrivateKey("/FYP/certs/private_key.pem")
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
	fakeData.Signature = base64.StdEncoding.EncodeToString(digitalSignatureByte)
	//! what is byte?
	//! what is stored theere ?
	//! how does signing work , how does encryption work here ? what is happening in decryption.
	signatureDecodedToBytes, err := base64.StdEncoding.DecodeString(fakeData.Signature)
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

	//then parse into out struct

	parsedBlock := ConvertintoModelsBlock(fakeData)
	// pow + nonce
	common.ProofOfWork(&parsedBlock.Header)
	merkelRootString, err := utils.CalculateMerkelRoot(parsedBlock.CertificateData[:])
	if err != nil {
		utils.LogErrorWithContext("Merkel Root Calculation", err)
		return
	}
	parsedBlock.Header.MerkleRoot = merkelRootString

	//insert 1st intial block
	utils.BlockChain = append(utils.BlockChain, utils.CreateGenesisBlock())
	//insert 1st data

	//print the data .

}
