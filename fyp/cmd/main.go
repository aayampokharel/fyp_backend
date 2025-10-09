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
	tcpPort := *currentPort + 1000 // e.g., 8001 -> 9001

	// logger.Logger.Infoln("[main] Info: Current Port::", *currentPort)
	env, err := config.NewEnv()
	if err != nil {
		logger.Logger.Errorw("[main] Error: Failed to load environment variables", zap.Error(err))
		return
	}

	peerPorts := env.GetValueForKey(constants.TCPPortsKey)
	fmt.Println("Peer Ports:", peerPorts)

	nodeSource := p2p.NewNodeSource(peerPorts)
	module := delivery.NewModule(source.NewBlockChainMemorySource(), nodeSource)

	mux := http.NewServeMux()
	delivery.RegisterRoutes(mux, module)

	addr := fmt.Sprintf(":%d", *currentPort)
	fmt.Printf("🚀 Starting blockchain node on http://localhost%s\n", addr)

	memSource := source.NewBlockChainMemorySource()
	service := service.NewService()
	// blockchainService := service.NewService()
	blockChainUseCase := usecase.NewBlockChainUseCase(memSource, nodeSource, service)

	go func() {
		for {
			er := blockChainUseCase.ReceiveBlockFromPeer(tcpPort)
			if er != nil {
				logger.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", zap.Error(er))
				fmt.Println("Error receiving block from peer:", er)
			}

		}
	}()
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("❌ Server failed:", err)
	}

}
