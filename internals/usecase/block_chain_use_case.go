package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	err "project/package/errors"
	"project/package/utils/common"
	logger "project/package/utils/pkg"
)

type BlockChainUseCase struct {
	BlockChainRepo repository.IBlockChainRepository
	Service        service.Service
}

func NewBlockChainUseCase(blockChainRepository repository.IBlockChainRepository, service service.Service) *BlockChainUseCase {
	return &BlockChainUseCase{
		BlockChainRepo: blockChainRepository,
		Service:        service,
	}
}

func (uc *BlockChainUseCase) GetBlockChainLength() int {
	return uc.BlockChainRepo.GetBlockChainLength()
}
func (uc *BlockChainUseCase) InsertGenesisBlock() error {
	return uc.BlockChainRepo.InsertGenesisBlock()
}

func (b *BlockChainUseCase) GetCertificateData() (entity.CertificateData, error) {
	return b.BlockChainRepo.GetCertificateData()
}
func (uc *BlockChainUseCase) CompleteBlockFromCertificate(certificate entity.CertificateData) (*entity.Block, error) {
	logger.InitLogger()
	latestBlock, err := uc.BlockChainRepo.GetLatestBlock()
	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}
	totalCertificateDataLength, err := common.CalculateCertificateDataLength(latestBlock.CertificateData)
	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}
	if totalCertificateDataLength == 4 || latestBlock.Header.BlockNumber == 0 {
		previousHash := latestBlock.Header.CurrentHash
		latestBlock, err = uc.BlockChainRepo.GenerateNewBlock()
		latestBlock.Header.PreviousHash = previousHash

		if err != nil {
			logger.Logger.Infoln(err)
			return nil, err
		}
	}

	blockCertificateDataLength, err := common.CalculateCertificateDataLength(latestBlock.CertificateData)
	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}

	updatedBlockAfterCertificateInsertion, err := uc.BlockChainRepo.InsertCertificateIntoBlock(certificate, latestBlock, blockCertificateDataLength)
	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}

	merkleRoot, err := uc.Service.CalculateMerkleRoot(updatedBlockAfterCertificateInsertion.CertificateData)
	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}

	powStructureParam, err := entity.NewPowStructure(updatedBlockAfterCertificateInsertion.Header.BlockNumber, updatedBlockAfterCertificateInsertion.Header.PreviousHash, merkleRoot, updatedBlockAfterCertificateInsertion.Header.TimeStamp)
	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}
	calculatedNonce, currentHash, err := uc.Service.CalculatePOW(powStructureParam)
	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}

	if err = uc.BlockChainRepo.InsertIntoBlockChain(*updatedBlockAfterCertificateInsertion); err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}

	completeBlock, err := uc.BlockChainRepo.UpdateCurrentBlock(calculatedNonce, merkleRoot, currentHash, *updatedBlockAfterCertificateInsertion)

	if err != nil {
		logger.Logger.Infoln(err)
		return nil, err
	}

	// logger.Logger.Infoln(err)
	return completeBlock, nil
}

func (uc *BlockChainUseCase) UpsertBlockChain(toBeUsedBlock entity.Block) error {
	latestBlockInChain, er := uc.BlockChainRepo.GetLatestBlock()
	logger.InitLogger()
	if er != nil {
		logger.Logger.Infoln(er)
		return er
	}
	switch latestBlockInChain.Header.BlockNumber {
	case toBeUsedBlock.Header.BlockNumber - 1:
		return uc.BlockChainRepo.InsertIntoBlockChain(toBeUsedBlock)
	case toBeUsedBlock.Header.BlockNumber:
		return uc.BlockChainRepo.UpdateLatestBlockInBlockchain(toBeUsedBlock)

	}

	return err.ErrBlockNumberMismatch
}
