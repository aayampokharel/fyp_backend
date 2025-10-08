package main

import (
	"fmt"
	"net/http"

	"project/constants"
	"project/internals/data/config"
	source "project/internals/data/memory"
	"project/internals/data/p2p"
	delivery "project/internals/delivery/blockchain"
	"project/package/utils/common"
)

func main() {

	port := common.GetPort()

	env, err := config.NewEnv()
	if err != nil {
		fmt.Println("❌ Failed to load environment variables:", err)
		return
	}

	ports := env.GetValueForKey(constants.NodePortsKey)
	fmt.Println("Node Ports:", ports)

	module := delivery.NewModule(source.NewBlockChainMemorySource(), p2p.NewNodeSource(ports))

	mux := http.NewServeMux()
	delivery.RegisterRoutes(mux, module)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("🚀 Starting blockchain node on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("❌ Server failed:", err)
	}
}
