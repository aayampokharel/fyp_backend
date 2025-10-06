package p2p

import (
	"encoding/json"
	"net"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	err "project/package/errors"
	"strconv"
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

func (n *NodeSource) SendBlockToPeer(block entity.Block, currentPort int) (map[int]string, error) {
	blockJson, er := json.Marshal(block)
	ackMap := make(map[int]string)
	if er != nil {
		return nil, err.ErrMarshaling
	}
	for _, portValString := range n.nodes {
		portValInt, er := strconv.Atoi(portValString)
		if er != nil {
			return nil, err.ErrIntParse
		}

		if currentPort != portValInt {

			conn, er := net.Dial("tcp", "localhost:"+portValString)
			if er != nil {
				ackMap[portValInt] = er.Error()
				continue
			}
			// conn.SetDeadline(time.Now().Add(3 * time.Second))//3 sec deadline may be imp if i dont want this node to be hanging if other side  doesnot respond
			_, er = conn.Write(blockJson)
			if er != nil {
				return nil, err.ErrMarshaling
			}
			ackbuf := make([]byte, 4096)
			count, er := conn.Read(ackbuf)
			if er != nil {
				return nil, err.ErrWithMoreInfo(err.ErrTcpListen, er.Error())
			}
			ackMap[portValInt] = string(ackbuf[:count])
			conn.Close()

		}
	}
	return ackMap, nil
}
