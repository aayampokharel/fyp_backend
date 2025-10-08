package common

import (
	"flag"
)

var tcpPort = flag.Int("tcp.port", 0, "TCP port to use")

func GetPort() *int {
	flag.Parse() // parse the CLI flags
	return tcpPort
}
