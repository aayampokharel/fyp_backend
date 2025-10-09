package common

import (
	"flag"
)

var tcpPort = flag.Int("port", 8000, "TCP port to use")

func GetPort() *int {
	flag.Parse()
	return tcpPort
}
