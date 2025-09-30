package main

import (
	"fmt"
	"project/internals/data/source"
	delivery "project/internals/delivery/blockchain"
	"project/internals/domain/service"
	"project/internals/usecase"
	"project/package/utils/common"
	logger "project/package/utils/pkg"
)

func main() {

	logger.InitLogger()
	memorySource := source.NewBlockChainMemorySource()
	blockchainService := service.Service{}
	usecase := usecase.NewBlockChainUseCase(memorySource, blockchainService)
	controller := delivery.NewController(*usecase)
	controller.InsertNewCertificateData()
	controller.InsertNewCertificateData()
	controller.InsertNewCertificateData()
	controller.InsertNewCertificateData()
	controller.InsertNewCertificateData()
	controller.InsertNewCertificateData()
	controller.InsertNewCertificateData()
	// controller.InsertNewCertificateData()
	finalBlockChain, _ := controller.InsertNewCertificateData()
	fmt.Print("pretty JSON ")
	fmt.Println(len(finalBlockChain))
	fmt.Print("pretty JSON ")
	common.PrintPrettyJSON(finalBlockChain[0])
	common.PrintPrettyJSON(finalBlockChain[1])
	common.PrintPrettyJSON(finalBlockChain[2])

	fmt.Print("end whooho ")

}
