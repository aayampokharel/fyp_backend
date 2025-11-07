package service

// import (
// 	"crypto/sha256"
// 	"project/internals/domain/entity"
// 	err "project/package/errors"
// 	"project/package/utils/common"
// 	"time"
// )

// func InitializeBlockFields(block *entity.Block, previousHash string,lastCertificateLength int) (*entity.Block, error) {

// 	if lastCertificateLength>=4 || lastCertificateLength<0 {
// 		s.Logger.Errorln("[initialize_block_fields_service] Error: InitializeBlockFieldsService::", err.ErrArrayOutOfBound)

// 		return nil,err.ErrArrayOutOfBound

// 	}
// 	block.Header.TimeStamp = time.Now()
// 	block.Header.PreviousHash = previousHash
//    //  block.CertificateData[lastCertificateLength].CertificateHash=
//     block.CertificateData[lastCertificateLength].InstitutionFacultyID=//from db ,
//     // block.CertificateData[lastCertificateLength].CertificateID==common.GenerateUUID(16)
//     block.CertificateData[lastCertificateLength].InstitutionID=//from db
//     block.CertificateData[lastCertificateLength].PDFCategoryID=//from db
//     block.CertificateData[lastCertificateLength].Position=lastCertificateLength+1
//     block.CertificateData[lastCertificateLength].IssuerPublicKey=//from db
//     block.CertificateData[lastCertificateLength].createdAt=time.Now().Format(time.RFC3339)
//     block.CertificateData[lastCertificateLength].CertificateHash=//to be calculated .

// 	return block, nil
// }
