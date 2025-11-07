package service

import (
	"bytes"
	"html/template"
	err "project/package/errors"
)

func (s *Service) ParseAndExecute(templatePath string, data interface{}) (string, error) {
	tmpl, er := template.ParseFiles(templatePath)
	if er != nil {
		s.Logger.Errorln("[service] Error: ParseAndExecute::", er)
		return "", err.ErrFileParsing
	}
	var buf bytes.Buffer
	er = tmpl.Execute(&buf, data)
	if er != nil {
		s.Logger.Errorln("[service] Error: ParseAndExecute::", er)
		return "", err.ErrFileExecuting
	}

	return buf.String(), nil
}
