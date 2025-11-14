package entity

import (
	err "project/package/errors"
	logger "project/package/utils/pkg"
	"time"
)

type PowStructure struct {
	BlockNumber           int
	PreviousPOWPuzzleHash string
	BlockMerkleRoot       string
	BlockTimeStamp        time.Time
}

func NewPowStructure(blockNumber int, previousPowPuzzleHash string, blockMerkleRoot string, blockTimeStamp time.Time) (PowStructure, error) {
	if blockNumber < 1 || previousPowPuzzleHash == "" || blockMerkleRoot == "" || blockTimeStamp.IsZero() {
		logger.Logger.Infoln(blockNumber, previousPowPuzzleHash, blockMerkleRoot)
		return PowStructure{}, err.ErrEmptyFields
	}

	return PowStructure{BlockNumber: blockNumber, PreviousPOWPuzzleHash: previousPowPuzzleHash, BlockMerkleRoot: blockMerkleRoot, BlockTimeStamp: blockTimeStamp}, nil
}
