package service

import (
	"project/internals/domain/entity"
	err "project/package/errors"
	"project/package/utils/common"
)

func (s *Service) GenerateCertificateData(cert entity.HashableData) (string, error) {
	if cert.CertificateID == "" || cert.StudentID == "" || cert.StudentName == "" ||
		cert.InstitutionID == "" || cert.InstitutionFacultyID == "" ||
		cert.FacultyPublicKey == "" || cert.IssueDate.IsZero() {
		return "", err.ErrEmptyString
	}

	hash, _, er := common.HashData(cert)
	if er != nil {
		return "", er
	}
	return hash, nil
}
