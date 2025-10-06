package p2p

import (
	"encoding/json"
	"net"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	err "project/package/errors"
	"strconv"
	"time"
)

type NodeSource struct {
	nodes []string
}

func NewNodeSource(nodes []string) *NodeSource {
	return &NodeSource{
		nodes: nodes,
	}
}

var _ repository.INodeRepository = (*NodeSource)(nil)

func (n *NodeSource) SendBlockToPeer(block entity.Block, currentPort int) error {
	blockJson, er := json.Marshal(block)
	if er != nil {
		return err.ErrMarshaling
	}
	for _, portValString := range n.nodes {
		portValInt, er := strconv.Atoi(portValString)
		if er != nil {
			return err.ErrIntParse
		}

		if currentPort != portValInt {

			conn, er := net.Dial("tcp", "localhost:"+portValString)
			if er != nil {
				return err.ErrTcpListen
			}
			conn.SetDeadline(time.Now().Add(3 * time.Second))//3 sec deadline
			_, er = conn.Write(blockJson)
			conn.Close()
			if er != nil {
				return err.ErrMarshaling
			}

		}
	}
	return nil
}
