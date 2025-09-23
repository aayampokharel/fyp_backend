package source

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	err "project/package/errors"
	"time"
)

type BlockChainMemorySource struct {
	blockChain []entity.Block
}

func NewBlockChainMemorySource() *BlockChainMemorySource {
	return &BlockChainMemorySource{
		blockChain: make([]entity.Block, 0),
	}
}

var _ repository.IBlockChainRepository = (*BlockChainMemorySource)(nil)

func (b *BlockChainMemorySource) GetCertificateData() (entity.CertificateData, error) {
	return entity.CertificateData{
		ID:                 "12345",
		StudentName:        "aayam pokharel",
		UniversityName:     "TU",
		Degree:             "Bachelor's",           //struct
		College:            "St. Xavier's College", //structs
		CertificateDate:    time.Now(),
		Division:           "first",
		PrincipalSignature: "123456",
		TuApproval:         "true",
	}, nil
}

func (b *BlockChainMemorySource) GetLatestBlock() (entity.Block, error) {
	if len(b.blockChain) <= 0 {
		return entity.Block{}, err.ErrEmptyBlockChain
	}
	return b.blockChain[len(b.blockChain)-1], nil

}

func (b *BlockChainMemorySource) UpdateLatestBlockInBlockchain(block entity.Block) error {
	lengthOfBlockChain := len(b.blockChain)
	if lengthOfBlockChain <= 0 {
		return err.ErrEmptyBlockChain
	}
	if lengthOfBlockChain == 1 {
		return err.ErrGenesisBlockUpdate
	}

	if b.blockChain[lengthOfBlockChain-1].Header.BlockNumber == block.Header.BlockNumber {
		b.blockChain[len(b.blockChain)-1] = block
		return nil
	}
	return err.ErrBlockNumberMismatch

}

func (b *BlockChainMemorySource) InsertIntoBlockChain(block entity.Block) error {
	lengthOfBlockChain := len(b.blockChain)
	if lengthOfBlockChain == 0 {
		return err.ErrEmptyBlockChain
	}
	if b.blockChain[lengthOfBlockChain-1].Header.BlockNumber != block.Header.BlockNumber-1 {
		return err.ErrBlockNumberMismatch
	}
	b.blockChain = append(b.blockChain, block)
	return nil
}

func (b *BlockChainMemorySource) InsertGenesisBlock() error {
	lengthOfBlockChain := len(b.blockChain)
	if lengthOfBlockChain != 0 {
		return err.ErrGenesisBlockInsert
	}
	genesisBlock := entity.Block{Header: entity.Header{
		BlockNumber:  0,
		PreviousHash: "0",
		Nonce:        "0",
		CurrentHash:  "0",
		TimeStamp:    time.Now(),
		MerkleRoot:   "0",
	}}
	b.blockChain = append(b.blockChain, genesisBlock)
	return nil
}

func (b *BlockChainMemorySource) GetTwoLatestBlocksInSlice() ([2]entity.Block, error) {
	if len(b.blockChain) < 2 {
		return [2]entity.Block{}, err.ErrNotEnoughBlocks
	}
	return [2]entity.Block{b.blockChain[len(b.blockChain)-2], b.blockChain[len(b.blockChain)-1]}, nil
}
func (b *BlockChainMemorySource) InsertCertificateIntoBlock(certificate entity.CertificateData, block entity.Block, certificateLength int) (*entity.Block, error) {
	// certificateLength, err := common.CalculateCertificateDataLength(block.CertificateData)
	// if err != nil {
	// 	return nil, err
	// }
	block.CertificateData[certificateLength] = certificate
	return &block, nil
}
