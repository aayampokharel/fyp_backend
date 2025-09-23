package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/package/utils/common"
)

type BlockChainUseCase struct {
	BlockChainRepo repository.IBlockChainRepository
}

func NewBlockChainUseCase(blockChainRepository repository.IBlockChainRepository) *BlockChainUseCase {
	return &BlockChainUseCase{
		BlockChainRepo: blockChainRepository,
	}
}

func (uc *BlockChainUseCase) InsertGenesisBlock() error {
	return uc.BlockChainRepo.InsertGenesisBlock()
}

func (uc *BlockChainUseCase) InsertCertificateIntoBlockChain(certificate entity.CertificateData) error {
	latestBlock, err := uc.BlockChainRepo.GetLatestBlock()
	if err != nil {
		return err
	}

	blockCertificateDataLength, err := common.CalculateCertificateDataLength(latestBlock.CertificateData)
	if err != nil {
		return err
	}

	updatedBlock, err := uc.BlockChainRepo.InsertCertificateIntoBlock(certificate, latestBlock, blockCertificateDataLength)
	if err != nil {
		return err
	}

	if err = uc.BlockChainRepo.InsertIntoBlockChain(*updatedBlock); err != nil {
		return err
	}

	return nil
}
func (uc *BlockChainUseCase) AddCertificateToBlockchain(cert entity.CertificateData) error {
	return nil
}
