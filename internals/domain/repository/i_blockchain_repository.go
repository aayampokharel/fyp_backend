package repository

import "project/internals/domain/entity"

type IBlockChainRepository interface {
	GetCertificateData() (entity.CertificateData, error)
	GetLatestBlock() (entity.Block, error)
	UpdateLatestBlockInBlockchain(block entity.Block) error
	InsertIntoBlockChain(block entity.Block) error
	InsertGenesisBlock() error
	GetTwoLatestBlocksInSlice() ([2]entity.Block, error)
	InsertCertificateIntoBlock(certificate entity.CertificateData, block entity.Block, certificateLength int) (*entity.Block, error)
	//! TransferBlockToOtherNodes(block entity.Block) error

}
