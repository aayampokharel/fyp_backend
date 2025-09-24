package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	err "project/package/errors"
	"project/package/utils/common"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

type BlockChainUseCase struct {
	BlockChainRepo repository.IBlockChainRepository
	Service        service.Service
	Logger         *zap.SugaredLogger
}

func NewBlockChainUseCase(blockChainRepository repository.IBlockChainRepository, service service.Service) *BlockChainUseCase {
	return &BlockChainUseCase{
		BlockChainRepo: blockChainRepository,
		Service:        service,
		Logger:         logger.Logger,
	}
}

func (uc *BlockChainUseCase) GetBlockChainLength() int {
	uc.Logger.Infoln("[block_chain_use_case] Info: GetBlockChainLength::", uc.BlockChainRepo.GetBlockChainLength())
	return uc.BlockChainRepo.GetBlockChainLength()
}
func (uc *BlockChainUseCase) InsertGenesisBlock() error {
	return uc.BlockChainRepo.InsertGenesisBlock()
}

func (b *BlockChainUseCase) GetCertificateData() (entity.CertificateData, error) {
	return b.BlockChainRepo.GetCertificateData()
}
func (uc *BlockChainUseCase) CompleteBlockFromCertificate(certificate entity.CertificateData) (*entity.Block, error) {
	latestBlock, err := uc.BlockChainRepo.GetLatestBlock()
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}
	totalCertificateDataLength, err := common.CalculateCertificateDataLength(latestBlock.CertificateData)
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}
	if totalCertificateDataLength == 4 || latestBlock.Header.BlockNumber == 0 {
		previousHash := latestBlock.Header.CurrentHash
		latestBlock, err = uc.BlockChainRepo.GenerateNewBlock()
		latestBlock.Header.PreviousHash = previousHash

		if err != nil {
			uc.Logger.Infoln(err)
			return nil, err
		}
	}

	blockCertificateDataLength, err := common.CalculateCertificateDataLength(latestBlock.CertificateData)
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}

	updatedBlockAfterCertificateInsertion, err := uc.BlockChainRepo.InsertCertificateIntoBlock(certificate, latestBlock, blockCertificateDataLength)
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}

	merkleRoot, err := uc.Service.CalculateMerkleRoot(updatedBlockAfterCertificateInsertion.CertificateData)
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}

	powStructureParam, err := entity.NewPowStructure(updatedBlockAfterCertificateInsertion.Header.BlockNumber, updatedBlockAfterCertificateInsertion.Header.PreviousHash, merkleRoot, updatedBlockAfterCertificateInsertion.Header.TimeStamp)
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}
	calculatedNonce, currentHash, err := uc.Service.CalculatePOW(powStructureParam)
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}

	if err = uc.BlockChainRepo.InsertIntoBlockChain(*updatedBlockAfterCertificateInsertion); err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}

	completeBlock, err := uc.BlockChainRepo.UpdateCurrentBlock(calculatedNonce, merkleRoot, currentHash, *updatedBlockAfterCertificateInsertion)

	if err != nil {
		uc.Logger.Infoln(err)
		return nil, err
	}
	return completeBlock, nil
}

func (uc *BlockChainUseCase) UpsertBlockChain(toBeUsedBlock entity.Block) error {
	latestBlockInChain, er := uc.BlockChainRepo.GetLatestBlock()
	if er != nil {
		uc.Logger.Infoln(er)
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
