package delivery

import (
	"log"
	"project/internals/usecase"
)

type Controller struct {
	useCase usecase.BlockChainUseCase
}

func NewController(useCase usecase.BlockChainUseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (c *Controller) InsertNewCertificateData() error {
	blockChainLength := c.useCase.GetBlockChainLength()
	if blockChainLength == 0 {
		if err := c.useCase.InsertGenesisBlock(); err != nil {
			log.Println(err)
			return err
		}
	}
	certificateData, err := c.useCase.GetCertificateData()
	if err != nil {
		log.Println(err)
		return err
	}
	finalBlock, err := c.useCase.CompleteBlockFromCertificate(certificateData)
	if err != nil {
		log.Println(err)
		return err
	}
	if err = c.useCase.UpsertBlockChain(*finalBlock); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
