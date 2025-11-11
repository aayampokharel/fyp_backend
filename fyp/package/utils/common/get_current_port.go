package common

import (
	"flag"
	"project/constants"
	"project/internals/data/config"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

var httpPort = flag.Int("port", 8000, "Http port to use")

func GetPort() *int {
	flag.Parse()
	return httpPort
}
func GetMappedTCPPort() int {
	flag.Parse()
	return (*httpPort + 1000)
}

func GetMappedTCPPBFTPort() int {
	flag.Parse()
	return (*httpPort + 1500)
}
func GetLeaderPort() *int {
	env, err := config.NewEnv()
	if err != nil {
		logger.Logger.Errorw("[main] Failed to load environment variables", zap.Error(err))
		return nil
	}
	leaderNode := env.GetValueForKey(constants.PbftLeaderNode)
	leaderNodeInt, er := ConvertToInt(leaderNode)
	if er != nil {
		logger.Logger.Errorw("[main] Failed to load environment variables", zap.Error(er))
		return nil
	}
	return &leaderNodeInt
}
