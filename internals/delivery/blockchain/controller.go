package delivery

import (
	"log"
	"project/internals/domain/entity"
	"project/internals/usecase"
)

type Controller struct {
	useCase usecase.BlockChainUseCase
}

func NewController(useCase usecase.BlockChainUseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (c *Controller) InsertNewCertificateData() ([]entity.Block, error) {
	blockChainLength := c.useCase.GetBlockChainLength()
	if blockChainLength == 0 {
		if err := c.useCase.InsertGenesisBlock(); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	certificateData, err := c.useCase.GetCertificateData()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, err := c.useCase.CompleteBlockFromCertificate(certificateData)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// latestBlockFromChain, err := c.useCase.BlockChainRepo.GetLatestBlock()
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	// latestBlockFromChainCertificateLength, er := common.CalculateCertificateDataLength(latestBlockFromChain.CertificateData)
	// if er != nil {
	// 	log.Println(er)
	// 	return nil, er
	// }
	// newBlockCertificateLength, er := common.CalculateCertificateDataLength(latestBlockFromChain.CertificateData)
	// if er != nil {
	// 	log.Println(er)
	// 	return nil, er
	// }
	if err = c.useCase.UpsertBlockChain(*latestBlockFromChain, *newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength); err != nil {
		log.Println(err)
		return nil, err
	}

	// blockchain, err := c.useCase.BlockChainRepo.GetBlockChain()
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	// common.PrintPrettyJSON(blockchain)
	blockchain, _ := c.useCase.BlockChainRepo.GetBlockChain()
	return blockchain, nil
}
