package usecase

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"project/internals/data/config"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase/helper"
	"project/package/enum"
	"project/package/utils/common"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type PBFTUseCase struct {
	operationChannels       map[int]chan entity.PBFTExecutionResultEntity
	SqlRepo                 repository.ISqlRepository
	NodeRepo                repository.INodeRepository
	BlockChainRepo          repository.IBlockChainRepository
	operationCounter        *int
	Service                 service.Service
	DigitalSignatureService service.DigitalSignatureService
	PBFTService             service.PBFTService
	countPrepareMap         map[int]int
	countCommitMap          map[int]int
}

func NewPBFTUseCase(service service.Service, sqlRepo repository.ISqlRepository, nodeRepo repository.INodeRepository, countPrepareMap map[int]int, countCommitMap map[int]int, operationCounter *int, pbftService service.PBFTService, blockchainRepo repository.IBlockChainRepository, operationChannels map[int]chan entity.PBFTExecutionResultEntity) *PBFTUseCase {
	return &PBFTUseCase{
		SqlRepo:           sqlRepo,
		BlockChainRepo:    blockchainRepo,
		NodeRepo:          nodeRepo,
		Service:           service,
		countPrepareMap:   countPrepareMap,
		countCommitMap:    countCommitMap,
		operationCounter:  operationCounter,
		operationChannels: operationChannels,
		PBFTService:       pbftService,
	}
}

// func (n *PBFTUseCase) RegisterOperationChannel(operationID int, ch chan entity.PBFTExecutionResultEntity) {
// 	n.channelMutex.Lock()
// 	defer n.channelMutex.Unlock()

// 	if n.operationChannels == nil {
// 		n.operationChannels = make(map[int]chan entity.PBFTExecutionResultEntity)
// 	}
// 	n.operationChannels[operationID] = ch
// }

// func (n *PBFTUseCase) GetOperationChannel(operationID int) (chan entity.PBFTExecutionResultEntity, bool) {
// 	n.channelMutex.RLock()
// 	defer n.channelMutex.RUnlock()

// 	ch, exists := n.operationChannels[operationID]
// 	return ch, exists
// }

func (uc *PBFTUseCase) SendPBFTMessageToPeerUseCase(pbftMessage entity.PBFTMessage) {

}

func (n *PBFTUseCase) SendPBFTMessageToPeer(pbftMessage entity.PBFTMessage, currentMappedTCPPort int, pbftChan chan entity.PBFTExecutionResultEntity) {
	env, er := config.NewEnv()
	if er != nil {
		return
	}

	if er := n.PBFTService.CheckPBFTMessageStructure(*n.operationCounter, pbftMessage); er != nil {
		n.Service.Logger.Errorln("[send_PBFT_Message_to_Peers]:error in PBFT structure Message ::", pbftMessage)
		return
	}

	//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case before SWITCH::", pbftMessage)

	switch pbftMessage.VerificationType {

	case enum.INITIAL: //pre-prepare actual process .( only for leader node)
		//n.receivedData = pbftMessage.QRVerificationRequestData

		_, hashByte, er := common.HashData(pbftMessage.QRVerificationRequestData)
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while hashing data", er)
			return
		}
		messageSignature, er := n.DigitalSignatureService.SignMessageWithHash(hashByte[:])
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while signing message", er)
			return
		}
		pbftMessage.Digest = hashByte
		pbftMessage.NodeID = currentMappedTCPPort
		pbftMessage.Signature = messageSignature
		pbftMessage.OperationID = *n.operationCounter
		(*n.operationCounter)++
		n.operationChannels[pbftMessage.OperationID] = pbftChan

		pbftMessage.VerificationType = enum.PREPREPARE
		//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case INITIAL::", pbftMessage)

		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while broadcasting", er)
			return
		}
		//prepare PREPARE case as well :
		//immediately send PREPARE type as well(for leader node)
		//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case INITIAL but for PREPARE case::", pbftMessage)
		pbftMessage.VerificationType = enum.PREPARE
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while broadcasting message", er)
			return
		}

	case enum.PREPREPARE: //prepare case
		// if pbftMessage.OperationID >= *n.operationCounter {
		// 	*n.operationCounter = pbftMessage.OperationID
		// }

		// if ok := n.PBFTService.IsSequenceNumberCorrect(n.operationCounter, pbftMessage); !ok {
		// 	return
		// }

		if ok := helper.VerifyDigitalSignature(env, strconv.Itoa(pbftMessage.NodeID), pbftMessage.Digest, pbftMessage.Signature); !ok {
			n.Service.Logger.Errorln("[send_PBFT_Message_to_Peers]:error in digital signature::", pbftMessage)
			n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:node nu", pbftMessage.NodeID)
			return

		}
		//!SET SIGNING OUT OF SWITCH INSTEAD OF ADDDING IN EACH CASE .
		// resultBlock, er := n.BlockChainRepo.ExtractBlockByHashAndCertificateID(string(pbftMessage.QRVerificationRequestData.CertificateHash), pbftMessage.QRVerificationRequestData.CertificateID)
		// if er != nil {
		// 	n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while extracting block by hash and certificate id", er)
		// 	return

		// }
		// previousBlock, _ := n.BlockChainRepo.GetBlockByBlockNumber(resultBlock.Header.BlockNumber - 1)
		// if er = n.Service.VerifyBlock(resultBlock, previousBlock); er != nil {
		// 	n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while verifying block", er)
		// 	return
		// }

		pbftMessage.Result = true

		_, hashByte, er := common.HashData(pbftMessage.QRVerificationRequestData)
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while hashing data", er)
			return
		}
		pbftMessage.Digest = hashByte
		messageSignature, er := n.DigitalSignatureService.SignMessageWithHash(hashByte[:])
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while signing message", er)
			return
		}
		pbftMessage.Signature = messageSignature
		//! other insertions do it in 1 helper method .
		pbftMessage.VerificationType = enum.PREPARE
		//n.receivedData = pbftMessage.QRVerificationRequestData
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		pbftMessage.NodeID = currentMappedTCPPort
		//pbftMessage.OperationID = *n.operationCounter
		//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPREPARE::", pbftMessage)
		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while broadcasting", er)
			return
		}

	case enum.PREPARE: //commit case
		// n.countPrepareMap[pbftMessage.OperationID]++
		// if ok := n.PBFTService.IsSequenceNumberCorrect(n.operationCounter, pbftMessage); !ok {
		// 	n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while checking sequence number", er)
		// 	return
		// }
		if ok := helper.VerifyDigitalSignature(env, strconv.Itoa(pbftMessage.NodeID), pbftMessage.Digest, pbftMessage.Signature); !ok {
			n.Service.Logger.Errorln("[send_PBFT_Message_to_Peers]:error in digital signature::", pbftMessage)
			n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:NODE NO ", pbftMessage.NodeID)
			return

		}
		countInt := n.countPrepareMap[pbftMessage.OperationID]
		//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE before count::", pbftMessage)
		//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE before count total till now==::", countInt)

		if countInt >= 2 {
			pbftMessage.VerificationType = enum.COMMIT
			pbftMessage.NodeID = currentMappedTCPPort
			// _, hashByte, er := common.HashData(pbftMessage.QRVerificationRequestData)
			// if er != nil {
			// 	n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while hashing data", er)
			// 	return
			// }
			// pbftMessage.Digest = hashByte
			messageSignature, er := n.DigitalSignatureService.SignMessageWithHash(pbftMessage.Digest[:])
			if er != nil {
				n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while signing message", er)
				return
			}
			pbftMessage.Signature = messageSignature
			//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE after count  count::", pbftMessage)
			if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
				return
			}
		}

	case enum.COMMIT:
		// if ok := n.PBFTService.IsSequenceNumberCorrect(n.operationCounter, pbftMessage); !ok {
		// 	n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while checking sequence number", er)
		// 	return
		// }
		if ok := helper.VerifyDigitalSignature(env, strconv.Itoa(pbftMessage.NodeID), pbftMessage.Digest, pbftMessage.Signature); !ok {
			n.Service.Logger.Errorln("[send_PBFT_Message_to_Peers]:error in digital signature::", pbftMessage)
			return

		}

		countInt := n.countCommitMap[pbftMessage.OperationID]
		//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT before count::", countInt)
		if countInt >= 2 {
			//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT after count::", pbftMessage)
			//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:SUCCESSFUL LOADING!!!!", pbftMessage)

			ch, exists := n.operationChannels[pbftMessage.OperationID]
			if exists && ch != nil {
				n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT after count::", "SUCCESSSS!!!!!")
				go func() {
					select {
					case ch <- entity.NewPBFTExecutionResultEntity(true, nil):
						n.Service.Logger.Infoln("Successfully sent to channel")
					case <-time.After(1 * time.Second):
						n.Service.Logger.Warn("Timeout sending to channel")
					}
				}()
			} else {
				// this is a replica node
				n.Service.Logger.Infoln("Replica node COMMIT reached, ignoring channel send")
			}

		}
	}
}

func (n *PBFTUseCase) ReceivePBFTMessageToPeer(listenPort int) (*entity.PBFTMessage, error) {
	//n.Service.Logger.Debugln("value received", listenPort)
	ln, er := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort))
	if er != nil {
		n.Service.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
		return nil, er
	}
	defer ln.Close()
	for {
		conn, er := ln.Accept()
		if er != nil {
			n.Service.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
			return nil, er
		}

		go func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			pbftMessageResponse := entity.PBFTMessage{}
			var pbftMessageJSON []byte
			defer func() {
				n.Service.Logger.Infow("CONNECTION CLOSED BY RECEIVER", listenPort)
				conn.Close()

			}()

			length, er := reader.ReadString('\n')
			if er != nil {
				n.Service.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
				// return nil, er
				return
			}
			lengthInt, er := strconv.Atoi(strings.TrimSpace(length))
			if er != nil {
				n.Service.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", zap.Error(er))
				// return nil, err.ErrIntParse
				return
			}
			pbftMessageJSON = make([]byte, lengthInt)
			_, er = io.ReadFull(reader, pbftMessageJSON)
			if er != nil {
				n.Service.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
				// return nil, er
				return
			}
			er = json.Unmarshal(pbftMessageJSON, &pbftMessageResponse)
			if er != nil {
				n.Service.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
				// return nil, er
				return
			}
			n.Service.Logger.Debugln("[node_source] DEBUG: VALUE RECEIVED", pbftMessageResponse)
			_, er = conn.Write([]byte("ack\n"))
			if er != nil {
				n.Service.Logger.Errorw("[node_source] Error: ReceiveBlockFromPeer::", er)
				// return nil, er
				return
			}

			switch pbftMessageResponse.VerificationType {
			case enum.PREPARE:
				n.countPrepareMap[pbftMessageResponse.OperationID]++

			case enum.COMMIT:
				n.countCommitMap[pbftMessageResponse.OperationID]++
				countInt := n.countPrepareMap[pbftMessageResponse.OperationID]
				if countInt >= 2 {
					//execution of the service.
					n.Service.Logger.Debugw("pbft message json for enum.COMMIT::", "operation_id", pbftMessageResponse.OperationID,
						"node_id", pbftMessageResponse.NodeID,
						"verification_type", pbftMessageResponse.VerificationType)
				}
			}
			// _, ok := n.operationChannels[pbftMessageResponse.OperationID]

			n.SendPBFTMessageToPeer(pbftMessageResponse, listenPort, n.operationChannels[pbftMessageResponse.OperationID])

		}(conn)
	}

}
