package entity

import (
	err "project/package/errors"
	logger "project/package/utils/pkg"
	"time"
)

type PowStructure struct {
	BlockNumber     int
	PreviousHash    string
	BlockMerkleRoot string
	BlockTimeStamp  time.Time
}

func NewPowStructure(blockNumber int, previousHash string, blockMerkleRoot string, blockTimeStamp time.Time) (PowStructure, error) {
	if blockNumber < 1 || previousHash == "" || blockMerkleRoot == "" || blockTimeStamp.IsZero() {
		logger.Logger.Infoln(blockNumber, previousHash, blockMerkleRoot)
		return PowStructure{}, err.ErrEmptyFields
	}

	return PowStructure{BlockNumber: blockNumber, PreviousHash: previousHash, BlockMerkleRoot: blockMerkleRoot, BlockTimeStamp: blockTimeStamp}, nil
}
