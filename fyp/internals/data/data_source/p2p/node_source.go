package p2p

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/package/enum"
	err "project/package/errors"
	logger "project/package/utils/pkg"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type NodeSource struct {
	globalSequenceNumber *int
	countPrepareMap      map[int]int
	receivedData         entity.QRVerificationRequestData
	countCommitMap       map[int]int
	peerPorts            []string
	pbftPeerPorts        []string
	logger               *zap.SugaredLogger
}

func NewNodeSource(peerPorts string, globalSequenceNumber *int, countCommitMap map[int]int, countPrepareMap map[int]int, pbftPeerPorts string) *NodeSource {
	return &NodeSource{
		peerPorts:            strings.Split(peerPorts, ","),
		logger:               logger.Logger,
		pbftPeerPorts:        strings.Split(pbftPeerPorts, ","),
		globalSequenceNumber: globalSequenceNumber,
	}
}

var _ repository.INodeRepository = (*NodeSource)(nil)

func (n *NodeSource) SendBlockToPeer(block entity.Block, currentMappedTCPPort int) (map[int]string, error) {
	n.logger.Debugw("[node_source] Debug: PeerPorts::", "ports", n.peerPorts)
	blockJson, er := json.Marshal(block)
	ackMap := make(map[int]string)
	if er != nil {
		return nil, err.ErrMarshaling
	}
	blockJsonLen := int32(len(blockJson))
	for _, tcpPortValString := range n.peerPorts {
		tcpPortValInt, er := strconv.Atoi(tcpPortValString)
		if er != nil {
			n.logger.Errorw("[node_source] Error: SendBlockToPeer::", zap.Error(er))
			return nil, err.ErrIntParse
		}

		if currentMappedTCPPort != tcpPortValInt {
			conn, er := net.Dial("tcp", "localhost:"+tcpPortValString)
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", zap.Error(er))
				ackMap[tcpPortValInt] = er.Error()
				continue
			}
			// conn.SetDeadline(time.Now().Add(3 * time.Second))//3 sec deadline may be imp if i dont want this node to be hanging if other side  doesnot respond
			_, er = conn.Write([]byte(strconv.Itoa(int(blockJsonLen)) + "\n"))
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(nil, er.Error())
			}
			_, er = conn.Write(blockJson)
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(nil, er.Error())
			}

			// ackbuf := make([]byte, 4096)
			ackStrResponse, er := bufio.NewReader(conn).ReadString('\n')
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(err.ErrTcpListen, er.Error())
			}
			ackMap[tcpPortValInt] = strings.TrimSpace(ackStrResponse)
			conn.Close()

		}
	}
	return ackMap, nil
}

func (n *NodeSource) ReceiveBlockFromPeer(listenPort int) (*entity.Block, error) {
	blockResponse := entity.Block{}
	var blockByte []byte
	ln, er := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort))
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	defer ln.Close()
	conn, er := ln.Accept()
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	length, er := reader.ReadString('\n')
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	lengthInt, er := strconv.Atoi(strings.TrimSpace(length))
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", zap.Error(er))
		return nil, err.ErrIntParse
	}
	blockByte = make([]byte, lengthInt)
	_, er = io.ReadFull(reader, blockByte)
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	er = json.Unmarshal(blockByte, &blockResponse)
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	_, er = conn.Write([]byte("ack\n"))
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}

	return &blockResponse, nil

}

func (n *NodeSource) SendPBFTMessageToPeer(pbftMessage entity.PBFTMessage, leaderNode int, currentMappedTCPPort int) (map[int]string, error) {
	n.logger.Debugw("[node_source] Debug: PeerPorts::", "ports", n.peerPorts)
	ackMap := make(map[int]string)

	//pbftMessageJSONLen := int32(len(pbftMessageJSON))
	//check structure
	if er := checkPBFTMessageStructure(*n.globalSequenceNumber, pbftMessage); er != nil {
		return nil, er
	}
	switch pbftMessage.VerificationType {

	case enum.INITIAL: //pre-prepare actual process .( only for leader node)
		pbftMessage.NodeID = leaderNode
		n.receivedData = pbftMessage.QRVerificationRequestData
		pbftMessage.SequenceNumber = *n.globalSequenceNumber
		(*n.globalSequenceNumber)++
		pbftMessage.VerificationType = enum.PREPREPARE
		if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
			return nil, er
		}

		//prepare PREPARE case as well :
		//immediately send PREPARE type as well(for leader node)
		pbftMessage.VerificationType = enum.PREPARE
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
			return nil, er
		}

	case enum.PREPREPARE: //prepare case
		pbftMessage.VerificationType = enum.PREPARE
		n.receivedData = pbftMessage.QRVerificationRequestData
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		pbftMessage.NodeID = currentMappedTCPPort
		//pbftMessage.SequenceNumber = *n.globalSequenceNumber
		if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
			return nil, er
		}

	case enum.PREPARE: //commit case
		// n.countPrepareMap[pbftMessage.SequenceNumber]++

		countInt := n.countPrepareMap[pbftMessage.SequenceNumber]
		if countInt >= 2 {
			pbftMessage.VerificationType = enum.COMMIT
			if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
				return nil, er
			}
		}

		//attach signature and sequenceNumber
		// pbftMessageJSON, er := json.Marshal(pbftMessage)
		// if er != nil {
		// 	return nil, err.ErrMarshaling
		// }
		// //broadcast to other peers
		// ackMap, er = broadcastToOtherPeers(n, pbftMessage, currentMappedTCPPort, pbftMessageJSON)
		// if er != nil {
		// 	return nil, er
		// }
		return ackMap, nil
	}
	return nil, nil
}

func broadcastToOtherPeers(n *NodeSource, currentMappedTCPPort int, pbftMessageJSON []byte) (map[int]string, error) {
	ackMap := make(map[int]string)
	pbftMessageJSONLen := int32(len(pbftMessageJSON))
	for _, tcpPortValString := range n.pbftPeerPorts { //! MAKE THIS GENERAL METHOD TO BE USED IN BOTH SEND FUNCTION
		tcpPortValInt, er := strconv.Atoi(tcpPortValString)
		if er != nil {
			n.logger.Errorw("[node_source] Error: SendBlockToPeer::", zap.Error(er))
			return nil, err.ErrIntParse
		}

		if currentMappedTCPPort != tcpPortValInt {
			conn, er := net.Dial("tcp", "localhost:"+tcpPortValString)
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", zap.Error(er))
				ackMap[tcpPortValInt] = er.Error()
				continue
			}
			// conn.SetDeadline(time.Now().Add(3 * time.Second))//3 sec deadline may be imp if i dont want this node to be hanging if other side  doesnot respond
			_, er = conn.Write([]byte(strconv.Itoa(int(pbftMessageJSONLen)) + "\n"))
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(nil, er.Error())
			}
			_, er = conn.Write(pbftMessageJSON)
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(nil, er.Error())
			}

			// ackbuf := make([]byte, 4096)
			ackStrResponse, er := bufio.NewReader(conn).ReadString('\n')
			if er != nil {
				n.logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(err.ErrTcpListen, er.Error())
			}
			ackMap[tcpPortValInt] = strings.TrimSpace(ackStrResponse)
			conn.Close()

		}
	}
	return ackMap, nil
}

func (n *NodeSource) ReceivePBFTMessageToPeer(listenPort int, leaderPort int) (*entity.PBFTMessage, error) {
	pbftMessageResponse := entity.PBFTMessage{}
	var pbftMessageJSON []byte
	ln, er := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort))
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	defer ln.Close()
	conn, er := ln.Accept()
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	length, er := reader.ReadString('\n')
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	lengthInt, er := strconv.Atoi(strings.TrimSpace(length))
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", zap.Error(er))
		return nil, err.ErrIntParse
	}
	pbftMessageJSON = make([]byte, lengthInt)
	_, er = io.ReadFull(reader, pbftMessageJSON)
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	er = json.Unmarshal(pbftMessageJSON, &pbftMessageResponse)
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	_, er = conn.Write([]byte("ack\n"))
	if er != nil {
		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}

	switch pbftMessageResponse.VerificationType {
	case enum.PREPARE:
		n.countPrepareMap[pbftMessageResponse.SequenceNumber]++

	case enum.COMMIT:
		n.countCommitMap[pbftMessageResponse.SequenceNumber]++
		countInt := n.countPrepareMap[pbftMessageResponse.SequenceNumber]
		if countInt >= 2 {
			//execution of the service.
			n.logger.Debugw("pbft message json for enum.COMMIT::", pbftMessageJSON)
		}
	}

	n.SendPBFTMessageToPeer(pbftMessageResponse, leaderPort, listenPort)
	return &pbftMessageResponse, nil

}

// private helper method
func checkPBFTMessageStructure(globalSequenceNumber int, pbftMessage entity.PBFTMessage) error {
	if pbftMessage.SequenceNumber != globalSequenceNumber || pbftMessage.SequenceNumber < 0 {
		return err.ErrNumberMismatch
	}
	if pbftMessage.VerificationType != enum.PREPREPARE &&
		pbftMessage.VerificationType != enum.PREPARE &&
		pbftMessage.VerificationType != enum.COMMIT {
		return err.ErrInvalidType
	}

	if pbftMessage.QRVerificationRequestData.CertificateHash == "" {
		return err.ErrInvalidType
	}
	if pbftMessage.QRVerificationRequestData.ClientID == "" {
		return err.ErrInvalidType
	}
	if pbftMessage.QRVerificationRequestData.Timestamp <= 0 {
		return err.ErrInvalidType
	}
	return nil
}

func marshalAndBroadcast(pbftMessage entity.PBFTMessage, n *NodeSource, currentMappedTCPPort int) (map[int]string, error) {
	pbftMessageJSON, er := json.Marshal(pbftMessage)
	if er != nil {
		return nil, err.ErrMarshaling
	}
	returnMap, er := broadcastToOtherPeers(n, currentMappedTCPPort, pbftMessageJSON)
	n.logger.Infof("[node_source] Debug: return map::", "return map", returnMap)
	if er != nil {
		return nil, er
	}
	return returnMap, nil
}
