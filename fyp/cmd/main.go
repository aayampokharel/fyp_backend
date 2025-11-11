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
	"project/internals/delivery/category"
	filehandling "project/internals/delivery/file_handling"
	"project/internals/delivery/sse"
	"project/internals/domain/entity"
	"project/internals/domain/service"
	"project/internals/usecase"
	"project/package/utils/common"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

func main() {
	// -------------------------------
	// 1️⃣ Initialize Logger,maps & Config
	// -------------------------------
	logger.InitLogger()

	currentPort := common.GetPort()
	tcpPort := *currentPort + 1000
	pbftTcpPort := *currentPort + 1500

	env, err := config.NewEnv()
	if err != nil {
		logger.Logger.Errorw("[main] Failed to load environment variables", zap.Error(err))
		return
	}
	peerPorts := env.GetValueForKey(constants.TCPPortsKey)
	pbftPeerPorts := env.GetValueForKey(constants.PbftPortsKey)
	globalSequenceNumber := 0
	countPrepareMap := make(map[int]int, 0)
	countCommitMap := make(map[int]int, 0)
	fmt.Println("Peer Ports:", peerPorts)
	fmt.Println(" PBFT Peer Ports:", pbftPeerPorts)

	// -------------------------------
	// 2️⃣ Initialize Channels
	// -------------------------------
	institutionChannel := make(chan entity.Institution)
	channelMap := make(map[string]chan<- entity.Institution)

	// -------------------------------
	// 3️⃣ Initialize Data Sources
	// -------------------------------
	dbConn := sql_source.NewDB()
	sqlSource := sql_source.NewSQLSource(dbConn)
	nodeSource := p2p.NewNodeSource(peerPorts, &globalSequenceNumber, countCommitMap, countPrepareMap, pbftPeerPorts)
	memSource := source.NewBlockChainMemorySource()

	// -------------------------------
	// 4️⃣ Initialize Services
	// -------------------------------
	sseService := service.NewSSEManager(channelMap)
	svc := service.NewService()

	// -------------------------------
	// 5️⃣ Initialize Modules
	// -------------------------------
	blockchainModule := delivery.NewModule(memSource, nodeSource, sqlSource)
	authModule := auth_delivery.NewModule(sqlSource, institutionChannel, channelMap, sseService)
	sseUseCase := usecase.NewSSEUseCase(sqlSource, sseService)
	sseModule := sse.NewModule(sqlSource, sseService, sseUseCase)
	adminModule := admin.NewModule(sqlSource, *svc, sseService)
	fileHandlingModule := filehandling.NewModule(*svc, memSource, nodeSource, sqlSource)
	categoryModule := category.NewModule(sqlSource, *svc)

	// -------------------------------
	// 6️⃣ Setup HTTP Server & Routes
	// -------------------------------
	mux := http.NewServeMux()

	// Register Routes
	authDeliveryRoutes := auth_delivery.RegisterRoutes(mux, authModule)
	sseRoutes := sse.RegisterRoutes(mux, sseModule)
	adminRoutes := admin.RegisterRoutes(mux, adminModule)
	fileHandlingRoutes := filehandling.RegisterRoutes(mux, fileHandlingModule)
	categoryRoutes := category.RegisterRoutes(mux, categoryModule)
	blockChainRoutes := delivery.RegisterRoutes(mux, blockchainModule)

	// Wrap routes for internal usage
	common.NewRouteWrapper(authDeliveryRoutes...)
	common.NewRouteWrapper(categoryRoutes...)
	common.NewRouteWrapper(blockChainRoutes...)
	common.NewRouteWrapper(adminRoutes...)
	common.NewFileRouteWrapper(fileHandlingRoutes...)
	common.NewSSERouteWrapper(sseRoutes)

	addr := fmt.Sprintf(":%d", *currentPort)
	fmt.Printf("🚀 Starting blockchain node on http://localhost%s\n", addr)

	// -------------------------------
	// 7️⃣ Initialize Use Cases
	// -------------------------------
	blockChainUseCase := usecase.NewBlockChainUseCase(memSource, nodeSource, sqlSource, *svc)

	// -------------------------------
	// 8️⃣ Start background goroutines
	// -------------------------------
	go receiveBlocks(blockChainUseCase, tcpPort)
	go receivePbftMessage(blockChainUseCase, pbftTcpPort)

	// -------------------------------
	// 9️⃣ Start HTTP Server
	// -------------------------------
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Logger.Errorw("❌ Server failed", zap.Error(err))
	}
}

// -------------------------------
// Helper: Background Block Receiver
// -------------------------------
func receiveBlocks(uc *usecase.BlockChainUseCase, tcpPort int) {
	for {
		if err := uc.ReceiveBlockFromPeer(tcpPort); err != nil {
			logger.Logger.Errorw("[node_source] Error receiving block", zap.Error(err))
			fmt.Println("Error receiving block from peer:", err)
		}
	}
}
func receivePbftMessage(uc *usecase.BlockChainUseCase, pbftTcpPort int) {
	// uc.Service.Logger.Infoln("stated in port::", pbftTcpPort)

	if _, er := uc.ReceivePBFTMessageFromPeers(pbftTcpPort); er != nil {
		logger.Logger.Errorw("[node_source] Error receiving pbft message", zap.Error(er))
		fmt.Println("Error receiving block from peer:", er)
	}

}
