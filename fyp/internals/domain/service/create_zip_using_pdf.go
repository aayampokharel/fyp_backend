package service

import (
	"archive/zip"
	"bytes"
	"project/internals/domain/entity"
	err "project/package/errors"
)

func (s *Service) CreateZipUsingPDF(listOfPDFs []entity.PDFFileEntity) ([]byte, error) {
	//create buffer
	zipBuffer := bytes.Buffer{}
	//create file
	zipWriter := zip.NewWriter(&zipBuffer)
	for _, val := range listOfPDFs {
		fileWriter, er := zipWriter.Create(val.FileName)
		if er != nil {
			return nil, err.ErrWritingZip
		}
		fileWriter.Write(val.PDFData)
	}
	er := zipWriter.Close()
	if er != nil {
		return nil, err.ErrClosingZipWriter
	}
	return zipBuffer.Bytes(), nil
}
