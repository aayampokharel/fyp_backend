package repository

import "project/internals/domain/entity"

type IBlockChainRepository interface {
	GetBlockChainLength() int
	GenerateNewBlock() (entity.Block, error)
	GetCertificateData() (entity.CertificateData, error)
	GetLatestBlock() (entity.Block, error)
	UpdateLatestBlockInBlockchain(block entity.Block) error
	InsertIntoBlockChain(block entity.Block) error
	InsertGenesisBlock() error
	GetTwoLatestBlocksInSlice() ([2]entity.Block, error)
	InsertCertificateIntoBlock(certificate entity.CertificateData, block entity.Block, certificateLength int) (*entity.Block, error)
	UpdateCurrentBlock(nonce int, merkleRoot string, currentHash string, block entity.Block) (*entity.Block, error)
	//! TransferBlockToOtherNodes(block entity.Block) error

}
