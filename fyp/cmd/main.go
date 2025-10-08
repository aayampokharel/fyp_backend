package main

import (
	"fmt"
	"net/http"

	"project/constants"
	"project/internals/data/config"
	source "project/internals/data/memory"
	"project/internals/data/p2p"
	delivery "project/internals/delivery/blockchain"
	"project/internals/domain/service"
	"project/internals/usecase"
	"project/package/utils/common"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	currentPort := common.GetPort()
	logger.Logger.Infoln("[main] Info: Current Port::", *currentPort)
	env, err := config.NewEnv()
	if err != nil {
		logger.Logger.Errorw("[main] Error: Failed to load environment variables", zap.Error(err))
		return
	}

	ports := env.GetValueForKey(constants.NodePortsKey)
	fmt.Println("Node Ports:", ports)

	module := delivery.NewModule(source.NewBlockChainMemorySource(), p2p.NewNodeSource(ports))

	mux := http.NewServeMux()
	delivery.RegisterRoutes(mux, module)

	addr := fmt.Sprintf(":%d", *currentPort)
	fmt.Printf("🚀 Starting blockchain node on http://localhost%s\n", addr)

	memSource := source.NewBlockChainMemorySource()
	nodeSource := p2p.NewNodeSource(ports)
	service := service.NewService()
	// blockchainService := service.NewService()
	blockChainUseCase := usecase.NewBlockChainUseCase(memSource, nodeSource, service)

	go func() {
		for {
			er := blockChainUseCase.ReceiveBlockFromPeer(*currentPort)
			if er != nil {
				logger.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", zap.Error(er))
				fmt.Println("Error receiving block from peer:", er)
			}

		}
	}()
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("❌ Server failed:", err)
	}

	// controller := delivery.NewController(*blockChainUseCase)
	// controller.InsertNewCertificateData()
	// fmt.Print("end whooho ")
}

// func main() {

// 	controller.InsertNewCertificateData()
// 	controller.InsertNewCertificateData()
// 	controller.InsertNewCertificateData()
// 	controller.InsertNewCertificateData()
// 	controller.InsertNewCertificateData()
// 	controller.InsertNewCertificateData()
// 	// controller.InsertNewCertificateData()
// 	finalBlockChain, _ := controller.InsertNewCertificateData()
// 	fmt.Print("pretty JSON ")
// 	fmt.Println(len(finalBlockChain))
// 	fmt.Print("pretty JSON ")
// 	common.PrintPrettyJSON(finalBlockChain[0])
// 	common.PrintPrettyJSON(finalBlockChain[1])
// 	common.PrintPrettyJSON(finalBlockChain[2])

// 	fmt.Print("end whooho ")

// }
