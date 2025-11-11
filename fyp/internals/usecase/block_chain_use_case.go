package usecase

import (
	"project/constants"
	"project/internals/data/config"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/package/enum"
	errorz "project/package/errors"
	"project/package/utils/common"
)

type BlockChainUseCase struct {
	BlockChainRepo repository.IBlockChainRepository
	NodeRepo       repository.INodeRepository
	SqlRepo        repository.ISqlRepository
	Service        service.Service
}

func NewBlockChainUseCase(blockChainRepository repository.IBlockChainRepository, nodeRepository repository.INodeRepository, sqlRepo repository.ISqlRepository, service service.Service) *BlockChainUseCase {
	return &BlockChainUseCase{
		BlockChainRepo: blockChainRepository,
		NodeRepo:       nodeRepository,
		SqlRepo:        sqlRepo,
		Service:        service,
	}
}

func (uc *BlockChainUseCase) InsertData(block *entity.Block) error {
	return uc.BlockChainRepo.InsertIntoBlockChain(*block)
}

func (uc *BlockChainUseCase) GetBlockChainLength() int {
	uc.Service.Logger.Infoln("[block_chain_use_case] Info: GetBlockChainLength::", uc.BlockChainRepo.GetBlockChainLength())
	return uc.BlockChainRepo.GetBlockChainLength()
}
func (uc *BlockChainUseCase) InsertGenesisBlock() error {
	return uc.BlockChainRepo.InsertGenesisBlock()
}

func (b *BlockChainUseCase) GetCertificateData() (*entity.CertificateData, error) {
	return b.BlockChainRepo.GetCertificateData()
}
func (uc *BlockChainUseCase) CompleteBlockFromCertificate(certificate *entity.CertificateData) (latestBlockFromChain, newBlock *entity.Block, latestBlockFromChainCertificateLength, newBlockCertificateLength int, err error) {
	//deals from getting certificate data to complete inserion at its proper position in the block.
	//auto handles to extract eiter new block or old block if space available
	latestBlock, err := uc.BlockChainRepo.GetLatestBlock()
	if err != nil {
		uc.Service.Logger.Infoln(err)
		return nil, nil, 0, 0, err
	}
	totalCertificateDataLength, err := common.CalculateCertificateDataLength(latestBlock.CertificateData)
	if err != nil {
		uc.Service.Logger.Infoln(err)
		return nil, nil, 0, 0, err
	}
	if totalCertificateDataLength == 4 || latestBlock.Header.BlockNumber == 0 {
		//new block relatin logic
		previousHash := latestBlock.Header.CurrentHash
		blockNumber := latestBlock.Header.BlockNumber + 1
		latestBlock, err = uc.BlockChainRepo.GenerateNewBlock(blockNumber, previousHash)
		if err != nil {
			uc.Service.Logger.Infoln(err)
			return nil, nil, 0, 0, err
		}
		//latestBlock.Header.PreviousHash = previousHash
		//for updated and completely new block
		totalCertificateDataLength = 0

	}

	updatedBlockAfterCertificateInsertion, totalCertificateDataLength, err := uc.BlockChainRepo.InsertCertificateIntoBlock(certificate, latestBlock)
	if err != nil {
		uc.Service.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}

	merkleRoot, err := uc.Service.CalculateMerkleRoot(updatedBlockAfterCertificateInsertion.CertificateData)
	if err != nil {
		uc.Service.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}

	powStructureParam, err := entity.NewPowStructure(updatedBlockAfterCertificateInsertion.Header.BlockNumber, updatedBlockAfterCertificateInsertion.Header.PreviousHash, merkleRoot, updatedBlockAfterCertificateInsertion.Header.TimeStamp)
	if err != nil {
		uc.Service.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}
	//!(SINGLETON PATTERN) for ENV
	//~instead of loading env always like this there should be some other solution as always reading env like this , my env file may furter lengthen more later , so initialize once , use in all other places
	env, err := config.NewEnv()
	if err != nil {
		uc.Service.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}
	//!-----
	powRuleString := env.GetValueForKey(constants.PowNumberRuleKey)

	calculatedNonce, currentHash, err := uc.Service.CalculatePOW(powStructureParam, powRuleString)
	if err != nil {
		uc.Service.Logger.Errorln(err)
		return nil, nil, 0, 0, err
	}

	lengthOfBlockChain := uc.BlockChainRepo.GetBlockChainLength()
	latestBlockFromBlockChain, _ := uc.BlockChainRepo.GetLatestBlock()
	if lengthOfBlockChain == 0 {
		uc.Service.Logger.Errorln("[block_chain_use_case] Error: CompleteBlockFromCertificate::", errorz.ErrEmptyBlockChain)
		return nil, nil, 0, 0, errorz.ErrEmptyBlockChain
	}

	latestBlockFromBlockChainCertificateLength, _ := common.CalculateCertificateDataLength(latestBlockFromBlockChain.CertificateData)

	if latestBlockFromBlockChain.Header.BlockNumber != latestBlock.Header.BlockNumber-1 && latestBlockFromBlockChainCertificateLength == 4 {
		uc.Service.Logger.Errorln("[block_chain_use_case] Error: CompleteBlockFromCertificate::", errorz.ErrBlockNumberMismatch)
		return nil, nil, 0, 0, errorz.ErrBlockNumberMismatch
	}

	if latestBlock.Header.BlockNumber != latestBlockFromBlockChain.Header.BlockNumber && latestBlockFromBlockChainCertificateLength != 4 && latestBlockFromBlockChain.Header.BlockNumber != 0 {
		uc.Service.Logger.Errorln("[block_chain_memory_source] Error: CompleteBlockFromCertificate::", errorz.ErrBlockNumberMismatch, latestBlock.Header.BlockNumber, latestBlockFromBlockChain.Header.BlockNumber)
		return nil, nil, 0, 0, errorz.ErrBlockNumberMismatch
	}

	completeBlock, err := uc.BlockChainRepo.UpdateCurrentBlock(calculatedNonce, merkleRoot, currentHash, *updatedBlockAfterCertificateInsertion)
	////completeBlock is the one to be transported to other nodes .
	if err != nil {
		uc.Service.Logger.Infoln(err)
		return nil, nil, 0, 0, err
	}
	return &latestBlockFromBlockChain, completeBlock, latestBlockFromBlockChainCertificateLength, totalCertificateDataLength, nil
}

func (uc *BlockChainUseCase) BroadcastNewBlock(completeBlock *entity.Block) (map[int]string, error) {
	return uc.NodeRepo.SendBlockToPeer(*completeBlock, common.GetMappedTCPPort())
}

func (uc *BlockChainUseCase) ReceiveBlockFromPeer(currentTCPPort int) error {
	receivedUpdatedBlock, er := uc.NodeRepo.ReceiveBlockFromPeer(currentTCPPort)
	if er != nil {
		uc.Service.Logger.Errorln(er)
		return er
	}
	common.PrintPrettyJSON(receivedUpdatedBlock)
	uc.Service.Logger.Debugw("[block_chain_use_case] Debug: ReceiveBlockFromPeer:: Received Block from peer:", "receivedUpdatedBlock", "receivedUpdatedBlock")

	blockChainLength := uc.BlockChainRepo.GetBlockChainLength()
	if blockChainLength == 0 {
		if err := uc.InsertGenesisBlock(); err != nil {
			uc.Service.Logger.Errorln(err)
			return err
		}
	}
	latestBlockFromChain, err := uc.BlockChainRepo.GetLatestBlock()
	if err != nil {
		uc.Service.Logger.Errorln(err)
		return err
	}
	latestBlockFromChainCertificateLength, _ := common.CalculateCertificateDataLength(latestBlockFromChain.CertificateData)
	receivedBlockCertificateLength, _ := common.CalculateCertificateDataLength(receivedUpdatedBlock.CertificateData)

	uc.UpsertBlockChain(latestBlockFromChain, *receivedUpdatedBlock, latestBlockFromChainCertificateLength, receivedBlockCertificateLength)
	return nil
}

func (uc *BlockChainUseCase) UpsertBlockChain(latestBlockFromChain, newBlock entity.Block, latestBlockFromChainCertificateLength, newBlockCertificateLength int) error {
	//This helps in recognizing whether a block needs to replace last block from chain(update) or be inserted as new block in chain.
	switch latestBlockFromChain.Header.BlockNumber {

	case newBlock.Header.BlockNumber - 1:
		if erz := uc.Service.EvaluateBlockChain(latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, enum.CREATE); erz != nil {
			uc.Service.Logger.Errorln(erz)
			return erz
		}
		return uc.BlockChainRepo.InsertIntoBlockChain(newBlock)

	case newBlock.Header.BlockNumber:
		if erz := uc.Service.EvaluateBlockChain(latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, enum.UPDATE); erz != nil {
			uc.Service.Logger.Errorln(erz)
			return erz
		}
		return uc.BlockChainRepo.UpdateLatestBlockInBlockchain(newBlock)

	}

	return errorz.ErrBlockNumberMismatch
}

func (uc *BlockChainUseCase) GetBlockChain() int {
	uc.Service.Logger.Infoln("[get_block_chain] Info: GetBlockChain::", uc.BlockChainRepo.GetBlockChainLength())
	return uc.BlockChainRepo.GetBlockChainLength()
}
func (uc *BlockChainUseCase) GetCertificateDataListUseCase(institutionID, institutionFacultyID, categoryID string) ([]entity.CertificateData, error) {
	return uc.BlockChainRepo.GetCertificateDataList(institutionID, institutionFacultyID, categoryID)
}

func (uc *BlockChainUseCase) SendPBFTMessageToPeer(pbftMessage entity.PBFTMessage) {
	env, err := config.NewEnv()
	if err != nil {
		uc.Service.Logger.Errorln(err)
		return
	}

	leaderNode := env.GetValueForKey(constants.PbftLeaderNode)
	leaderNodeInt, er := common.ConvertToInt(leaderNode)
	if er != nil {
		uc.Service.Logger.Errorln(er)
		return
	}
	uc.NodeRepo.SendPBFTMessageToPeer(pbftMessage, leaderNodeInt, common.GetMappedTCPPBFTPort())
}

func (uc *BlockChainUseCase) ReceivePBFTMessageFromPeers(currentTCPPort int) (*entity.PBFTMessage, error) {
	leaderPort := common.GetLeaderPort()
	if leaderPort == nil {
		return nil, errorz.ErrLeaderPort
	}
	return uc.NodeRepo.ReceivePBFTMessageToPeer(common.GetMappedTCPPBFTPort(), *leaderPort)
}
