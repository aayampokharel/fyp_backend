package main

import (
	"fmt"
	"project/internals/data/source"
	delivery "project/internals/delivery/blockchain"
	"project/internals/domain/service"
	"project/internals/usecase"
)

func main() {
	// // init huma
	// api := huma.NewAPI("Blockchain API", "1.0.0")

	// // init repo
	// repo := postgres.NewBlockChainRepository()

	// // init module
	// module := delivery.NewModule(repo)

	// // register routes
	// delivery.RegisterRoutes(api, module)

	// // start server
	// huma.Listen(api)

	memorySource := source.NewBlockChainMemorySource()

	// 2️⃣ Create service (for MerkleRoot, POW, etc.)
	blockchainService := service.Service{}

	// 3️⃣ Create usecase (just like normal object)
	usecase := usecase.NewBlockChainUseCase(memorySource, blockchainService)

	// 4️⃣ Insert Genesis block
	controller := delivery.NewController(*usecase)
	controller.InsertNewCertificateData()
	fmt.Print("end whooho ")

}
