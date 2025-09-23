package entity

import "time"

type PowStructure struct {
	BlockNumber     int
	PreviousHash    string
	BlockMerkleRoot string
	BlockTimeStamp  time.Time
}

func NewPowStructure(blockNumber int, previousHash string, blockMerkleRoot string, blockTimeStamp time.Time) PowStructure {
	return PowStructure{BlockNumber: blockNumber, PreviousHash: previousHash, BlockMerkleRoot: blockMerkleRoot, BlockTimeStamp: blockTimeStamp}
}
