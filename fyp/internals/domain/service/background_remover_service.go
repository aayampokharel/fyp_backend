package service

import (
	"bytes"
	"log"
	"os/exec"
	err "project/package/errors"
)

func (s *Service) RemoveBackgroundService(imgBytes []byte) ([]byte, error) {
	cmd := exec.Command("python", "../scripts/bg_remover.py")
	cmd.Stdin = bytes.NewReader(imgBytes)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if er := cmd.Run(); er != nil {
		log.Println("Python error:", stderr.String())
		log.Println("Go error:", er)
		return nil, er
	}
	if out.Len() == 0 {
		return nil, err.ErrPythonScriptReturnedEmpty
	}
	return out.Bytes(), nil
}
