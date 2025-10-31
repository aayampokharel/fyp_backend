// package main

// import (
// 	"fmt"
// 	"net/http"

// 	"project/constants"
// 	"project/internals/data/config"
// 	source "project/internals/data/data_source/memory"
// 	"project/internals/data/data_source/p2p"
// 	sql_source "project/internals/data/data_source/sql"
// 	"project/internals/delivery/admin"
// 	auth_delivery "project/internals/delivery/authentication"
// 	delivery "project/internals/delivery/blockchain"
// 	filehandling "project/internals/delivery/file_handling"
// 	"project/internals/delivery/sse"
// 	"project/internals/domain/entity"
// 	"project/internals/domain/service"
// 	"project/internals/usecase"
// 	"project/package/utils/common"
// 	logger "project/package/utils/pkg"

// 	"go.uber.org/zap"
// )

// func main() {
// 	logger.InitLogger()
// 	currentPort := common.GetPort()
// 	tcpPort := *currentPort + 1000 // e.g., 8001 -> 9001

// 	// logger.Logger.Infoln("[main] Info: Current Port::", *currentPort)
// 	env, err := config.NewEnv()
// 	if err != nil {
// 		logger.Logger.Errorw("[main] Error: Failed to load environment variables", zap.Error(err))
// 		return
// 	}
// 	peerPorts := env.GetValueForKey(constants.TCPPortsKey)
// 	fmt.Println("Peer Ports:", peerPorts)

// 	institutionchannel := make(chan entity.Institution)
// 	channelMap := make(map[string]chan<- entity.Institution)

// 	nodeSource := p2p.NewNodeSource(peerPorts)
// 	dbConn := sql_source.NewDB()
// 	sqlSource := sql_source.NewSQLSource(dbConn)
// 	sseService := service.NewSSEManager(channelMap)
// 	blockchainsource := source.NewBlockChainMemorySource()
// 	module := delivery.NewModule(blockchainsource, nodeSource, sqlSource)
// 	authModule := auth_delivery.NewModule(sqlSource, institutionchannel, channelMap, sseService)
// 	sseUseCase := usecase.NewSSEUseCase(sqlSource, sseService)
// 	sseModule := sse.NewModule(sqlSource, sseService, sseUseCase)
// 	adminModule := admin.NewModule(sqlSource, *service.NewService(), sseService)
// 	mux := http.NewServeMux()
// 	delivery.RegisterRoutes(mux, module)
// 	auth_delivery_routes := auth_delivery.RegisterRoutes(mux, authModule)
// 	sse_routes := sse.RegisterRoutes(mux, sseModule)
// 	admin_routes := admin.RegisterRoutes(mux, adminModule)
// 	fileHandlingModule := filehandling.NewModule(*service.NewService(), blockchainsource, nodeSource, sqlSource)
// 	fileHandling_routes := filehandling.RegisterRoutes(mux, fileHandlingModule)

// 	//! structurize main.go as WELL
// 	var allNormalRoutes []common.RouteWrapper
// 	allNormalRoutes = append(allNormalRoutes, auth_delivery_routes...)
// 	allNormalRoutes = append(allNormalRoutes, admin_routes...)

// 	var allFileHandlingRoutes []common.FileRouteWrapper
// 	allFileHandlingRoutes = append(allFileHandlingRoutes, fileHandling_routes...)
// 	common.NewRouteWrapper(allNormalRoutes...)
// 	common.NewFileRouteWrapper(allFileHandlingRoutes...)
// 	common.NewSSERouteWrapper(sse_routes)

// 	addr := fmt.Sprintf(":%d", *currentPort)
// 	fmt.Printf("🚀 Starting blockchain node on http://localhost%s\n", addr)

// 	memSource := source.NewBlockChainMemorySource()
// 	service := service.NewService()
// 	// blockchainService := service.NewService()
// 	blockChainUseCase := usecase.NewBlockChainUseCase(memSource, nodeSource, sqlSource, *service)

// 	go func() {
// 		for {
// 			er := blockChainUseCase.ReceiveBlockFromPeer(tcpPort)
// 			if er != nil {
// 				logger.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", zap.Error(er))
// 				fmt.Println("Error receiving block from peer:", er)
// 			}

// 		}
// 	}()
// 	if err := http.ListenAndServe(addr, mux); err != nil {
// 		fmt.Println("❌ Server failed:", err)
// 	}

// }

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
	filehandling "project/internals/delivery/file_handling"
	"project/internals/delivery/sse"
	"project/internals/domain/entity"
	"project/internals/domain/service"
	"project/internals/usecase"
	"project/package/utils/common"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

type Application struct { //wrappers of wrappers
	env         *config.Env
	currentPort int
	tcpPort     int
	mux         *http.ServeMux
	modules     *Modules
	dataSources *DataSources
	services    *Services
	useCases    *UseCases
	channels    *Channels
}

type DataSources struct {
	nodeSource       *p2p.NodeSource
	sqlSource        *sql_source.SQLSource
	blockchainSource *source.BlockChainMemorySource
}

type Services struct {
	sseService  *service.SSEManager
	baseService *service.Service
}

type UseCases struct {
	sseUseCase        *usecase.SSEUseCase
	blockchainUseCase *usecase.BlockChainUseCase
}

type Modules struct {
	blockchainModule   *delivery.Module
	authModule         *auth_delivery.Module
	sseModule          *sse.Module
	adminModule        *admin.Module
	fileHandlingModule *filehandling.Module
}

type Channels struct {
	institutionChannel chan entity.Institution
	channelMap         map[string]chan<- entity.Institution
}

func NewApplication() (*Application, error) {
	app := &Application{}

	if err := app.initialize(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) initialize() error {
	logger.InitLogger()
	if err := app.setupPorts(); err != nil {
		return err
	}

	if err := app.setupEnvironment(); err != nil {
		return err
	}

	app.setupDataSources()
	app.setupChannels()
	app.setupServices()
	app.setupUseCases()
	app.setupModules()
	app.setupHTTPServer()
	app.registerRoutes()

	return nil
}

func (app *Application) setupPorts() error {
	currentPort := common.GetPort()
	if currentPort == nil {
		return fmt.Errorf("failed to get current port")
	}

	app.currentPort = *currentPort
	app.tcpPort = app.currentPort + 1000 // e.g., 8001 -> 9001

	return nil
}

func (app *Application) setupEnvironment() error {
	env, err := config.NewEnv()
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	app.env = env
	peerPorts := app.env.GetValueForKey(constants.TCPPortsKey)
	fmt.Println("Peer Ports:", peerPorts)

	return nil
}

func (app *Application) setupDataSources() {
	peerPorts := app.env.GetValueForKey(constants.TCPPortsKey)

	app.dataSources = &DataSources{
		nodeSource:       p2p.NewNodeSource(peerPorts),
		sqlSource:        sql_source.NewSQLSource(sql_source.NewDB()),
		blockchainSource: source.NewBlockChainMemorySource(),
	}
}

func (app *Application) setupChannels() {
	app.channels = &Channels{
		institutionChannel: make(chan entity.Institution),
		channelMap:         make(map[string]chan<- entity.Institution),
	}
}
func (app *Application) setupServices() {
	app.services = &Services{
		sseService:  service.NewSSEManager(app.channels.channelMap),
		baseService: service.NewService(),
	}
}

func (app *Application) setupUseCases() {
	app.useCases = &UseCases{
		sseUseCase: usecase.NewSSEUseCase(app.dataSources.sqlSource, app.services.sseService),
		blockchainUseCase: usecase.NewBlockChainUseCase(
			app.dataSources.blockchainSource,
			app.dataSources.nodeSource,
			app.dataSources.sqlSource,
			*app.services.baseService,
		),
	}
}

func (app *Application) setupModules() {
	app.modules = &Modules{
		blockchainModule: delivery.NewModule(
			app.dataSources.blockchainSource,
			app.dataSources.nodeSource,
			app.dataSources.sqlSource,
		),
		authModule: auth_delivery.NewModule(
			app.dataSources.sqlSource,
			app.channels.institutionChannel,
			app.channels.channelMap,
			app.services.sseService,
		),
		sseModule: sse.NewModule(
			app.dataSources.sqlSource,
			app.services.sseService,
			app.useCases.sseUseCase,
		),
		adminModule: admin.NewModule(
			app.dataSources.sqlSource,
			*app.services.baseService,
			app.services.sseService,
		),
		fileHandlingModule: filehandling.NewModule(
			*app.services.baseService,
			app.dataSources.blockchainSource,
			app.dataSources.nodeSource,
			app.dataSources.sqlSource,
		),
	}
}

func (app *Application) setupHTTPServer() {
	app.mux = http.NewServeMux()
}

func (app *Application) registerRoutes() {
	delivery.RegisterRoutes(app.mux, app.modules.blockchainModule)
	authRoutes := auth_delivery.RegisterRoutes(app.mux, app.modules.authModule)
	blockChainRoutes := delivery.RegisterRoutes(app.mux, app.modules.blockchainModule)
	sseRoutes := sse.RegisterRoutes(app.mux, app.modules.sseModule)
	adminRoutes := admin.RegisterRoutes(app.mux, app.modules.adminModule)
	fileHandlingRoutes := filehandling.RegisterRoutes(app.mux, app.modules.fileHandlingModule)

	app.wrapAndRegisterRoutes(authRoutes, blockChainRoutes, adminRoutes, fileHandlingRoutes, sseRoutes)
}

func (app *Application) wrapAndRegisterRoutes(
	authRoutes []common.RouteWrapper,
	blockchainRoutes []common.RouteWrapper,
	adminRoutes []common.RouteWrapper,
	fileHandlingRoutes []common.FileRouteWrapper,
	sseRoutes common.SSERouteWrapper,
) {
	var allNormalRoutes []common.RouteWrapper
	allNormalRoutes = append(allNormalRoutes, authRoutes...)
	allNormalRoutes = append(allNormalRoutes, adminRoutes...)
	allNormalRoutes = append(allNormalRoutes, blockchainRoutes...)

	common.NewRouteWrapper(allNormalRoutes...)
	common.NewFileRouteWrapper(fileHandlingRoutes...)
	common.NewSSERouteWrapper(sseRoutes)
}

func (app *Application) startPeerReceiver() {
	go func() {
		for {
			err := app.useCases.blockchainUseCase.ReceiveBlockFromPeer(app.tcpPort)
			if err != nil {
				logger.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer", zap.Error(err))
				fmt.Println("Error receiving block from peer:", err)
			}
		}
	}()
}

func (app *Application) Run() error {
	app.startPeerReceiver()

	addr := fmt.Sprintf(":%d", app.currentPort)
	fmt.Printf("🚀 Starting blockchain node on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, app.mux); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func main() {
	app, err := NewApplication()
	if err != nil {
		logger.Logger.Errorw("[main] Error: Failed to initialize application", zap.Error(err))
		return
	}

	if err := app.Run(); err != nil {
		logger.Logger.Errorw("[main] Error: Application failed to run", zap.Error(err))
		fmt.Println("❌ Application failed:", err)
	}
}
