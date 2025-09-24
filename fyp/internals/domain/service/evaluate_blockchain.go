package service

import (
	"project/internals/domain/entity"
	"project/package/enum"
	err "project/package/errors"
)

func (s *Service) EvaluateBlockChain(latestBlockFromChain, newBlock entity.Block, latestBlockFromBlockChainCertificateLength, newBlockCertificateLength int, status enum.CREATEUPDATE) error {
	switch status {
	case enum.CREATE:
		if (latestBlockFromBlockChainCertificateLength != 4 && latestBlockFromChain.Header.BlockNumber != 0) || newBlockCertificateLength != 1 {
			return err.ErrBlockNumberMismatch
		}

	case enum.UPDATE:
		if latestBlockFromBlockChainCertificateLength == 4 || latestBlockFromBlockChainCertificateLength+1 != newBlockCertificateLength {
			return err.ErrBlockNumberMismatch
		}

	}
	return nil
}
