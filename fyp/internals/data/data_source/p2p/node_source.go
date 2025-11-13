package p2p

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	err "project/package/errors"
	logger "project/package/utils/pkg"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type NodeSource struct {
	// lastSequenceNumber      int
	peerPorts               []string
	digitalSignatureService service.DigitalSignatureService
	logger                  *zap.SugaredLogger
}

func NewNodeSource(peerPorts string, operationCounter *int, countCommitMap map[int]int, countPrepareMap map[int]int, pbftPeerPorts string) *NodeSource {
	return &NodeSource{
		digitalSignatureService: *service.NewDigitalSignature(),
		peerPorts:               strings.Split(peerPorts, ","),
		logger:                  logger.Logger,
		// countPrepareMap:         make(map[int]int),
		// countCommitMap:          make(map[int]int),
		// pbftPeerPorts: strings.Split(pbftPeerPorts, ","),
		// operationCounter:        operationCounter,
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

// func (n *NodeSource) SendPBFTMessageToPeer(pbftMessage entity.PBFTMessage, leaderNode int, currentMappedTCPPort int) (map[int]string, error) {
// 	ackMap := make(map[int]string)

// 	if er := checkPBFTMessageStructure(*n.operationCounter, pbftMessage); er != nil {
// 		n.logger.Errorln("[send_PBFT_Message_to_Peers]:error in PBFT structure Message ::", pbftMessage)
// 		return nil, er
// 	}
// 	n.logger.Infoln("[send_PBFT_Message_to_Peers]:case before SWITCH::", pbftMessage)
// 	switch pbftMessage.VerificationType {

// 	case enum.INITIAL: //pre-prepare actual process .( only for leader node)
// 		//n.receivedData = pbftMessage.QRVerificationRequestData
// 		messageSignature, er := n.digitalSignatureService.SignMessage(pbftMessage.QRVerificationRequestData)
// 		if er != nil {
// 			return nil, er
// 		}
// 		//! insert using HandleData or some function do these things .
// 		pbftMessage.NodeID = leaderNode
// 		pbftMessage.Signature = messageSignature
// 		pbftMessage.OperationID = *n.operationCounter
// 		//n.lastSequenceNumber = *n.operationCounter
// 		(*n.operationCounter)++
// 		pbftMessage.VerificationType = enum.PREPREPARE
// 		n.logger.Infoln("[send_PBFT_Message_to_Peers]:case INITIAL::", pbftMessage)
// 		_, hashByte, er := common.HashData(pbftMessage.QRVerificationRequestData)
// 		if er != nil {
// 			return nil, er
// 		}
// 		pbftMessage.QRVerificationRequestData.CertificateHash = hashByte
// 		if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
// 			return nil, er
// 		}
// 		//prepare PREPARE case as well :
// 		//immediately send PREPARE type as well(for leader node)
// 		n.logger.Infoln("[send_PBFT_Message_to_Peers]:case INITIAL but for PREPARE case::", pbftMessage)
// 		pbftMessage.VerificationType = enum.PREPARE
// 		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
// 		if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
// 			return nil, er
// 		}

// 	case enum.PREPREPARE: //prepare case

// 		if ok := isSequenceNumberCorrect(*n, pbftMessage); !ok {
// 			return nil, err.ErrNumberMismatch
// 		}

// 		pbftMessage.VerificationType = enum.PREPARE
// 		//n.receivedData = pbftMessage.QRVerificationRequestData
// 		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
// 		pbftMessage.NodeID = currentMappedTCPPort
// 		//pbftMessage.OperationID = *n.operationCounter
// 		n.logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPREPARE::", pbftMessage)
// 		if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
// 			return nil, er
// 		}

// 	case enum.PREPARE: //commit case
// 		// n.countPrepareMap[pbftMessage.OperationID]++
// 		if ok := isSequenceNumberCorrect(*n, pbftMessage); !ok {
// 			return nil, err.ErrNumberMismatch
// 		}
// 		countInt := n.countPrepareMap[pbftMessage.OperationID]
// 		n.logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE before count::", pbftMessage)
// 		n.logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE before count total till now==::", countInt)
// 		if countInt >= 2 {
// 			pbftMessage.VerificationType = enum.COMMIT
// 			pbftMessage.NodeID = currentMappedTCPPort
// 			n.logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE after count  count::", pbftMessage)
// 			if _, er := marshalAndBroadcast(pbftMessage, n, currentMappedTCPPort); er != nil {
// 				return nil, er
// 			}
// 		}
// 		return ackMap, nil
// 	case enum.COMMIT:
// 		if ok := isSequenceNumberCorrect(*n, pbftMessage); !ok {
// 			return nil, err.ErrNumberMismatch
// 		}
// 		countInt := n.countCommitMap[pbftMessage.OperationID]
// 		n.logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT before count::", countInt)
// 		if countInt >= 2 {
// 			n.logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT after count::", pbftMessage)

// 			return ackMap, nil
// 		}
// 	}

// 	return nil, nil
// }

// func (n *NodeSource) ReceivePBFTMessageToPeer(listenPort int, leaderPort int) (*entity.PBFTMessage, error) {
// 	n.logger.Debugln("value received", listenPort)
// 	// pbftMessageResponse := entity.PBFTMessage{}
// 	// var pbftMessageJSON []byte
// 	//waiting below as listenAndServe.
// 	ln, er := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort))
// 	if er != nil {
// 		n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
// 		return nil, er
// 	}
// 	defer ln.Close()
// 	for {
// 		conn, er := ln.Accept()
// 		if er != nil {
// 			n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
// 			return nil, er
// 		}

// 		go func(conn net.Conn) {
// 			reader := bufio.NewReader(conn)
// 			pbftMessageResponse := entity.PBFTMessage{}
// 			var pbftMessageJSON []byte
// 			defer func() {
// 				n.logger.Infow("CONNECTION CLOSED BY RECEIVER")
// 				conn.Close()

// 			}()

// 			length, er := reader.ReadString('\n')
// 			if er != nil {
// 				n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
// 				// return nil, er
// 				return
// 			}
// 			lengthInt, er := strconv.Atoi(strings.TrimSpace(length))
// 			if er != nil {
// 				n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", zap.Error(er))
// 				// return nil, err.ErrIntParse
// 				return
// 			}
// 			pbftMessageJSON = make([]byte, lengthInt)
// 			_, er = io.ReadFull(reader, pbftMessageJSON)
// 			if er != nil {
// 				n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
// 				// return nil, er
// 				return
// 			}
// 			er = json.Unmarshal(pbftMessageJSON, &pbftMessageResponse)
// 			if er != nil {
// 				n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
// 				// return nil, er
// 				return
// 			}
// 			n.logger.Debugln("[node_source] DEBUG: VALUE RECEIVED", pbftMessageResponse)
// 			_, er = conn.Write([]byte("ack\n"))
// 			if er != nil {
// 				n.logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
// 				// return nil, er
// 				return
// 			}

// 			switch pbftMessageResponse.VerificationType {
// 			case enum.PREPARE:
// 				n.countPrepareMap[pbftMessageResponse.OperationID]++

// 			case enum.COMMIT:
// 				n.countCommitMap[pbftMessageResponse.OperationID]++
// 				countInt := n.countPrepareMap[pbftMessageResponse.OperationID]
// 				if countInt >= 2 {
// 					//execution of the service.
// 					n.logger.Debugw("pbft message json for enum.COMMIT::", pbftMessageJSON)
// 				}
// 			}

// 			n.SendPBFTMessageToPeer(pbftMessageResponse, leaderPort, listenPort)
// 			//return &pbftMessageResponse, nil
// 		}(conn)
// 	}

// }

// // private helper method
