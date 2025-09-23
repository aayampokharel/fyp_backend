package common

import (
	"project/internals/domain/entity"
	err "project/package/errors"
)

func ValidateChain(blockChain []entity.Block) error {
	//should be called right before insertion or update
	if len(blockChain) == 0 {
		return err.ErrEmptyBlockChain
	}
	if blockChain[0].Header.BlockNumber != 0 {
		return err.ErrGenesisBlockMismatch
	}
	return nil
}
