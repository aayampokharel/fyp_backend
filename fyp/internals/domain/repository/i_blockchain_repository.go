package repository

import "project/internals/domain/entity"

type IBlockChainRepository interface {
	GetBlockChainLength() int

	GetBlockChain() ([]entity.Block, error)

	GenerateNewBlock(blockNumber int, previousHash string) (entity.Block, error)

	GetCertificateData() (*entity.CertificateData, error)

	GetLatestBlock() (entity.Block, error)

	UpdateLatestBlockInBlockchain(block entity.Block) error

	InsertIntoBlockChain(block entity.Block) error

	InsertGenesisBlock() error

	GetTwoLatestBlocksInSlice() ([2]entity.Block, error)

	InsertCertificateIntoBlock(certificate *entity.CertificateData, block entity.Block) (*entity.Block, int, error)

	UpdateCurrentBlock(nonce int, merkleRoot string, currentHash string, block entity.Block) (*entity.Block, error)

	GetAnyBlockHeader(length int) (entity.Header, error)
	ReceiveFromPeer(currentPort string) error
}
