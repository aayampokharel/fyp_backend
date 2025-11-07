package memory_source

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	err "project/package/errors"
	"project/package/seed"
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

func (b *BlockChainMemorySource) GenerateNewBlock(blockNumber int, previousHash string) (entity.Block, error) {
	return entity.Block{
		Header: entity.Header{
			BlockNumber:  blockNumber,
			PreviousHash: previousHash,
			TimeStamp:    time.Now(),
		},
	}, nil

}

func (b *BlockChainMemorySource) GetAnyBlockHeader(length int) (entity.Header, error) {
	return b.blockChain[length].Header, nil
}
func (b *BlockChainMemorySource) GetBlockChain() ([]entity.Block, error) {

	return b.blockChain, nil

}
func (b *BlockChainMemorySource) GetCertificateData() (*entity.CertificateData, error) {
	return seed.GenerateRandomCertificateData(), nil
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
	b.blockChain[len(b.blockChain)-1] = block

	return nil
}

func (b *BlockChainMemorySource) InsertIntoBlockChain(block entity.Block) error {

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
func (b *BlockChainMemorySource) InsertCertificateIntoBlock(certificate *entity.CertificateData, block entity.Block) (*entity.Block, int, error) {
	nextIndex, err := common.CalculateCertificateDataLength(block.CertificateData)
	if err != nil {
		return nil, -1, err
	}
	block.CertificateData[nextIndex] = *certificate
	b.logger.Infoln("[block_chain_memory_source] Info: InsertCertificateIntoBlock::", block)
	return &block, nextIndex + 1, nil
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
	// common.PrintPrettyJSON(block)
	return &block, nil
}

func (b *BlockChainMemorySource) ReceiveFromPeer(currentPort string) error {
	return nil
}
func (b *BlockChainMemorySource) GetInfoFromPdfFilesCategories(categoryID string) (*entity.PDFFileCategoryEntity, error) {
	return nil, nil
}
func (b *BlockChainMemorySource) GetCertificateDataList(institutionID, institutionFacultyID, categoryID string) ([]entity.CertificateData, error) {

	var certificateDataList []entity.CertificateData
	for _, block := range b.blockChain {
		for _, certificate := range block.CertificateData {
			if certificate.InstitutionID == institutionID && certificate.InstitutionFacultyID == institutionFacultyID && certificate.PDFCategoryID == categoryID {
				certificateDataList = append(certificateDataList, certificate)
			}
		}
	}
	return certificateDataList, nil
}
