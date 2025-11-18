package usecase

import (
	"bufio"
	"encoding/hex"
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
	"sync"
	"time"
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
	mutex                   *sync.Mutex
}

func NewPBFTUseCase(service service.Service, sqlRepo repository.ISqlRepository, nodeRepo repository.INodeRepository, countPrepareMap map[int]int, countCommitMap map[int]int, operationCounter *int, pbftService service.PBFTService, blockchainRepo repository.IBlockChainRepository, operationChannels map[int]chan entity.PBFTExecutionResultEntity, mutex *sync.Mutex) *PBFTUseCase {
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
		mutex:             mutex,
	}
}

func (n *PBFTUseCase) SendPBFTMessageToPeer(pbftMessage entity.PBFTMessage, currentMappedTCPPort int, pbftChan chan entity.PBFTExecutionResultEntity) {
	env, er := config.NewEnv()
	if er != nil {
		n.Service.Logger.Errorw("[send_PBFT_Message_to_Peers] env load failed", "error", er)
		return
	}

	if er := n.PBFTService.CheckPBFTMessageStructure(*n.operationCounter, pbftMessage); er != nil {
		n.Service.Logger.Errorw("[send_PBFT_Message_to_Peers] invalid PBFT structure", "error", er, "message", pbftMessage)
		return

	}

	// n.Service.Logger.Infow("[send_PBFT_Message_to_Peers] processing message", "verification_type", pbftMessage.VerificationType, "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)

	switch pbftMessage.VerificationType {

	case enum.INITIAL: //pre-prepare actual process .( only for leader node)
		n.Service.Logger.Infow("[pbft][INITIAL] leader preparing pre-prepare", "current_port", currentMappedTCPPort, "operation_counter", *n.operationCounter)

		_, hashByte, er := common.HashData(pbftMessage.QRVerificationRequestData)
		if er != nil {
			n.PBFTService.Logger.Errorw("[pbft][INITIAL] error hashing request data", "error", er)

			return
		}
		// n.Service.Logger.Debugw("[pbft][INITIAL] hashed request data", "digest_len", len(hashByte))

		messageSignature, er := n.DigitalSignatureService.SignMessageWithHash(hashByte[:])
		if er != nil {
			n.PBFTService.Logger.Errorw("[pbft][INITIAL] signing failed", "error", er)

			return
		}
		pbftMessage.Digest = hashByte
		pbftMessage.NodeID = currentMappedTCPPort
		pbftMessage.Signature = messageSignature
		n.mutex.Lock()

		pbftMessage.OperationID = *n.operationCounter
		(*n.operationCounter)++
		n.operationChannels[pbftMessage.OperationID] = pbftChan

		n.mutex.Unlock()
		n.Service.Logger.Debugw("[pbft][INITIAL] broadcasting PREPREPARE & PREPARE", "operation_id", pbftMessage.OperationID)

		pbftMessage.VerificationType = enum.PREPREPARE
		//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case INITIAL::", pbftMessage)

		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorw("[pbft][INITIAL] broadcast PREPREPARE failed", "error", er)

			return
		}
		n.Service.Logger.Infow("[pbft][INITIAL] PREPREPARE broadcasted", "operation_id", pbftMessage.OperationID)

		pbftMessage.VerificationType = enum.PREPARE
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorw("[pbft][INITIAL] broadcast PREPARE failed", "error", er)

			return
		}

	case enum.PREPREPARE: //prepare case
		n.Service.Logger.Infow("[pbft][PREPREPARE] received PREPREPARE", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)

		if ok := helper.VerifyDigitalSignature(env, strconv.Itoa(pbftMessage.NodeID), pbftMessage.Digest, pbftMessage.Signature); !ok {
			n.Service.Logger.Errorw("[pbft][PREPREPARE] digital signature verification failed", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)
			n.Service.Logger.Infow("[pbft][PREPREPARE] dropping message due to signature failure", "operation_id", pbftMessage.OperationID)

			return

		}
		//!SET SIGNING OUT OF SWITCH INSTEAD OF ADDDING IN EACH CASE .

		resultBlock, er := n.BlockChainRepo.ExtractBlockByHashAndCertificateID(hex.EncodeToString(pbftMessage.QRVerificationRequestData.CertificateHash), pbftMessage.QRVerificationRequestData.CertificateID)
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while extracting block by hash and certificate id", er)
			return

		}
		previousBlock, _ := n.BlockChainRepo.GetBlockByBlockNumber(resultBlock.Header.BlockNumber - 1)
		if er = n.Service.VerifyBlock(resultBlock, previousBlock); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while verifying block", er)
			return
		}

		pbftMessage.Result = true

		_, hashByte, er := common.HashData(pbftMessage.QRVerificationRequestData)
		if er != nil {
			n.PBFTService.Logger.Errorw("[pbft][PREPREPARE] hashing for PREPARE failed", "error", er, "operation_id", pbftMessage.OperationID)

			return
		}
		pbftMessage.Digest = hashByte
		messageSignature, er := n.DigitalSignatureService.SignMessageWithHash(hashByte[:])
		if er != nil {
			n.PBFTService.Logger.Errorw("[pbft][PREPREPARE] signing PREPARE failed", "error", er, "operation_id", pbftMessage.OperationID)

			return
		}
		pbftMessage.Signature = messageSignature
		//! other insertions do it in 1 helper method .
		pbftMessage.VerificationType = enum.PREPARE
		//n.receivedData = pbftMessage.QRVerificationRequestData
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		pbftMessage.NodeID = currentMappedTCPPort
		n.Service.Logger.Infow("[pbft][PREPREPARE] broadcasting PREPARE from replica", "operation_id", pbftMessage.OperationID, "from_node", currentMappedTCPPort)

		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorw("[pbft][PREPREPARE] broadcast PREPARE failed", "error", er, "operation_id", pbftMessage.OperationID)

			return
		}
		n.Service.Logger.Debugw("[pbft][PREPREPARE] PREPARE broadcast sent", "operation_id", pbftMessage.OperationID)

	case enum.PREPARE: //commit case
		n.Service.Logger.Infow("[pbft][PREPARE] received PREPARE", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)

		if ok := helper.VerifyDigitalSignature(env, strconv.Itoa(pbftMessage.NodeID), pbftMessage.Digest, pbftMessage.Signature); !ok {
			n.Service.Logger.Errorw("[pbft][PREPARE] signature verification failed", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)
			n.Service.Logger.Infow("[pbft][PREPARE] ignoring PREPARE due to invalid signature", "operation_id", pbftMessage.OperationID)

			return

		}
		n.Service.Logger.Debugw("[pbft][PREPARE] signature valid for PREPARE", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)

		countInt := n.countPrepareMap[pbftMessage.OperationID]
		n.Service.Logger.Infow("[pbft][PREPARE] prepare count check", "operation_id", pbftMessage.OperationID, "prepare_count", countInt)

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
				n.PBFTService.Logger.Errorw("[pbft][PREPARE] signing COMMIT failed", "error", er, "operation_id", pbftMessage.OperationID)

				return
			}
			pbftMessage.Signature = messageSignature
			n.Service.Logger.Infow("[pbft][PREPARE] broadcasting COMMIT", "operation_id", pbftMessage.OperationID, "from_node", currentMappedTCPPort)

			if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
				n.PBFTService.Logger.Errorw("[pbft][PREPARE] broadcast COMMIT failed", "error", er, "operation_id", pbftMessage.OperationID)

				return
			}
		} else {
			n.Service.Logger.Debugw("[pbft][PREPARE] not enough PREPARE votes yet", "operation_id", pbftMessage.OperationID, "prepare_count", countInt)

		}

	case enum.COMMIT:
		n.Service.Logger.Infow("[pbft][COMMIT] received COMMIT", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)

		if ok := helper.VerifyDigitalSignature(env, strconv.Itoa(pbftMessage.NodeID), pbftMessage.Digest, pbftMessage.Signature); !ok {
			n.Service.Logger.Errorw("[pbft][COMMIT] signature verification failed for COMMIT", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)

			return

		}
		n.Service.Logger.Debugw("[pbft][COMMIT] signature valid", "node_id", pbftMessage.NodeID, "operation_id", pbftMessage.OperationID)

		countInt := n.countCommitMap[pbftMessage.OperationID]
		n.Service.Logger.Infow("[pbft][COMMIT] commit count check", "operation_id", pbftMessage.OperationID, "commit_count", countInt)

		if countInt >= 2 {
			//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT after count::", pbftMessage)
			//n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:SUCCESSFUL LOADING!!!!", pbftMessage)

			ch, exists := n.operationChannels[pbftMessage.OperationID]
			if exists && ch != nil {
				n.Service.Logger.Infow("[pbft][COMMIT] SUCCESS!!..sending execution result to waiting channel", "operation_id", pbftMessage.OperationID)

				go func() {
					select {
					case ch <- entity.NewPBFTExecutionResultEntity(true, nil):
						n.Service.Logger.Infow("[pbft][COMMIT] successfully delivered execution result", "operation_id", pbftMessage.OperationID)

					case <-time.After(3 * time.Second):
						n.Service.Logger.Infow("[pbft][COMMIT] timeout sending execution result to channel", "operation_id", pbftMessage.OperationID)

					}
				}()
			} else {
				// this is a replica node
				n.Service.Logger.Infow("[pbft][COMMIT] no waiting channel found; replica may ignore", "operation_id", pbftMessage.OperationID)
			}

		} else {
			n.Service.Logger.Debugw("[pbft][COMMIT] insufficient COMMIT votes to finalize", "operation_id", pbftMessage.OperationID, "commit_count", countInt)

		}
	}
}

func (n *PBFTUseCase) ReceivePBFTMessageToPeer(listenPort int) (*entity.PBFTMessage, error) {
	//n.Service.Logger.Debugln("value received", listenPort)
	ln, er := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort))
	if er != nil {
		n.Service.Logger.Errorw("[pbft][RECEIVER] TCP listen failed", "listen_port", listenPort, "error", er)

		return nil, er
	}
	defer ln.Close()
	for {
		conn, er := ln.Accept()
		if er != nil {
			n.Service.Logger.Errorw("[pbft][RECEIVER] accept failed", "listen_port", listenPort, "error", er)

			return nil, er
		}
		n.Service.Logger.Infow("[pbft][RECEIVER] connection accepted", "listen_port", listenPort, "remote_addr", conn.RemoteAddr().String())

		go func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			pbftMessageResponse := entity.PBFTMessage{}
			var pbftMessageJSON []byte
			defer func() {
				//n.Service.Logger.Infow("[pbft][RECEIVER] connection closing", "listen_port", listenPort)

				conn.Close()

			}()

			length, er := reader.ReadString('\n')
			if er != nil {
				n.Service.Logger.Errorw("[pbft][RECEIVER] read length failed", "error", er, "listen_port", listenPort)

				// return nil, er
				return
			}
			lengthInt, er := strconv.Atoi(strings.TrimSpace(length))
			if er != nil {
				n.Service.Logger.Errorw("[pbft][RECEIVER] parse length failed", "error", er, "length_str", length)

				// return nil, err.ErrIntParse
				return
			}
			//n.Service.Logger.Debugw("[pbft][RECEIVER] incoming message length", "length", lengthInt, "listen_port", listenPort)

			pbftMessageJSON = make([]byte, lengthInt)
			_, er = io.ReadFull(reader, pbftMessageJSON)
			if er != nil {
				n.Service.Logger.Errorw("[pbft][RECEIVER] read payload failed", "error", er, "expected_len", lengthInt)

				// return nil, er
				return
			}
			er = json.Unmarshal(pbftMessageJSON, &pbftMessageResponse)
			if er != nil {
				n.Service.Logger.Errorw("[pbft][RECEIVER] json unmarshal failed", "error", er, "payload", string(pbftMessageJSON))

				return
			}
			n.Service.Logger.Infow("[pbft][RECEIVER] message received and parsed", "verification_type", pbftMessageResponse.VerificationType, "operation_id", pbftMessageResponse.OperationID, "node_id", pbftMessageResponse.NodeID)

			_, er = conn.Write([]byte("ack\n"))
			if er != nil {
				n.Service.Logger.Errorw("[pbft][RECEIVER] ack write failed", "error", er, "listen_port", listenPort)

				return
			}
			//n.Service.Logger.Debugw("[pbft][RECEIVER] ack sent", "listen_port", listenPort)

			switch pbftMessageResponse.VerificationType {
			case enum.PREPARE:
				n.countPrepareMap[pbftMessageResponse.OperationID]++
				n.Service.Logger.Infow("[pbft][RECEIVER] PREPARE count incremented", "operation_id", pbftMessageResponse.OperationID, "new_prepare_count", n.countPrepareMap[pbftMessageResponse.OperationID])

			case enum.COMMIT:
				n.countCommitMap[pbftMessageResponse.OperationID]++
				n.Service.Logger.Infow("[pbft][RECEIVER] COMMIT count incremented", "operation_id", pbftMessageResponse.OperationID, "new_commit_count", n.countCommitMap[pbftMessageResponse.OperationID])

				countInt := n.countPrepareMap[pbftMessageResponse.OperationID]
				if countInt >= 2 {
					//execution of the service.
					n.Service.Logger.Debugw("[pbft][RECEIVER] execution ready on COMMIT", "operation_id", pbftMessageResponse.OperationID, "prepare_count", countInt, "commit_count", n.countCommitMap[pbftMessageResponse.OperationID])

				} else {
					n.Service.Logger.Debugw("[pbft][RECEIVER] commit received but prepare not sufficient", "operation_id", pbftMessageResponse.OperationID, "prepare_count", countInt)
				}
			}
			// _, ok := n.operationChannels[pbftMessageResponse.OperationID]

			n.SendPBFTMessageToPeer(pbftMessageResponse, listenPort, n.operationChannels[pbftMessageResponse.OperationID])

		}(conn)
	}

}
