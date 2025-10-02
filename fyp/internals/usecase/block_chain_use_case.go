package usecase

import (
	"project/constants"
	"project/internals/data/config"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/package/enum"
	err "project/package/errors"
	errorz "project/package/errors"
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
func (uc *BlockChainUseCase) CompleteBlockFromCertificate(certificate entity.CertificateData) (latestBlockFromChain, newBlock *entity.Block, latestBlockFromChainCertificateLength, newBlockCertificateLength int, err error) {
	latestBlock, err := uc.BlockChainRepo.GetLatestBlock()
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, nil, 0, 0, err
	}
	totalCertificateDataLength, err := common.CalculateCertificateDataLength(latestBlock.CertificateData)
	if err != nil {
		uc.Logger.Infoln(err)
		return nil, nil, 0, 0, err
	}
	if totalCertificateDataLength == 4 || latestBlock.Header.BlockNumber == 0 {
		previousHash := latestBlock.Header.CurrentHash
		blockNumber := latestBlock.Header.BlockNumber + 1
		latestBlock, err = uc.BlockChainRepo.GenerateNewBlock(blockNumber, previousHash)
		if err != nil {
			uc.Logger.Infoln(err)
			return nil, nil, 0, 0, err
		}
		latestBlock.Header.PreviousHash = previousHash
		//for updated and completely new block
		totalCertificateDataLength = 0

	}

	updatedBlockAfterCertificateInsertion, totalCertificateDataLength, err := uc.BlockChainRepo.InsertCertificateIntoBlock(certificate, latestBlock)
	if err != nil {
		uc.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}

	merkleRoot, err := uc.Service.CalculateMerkleRoot(updatedBlockAfterCertificateInsertion.CertificateData)
	if err != nil {
		uc.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}

	powStructureParam, err := entity.NewPowStructure(updatedBlockAfterCertificateInsertion.Header.BlockNumber, updatedBlockAfterCertificateInsertion.Header.PreviousHash, merkleRoot, updatedBlockAfterCertificateInsertion.Header.TimeStamp)
	if err != nil {
		uc.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}
	env, err := config.NewEnv()
	if err != nil {
		uc.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}
	powRuleString := env.GetValueForKey(constants.PowNumberRuleKey)

	calculatedNonce, currentHash, err := uc.Service.CalculatePOW(powStructureParam, powRuleString)
	if err != nil {
		uc.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}

	lengthOfBlockChain := uc.BlockChainRepo.GetBlockChainLength()
	latestBlockFromBlockChain, _ := uc.BlockChainRepo.GetLatestBlock()
	if lengthOfBlockChain == 0 {
		uc.Logger.Errorln("[block_chain_use_case] Error: CompleteBlockFromCertificate::", errorz.ErrEmptyBlockChain)
		return nil, nil, 0, 0, errorz.ErrEmptyBlockChain
	}

	latestBlockFromBlockChainCertificateLength, _ := common.CalculateCertificateDataLength(latestBlockFromBlockChain.CertificateData)

	if latestBlockFromBlockChain.Header.BlockNumber != latestBlock.Header.BlockNumber-1 && latestBlockFromBlockChainCertificateLength == 4 {
		uc.Logger.Errorln("[block_chain_use_case] Error: CompleteBlockFromCertificate::", errorz.ErrBlockNumberMismatch)
		return nil, nil, 0, 0, errorz.ErrBlockNumberMismatch
	}

	if latestBlock.Header.BlockNumber != latestBlockFromBlockChain.Header.BlockNumber && latestBlockFromBlockChainCertificateLength != 4 && latestBlockFromBlockChain.Header.BlockNumber != 0 {
		uc.Logger.Errorln("[block_chain_memory_source] Error: CompleteBlockFromCertificate::", errorz.ErrBlockNumberMismatch, latestBlock.Header.BlockNumber, latestBlockFromBlockChain.Header.BlockNumber)
		return nil, nil, 0, 0, errorz.ErrBlockNumberMismatch
	}

	completeBlock, err := uc.BlockChainRepo.UpdateCurrentBlock(calculatedNonce, merkleRoot, currentHash, *updatedBlockAfterCertificateInsertion)

	if err != nil {
		uc.Logger.Infoln(err)
		return nil, nil, 0, 0, err
	}
	return &latestBlockFromBlockChain, completeBlock, latestBlockFromBlockChainCertificateLength, totalCertificateDataLength, nil
}

func (uc *BlockChainUseCase) UpsertBlockChain(latestBlockFromChain, newBlock entity.Block, latestBlockFromChainCertificateLength, newBlockCertificateLength int) error {

	switch latestBlockFromChain.Header.BlockNumber {

	case newBlock.Header.BlockNumber - 1:
		if erz := uc.Service.EvaluateBlockChain(latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, enum.CREATE); erz != nil {
			uc.Logger.Errorln(erz)
			return erz
		}
		return uc.BlockChainRepo.InsertIntoBlockChain(newBlock)

	case newBlock.Header.BlockNumber:
		if erz := uc.Service.EvaluateBlockChain(latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, enum.UPDATE); erz != nil {
			uc.Logger.Errorln(erz)
			return erz
		}
		return uc.BlockChainRepo.UpdateLatestBlockInBlockchain(newBlock)

	}

	return err.ErrBlockNumberMismatch
}
