package p2p

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	err "project/package/errors"
	"strconv"
	"strings"
)

type NodeSource struct {
	nodes []string
}

func NewNodeSource(nodes string) *NodeSource {
	return &NodeSource{
		nodes: strings.Split(nodes, ","),
	}
}

var _ repository.INodeRepository = (*NodeSource)(nil)

func (n *NodeSource) SendBlockToPeer(block entity.Block, currentPort int) (map[int]string, error) {
	blockJson, er := json.Marshal(block)
	ackMap := make(map[int]string)
	if er != nil {
		return nil, err.ErrMarshaling
	}
	blockJsonLen := int32(len(blockJson))
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
			_, er = conn.Write([]byte(strconv.Itoa(int(blockJsonLen)) + "\n"))
			if er != nil {
				return nil, err.ErrWithMoreInfo(nil, er.Error())
			}
			_, er = conn.Write(blockJson)
			if er != nil {
				return nil, err.ErrWithMoreInfo(nil, er.Error())
			}

			// ackbuf := make([]byte, 4096)
			ackStrResponse, er := bufio.NewReader(conn).ReadString('\n')
			if er != nil {
				return nil, err.ErrWithMoreInfo(err.ErrTcpListen, er.Error())
			}
			ackMap[portValInt] = strings.TrimSpace(ackStrResponse)
			conn.Close()

		}
	}
	return ackMap, nil
}

func (n *NodeSource) ReceiveBlockFromPeer(listenPort int) (*entity.Block, error) {
	blockResponse := entity.Block{}
	var blockByte []byte
	ln, er := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort))
	defer ln.Close()
	if er != nil {
		return nil, er
	}
	conn, er := ln.Accept()
	defer conn.Close()
	if er != nil {
		return nil, er
	}
	reader := bufio.NewReader(conn)
	length, er := reader.ReadString('\n')
	if er != nil {
		return nil, er
	}
	lengthInt, er := strconv.Atoi(strings.TrimSpace(length))
	if er != nil {
		return nil, err.ErrIntParse
	}
	blockByte = make([]byte, lengthInt)
	_, er = io.ReadFull(reader, blockByte)
	if er != nil {
		return nil, er
	}
	er = json.Unmarshal(blockByte, &blockResponse)
	if er != nil {
		return nil, er
	}
	_, er = conn.Write([]byte("ack\n"))
	if er != nil {
		return nil, er
	}

	return &blockResponse, nil

}
