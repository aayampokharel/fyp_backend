package common

import (
	"project/internals/domain/entity"
	err "project/package/errors"
)

func ValidateBlock(block entity.Block) error {
	if block.Header.BlockNumber < 1 {
		return err.ErrInvalidBlockNumber
	}

	if block.Header.CurrentHash == "" {
		return err.ErrInvalidHash
	}
	return nil
}
