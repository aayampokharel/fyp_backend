package source

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	err "project/package/errors"
	"project/package/utils/common"
	logger "project/package/utils/pkg"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type BlockChainMemorySource struct {
	blockChain []entity.Block
	logger     *zap.SugaredLogger
}

func NewBlockChainMemorySource() *BlockChainMemorySource {
	return &BlockChainMemorySource{
		logger:     logger.Logger,
		blockChain: make([]entity.Block, 0),
	}
}

var _ repository.IBlockChainRepository = (*BlockChainMemorySource)(nil)

func (b *BlockChainMemorySource) GenerateNewBlock() (entity.Block, error) {

	blockNumber := len(b.blockChain)
	previousHash := b.blockChain[len(b.blockChain)-1].Header.CurrentHash
	block := entity.Block{}
	block.Header.BlockNumber = blockNumber
	block.Header.PreviousHash = previousHash
	block.Header.TimeStamp = time.Now()
	b.logger.Infoln("[block_chain_memory_source] Info: GenerateNewBLock::", block.Header)
	return block, nil

}
func (b *BlockChainMemorySource) GetCertificateData() (entity.CertificateData, error) {
	return entity.CertificateData{
		ID:                 "123456",
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

		b.logger.Errorln("[block_chain_memory_source] Error: GetLatestBlock::", len(b.blockChain))
		return entity.Block{}, err.ErrEmptyBlockChain
	}
	return b.blockChain[len(b.blockChain)-1], nil

}

func (b *BlockChainMemorySource) UpdateLatestBlockInBlockchain(block entity.Block) error {
	lengthOfBlockChain := len(b.blockChain)
	if lengthOfBlockChain <= 0 {
		b.logger.Errorln("[block_chain_memory_source] Error: UpdateLatestBlockInBlockChain::", len(b.blockChain))
		return err.ErrEmptyBlockChain
	}

	if lengthOfBlockChain == 1 {
		b.logger.Errorln("[block_chain_memory_source] Error: UpdateLatestBlockInBlockChain::", len(b.blockChain))
		return err.ErrGenesisBlockUpdate
	}

	if b.blockChain[lengthOfBlockChain-1].Header.BlockNumber == block.Header.BlockNumber {
		b.blockChain[len(b.blockChain)-1] = block
		return nil
	}
	b.logger.Errorln("[block_chain_memory_source] Error: UpdateLatestBlockInBlockChain::", err.ErrBlockNumberMismatch)
	return err.ErrBlockNumberMismatch

}

func (b *BlockChainMemorySource) InsertIntoBlockChain(block entity.Block) error {
	lengthOfBlockChain := len(b.blockChain)
	if lengthOfBlockChain == 0 {
		b.logger.Errorln("[block_chain_memory_source] Error: InsertIntoBlockChain::", err.ErrEmptyBlockChain)
		return err.ErrEmptyBlockChain
	}
	if b.blockChain[lengthOfBlockChain-1].Header.BlockNumber != block.Header.BlockNumber-1 {
		b.logger.Errorln("[block_chain_memory_source] Error: InsertIntoBlockChain::", err.ErrBlockNumberMismatch)
		return err.ErrBlockNumberMismatch
	}
	b.logger.Infoln("[block_chain_memory_source] Info: InsertIntoBlockChain::", block)
	b.blockChain = append(b.blockChain, block)
	return nil
}

func (b *BlockChainMemorySource) InsertGenesisBlock() error {
	lengthOfBlockChain := len(b.blockChain)
	if lengthOfBlockChain != 0 {
		b.logger.Errorln("[block_chain_memory_source] Error: InsertGenesisBlock::", err.ErrGenesisBlockInsert)
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
		b.logger.Errorln("[block_chain_memory_source] Error: GetTwoLatestBlocksInSlice::", err.ErrNotEnoughBlocks)
		return [2]entity.Block{}, err.ErrNotEnoughBlocks
	}
	return [2]entity.Block{b.blockChain[len(b.blockChain)-2], b.blockChain[len(b.blockChain)-1]}, nil
}
func (b *BlockChainMemorySource) InsertCertificateIntoBlock(certificate entity.CertificateData, block entity.Block, certificateLength int) (*entity.Block, error) {
	block.CertificateData[certificateLength] = certificate

	b.logger.Infoln("[block_chain_memory_source] Info: InsertCertificateIntoBlock::", block)
	return &block, nil
}

func (b *BlockChainMemorySource) GetBlockChainLength() int {
	return len(b.blockChain)
}

func (b *BlockChainMemorySource) UpdateCurrentBlock(nonce int, merkleRoot string, currentHash string, block entity.Block) (*entity.Block, error) {

	if nonce < 0 || merkleRoot == "" || currentHash == "" {
		b.logger.Errorln("[block_chain_memory_source] Error: UpdateCurrentBlock::", err.ErrEmptyFields)
		return nil, err.ErrEmptyFields

	}
	block.Header.Nonce = strconv.Itoa(nonce)
	block.Header.CurrentHash = currentHash
	block.Header.MerkleRoot = merkleRoot

	b.logger.Infoln("[block_chain_memory_source] Info: UpdateCurrentBlock::")
	common.PrintPrettyJSON(block)
	return &block, nil
}
