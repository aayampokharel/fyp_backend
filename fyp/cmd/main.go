package main

import (
	"fmt"
	"net/http"

	"project/constants"
	"project/internals/data/config"
	source "project/internals/data/data_source/memory"
	"project/internals/data/data_source/p2p"
	sql_source "project/internals/data/data_source/sql"
	"project/internals/delivery/admin"
	auth_delivery "project/internals/delivery/authentication"
	delivery "project/internals/delivery/blockchain"
	"project/internals/delivery/sse"
	"project/internals/domain/entity"
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

	institutionchannel := make(chan entity.Institution)
	channelMap := make(map[string]chan<- entity.Institution)

	nodeSource := p2p.NewNodeSource(peerPorts)
	dbConn := sql_source.NewDB()
	sqlSource := sql_source.NewSQLSource(dbConn)
	sseService := service.NewSSEManager(channelMap)
	module := delivery.NewModule(source.NewBlockChainMemorySource(), nodeSource, sqlSource)
	authModule := auth_delivery.NewModule(sqlSource, institutionchannel, channelMap, sseService)
	sseUseCase := usecase.NewSSEUseCase(sqlSource, sseService)
	sseModule := sse.NewModule(sqlSource, sseService, sseUseCase)
	adminModule := admin.NewModule(sqlSource, service.Service{}, sseService)
	mux := http.NewServeMux()
	delivery.RegisterRoutes(mux, module)
	auth_delivery_routes := auth_delivery.RegisterRoutes(mux, authModule)
	sse_routes := sse.RegisterRoutes(mux, sseModule)
	admin_routes := admin.RegisterRoutes(mux, adminModule)

	//! structurize main.go as WELL
	var allNormalRoutes []common.RouteWrapper
	allNormalRoutes = append(allNormalRoutes, auth_delivery_routes...)
	allNormalRoutes = append(allNormalRoutes, admin_routes...)

	common.NewRouteWrapper(allNormalRoutes...)
	common.NewSSERouteWrapper(sse_routes)

	addr := fmt.Sprintf(":%d", *currentPort)
	fmt.Printf("🚀 Starting blockchain node on http://localhost%s\n", addr)

	memSource := source.NewBlockChainMemorySource()
	service := service.NewService()
	// blockchainService := service.NewService()
	blockChainUseCase := usecase.NewBlockChainUseCase(memSource, nodeSource, sqlSource, service)

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
