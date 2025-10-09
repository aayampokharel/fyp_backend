package common

import (
	"flag"
)

var httpPort = flag.Int("port", 8001, "Http port to use")

func GetPort() *int {
	flag.Parse()
	return httpPort
}
func GetMappedTCPPort() int {
	flag.Parse()
	return (*httpPort + 1000)
}
