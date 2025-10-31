package delivery

import (
	"log"
	"project/internals/domain/entity"
	"project/internals/usecase"
	"project/package/utils/common"
)

type Controller struct {
	useCase usecase.BlockChainUseCase
}

func NewController(useCase usecase.BlockChainUseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (c *Controller) InsertNewCertificateData(request CreateCertificateDataRequest) entity.Response {
	blockChainLength := c.useCase.GetBlockChainLength()
	if blockChainLength == 0 {
		//mock data
		if err := c.useCase.InsertGenesisBlock(); err != nil {
			log.Println(err)
			return common.HandleErrorResponse(401, err.Error(), err)
		}
	}
	// certificateData, err := c.useCase.GetCertificateData()
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	for i := 0; i < len(request.CertificateData); i++ {

		certificateData := request.CertificateData[i]
		latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, err := c.useCase.CompleteBlockFromCertificate(certificateData)
		if err != nil {
			log.Println(err)
			return common.HandleErrorResponse(401, err.Error(), err)
		}
		var strNodeInfoMap map[int]string
		strNodeInfoMap, err = c.useCase.BroadcastNewBlock(newBlock)
		if err != nil {
			log.Println(err)
			return common.HandleErrorResponse(401, err.Error(), err)
		}
		log.Println("Acknowledgement from nodes: ", strNodeInfoMap)

		if err = c.useCase.UpsertBlockChain(*latestBlockFromChain, *newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength); err != nil {
			log.Println(err)
			return common.HandleErrorResponse(401, err.Error(), err)
		}

		// blockchain, _ := c.useCase.BlockChainRepo.GetBlockChain()
	}
	return common.HandleSuccessResponse(CreateCertificateResponse{
		Message: "All Certificates Inserted Successfully",
	})
}
