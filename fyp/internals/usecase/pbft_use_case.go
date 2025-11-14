package usecase

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"io"
	"net"
	"project/constants"
	"project/internals/data/config"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase/helper"
	"project/package/enum"
	err "project/package/errors"
	"project/package/utils/common"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type PBFTUseCase struct {
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

func NewPBFTUseCase(service service.Service, sqlRepo repository.ISqlRepository, nodeRepo repository.INodeRepository, countPrepareMap map[int]int, countCommitMap map[int]int, operationCounter *int, pbftService service.PBFTService, blockchainRepo repository.IBlockChainRepository) *PBFTUseCase {
	return &PBFTUseCase{
		SqlRepo:          sqlRepo,
		BlockChainRepo:   blockchainRepo,
		NodeRepo:         nodeRepo,
		Service:          service,
		countPrepareMap:  countPrepareMap,
		countCommitMap:   countCommitMap,
		operationCounter: operationCounter,
		PBFTService:      pbftService,
	}
}

func (uc *PBFTUseCase) SendPBFTMessageToPeerUseCase(pbftMessage entity.PBFTMessage) {

}

func (n *PBFTUseCase) SendPBFTMessageToPeer(pbftMessage entity.PBFTMessage, currentMappedTCPPort int) (map[int]string, error) {
	env, er := config.NewEnv()
	if er != nil {
		return nil, er
	}
	ackMap := make(map[int]string)
	leaderNodeString := env.GetValueForKey(constants.PbftLeaderNode)
	leaderNode, er := common.ConvertToInt(leaderNodeString)
	if er != nil {
		return nil, er
	}
	if er := n.PBFTService.CheckPBFTMessageStructure(*n.operationCounter, pbftMessage); er != nil {
		n.Service.Logger.Errorln("[send_PBFT_Message_to_Peers]:error in PBFT structure Message ::", pbftMessage)
		return nil, er
	}

	n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case before SWITCH::", pbftMessage)
	switch pbftMessage.VerificationType {

	case enum.INITIAL: //pre-prepare actual process .( only for leader node)
		//n.receivedData = pbftMessage.QRVerificationRequestData
		_, hashByte, er := common.HashData(pbftMessage.QRVerificationRequestData)
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while hashing data", er)
			return nil, er
		}
		pbftMessage.Digest = hex.EncodeToString(hashByte[:])

		messageSignature, er := n.DigitalSignatureService.SignMessageWithHash(hashByte[:])
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while signing message", er)
			return nil, er
		}
		//! insert using HandleData or some function do these things .
		pbftMessage.NodeID = leaderNode
		pbftMessage.Signature = messageSignature
		pbftMessage.OperationID = *n.operationCounter
		//n.lastSequenceNumber = *n.operationCounter
		(*n.operationCounter)++
		pbftMessage.VerificationType = enum.PREPREPARE
		n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case INITIAL::", pbftMessage)

		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while broadcasting", er)
			return nil, er
		}
		//prepare PREPARE case as well :
		//immediately send PREPARE type as well(for leader node)
		n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case INITIAL but for PREPARE case::", pbftMessage)
		pbftMessage.VerificationType = enum.PREPARE
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while broadcasting message", er)
			return nil, er
		}

	case enum.PREPREPARE: //prepare case

		if ok := n.PBFTService.IsSequenceNumberCorrect(n.operationCounter, pbftMessage); !ok {
			return nil, err.ErrNumberMismatch
		}
		if ok := helper.VerifyDigitalSignature(env, n.PBFTService.PbftPeerPorts, leaderNodeString, pbftMessage.Digest, pbftMessage.Signature); !ok {
			n.Service.Logger.Errorln("[send_PBFT_Message_to_Peers]:error in digital signature::", pbftMessage)
			return nil, err.ErrDigitalSignature

		}
		//!  i have to extract exact data from blockchian and set result .
		resultBlock, er := n.BlockChainRepo.ExtractBlockByHashAndCertificateID(string(pbftMessage.QRVerificationRequestData.CertificateHash), pbftMessage.QRVerificationRequestData.CertificateID)
		if er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while extracting block by hash and certificate id", er)
			return nil, er
		}
		previousBlock, _ := n.BlockChainRepo.GetBlockByBlockNumber(resultBlock.Header.BlockNumber - 1)
		if er = n.Service.VerifyBlock(resultBlock, previousBlock); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while verifying block", er)
			return nil, er
		}

		pbftMessage.Result = true
		//! other insertions do it in 1 helper method .
		pbftMessage.VerificationType = enum.PREPARE
		//n.receivedData = pbftMessage.QRVerificationRequestData
		pbftMessage.QRVerificationRequestData = entity.QRVerificationRequestData{}
		pbftMessage.NodeID = currentMappedTCPPort
		//pbftMessage.OperationID = *n.operationCounter
		n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPREPARE::", pbftMessage)
		if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while broadcasting", er)
			return nil, er
		}

	case enum.PREPARE: //commit case
		// n.countPrepareMap[pbftMessage.OperationID]++
		if ok := n.PBFTService.IsSequenceNumberCorrect(n.operationCounter, pbftMessage); !ok {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while checking sequence number", er)
			return nil, err.ErrNumberMismatch
		}
		countInt := n.countPrepareMap[pbftMessage.OperationID]
		n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE before count::", pbftMessage)
		n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE before count total till now==::", countInt)
		if countInt >= 2 {
			pbftMessage.VerificationType = enum.COMMIT
			pbftMessage.NodeID = currentMappedTCPPort
			n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case PREPARE after count  count::", pbftMessage)
			if _, er := n.PBFTService.MarshalAndBroadcast(pbftMessage, n.PBFTService.PbftPeerPorts, currentMappedTCPPort); er != nil {
				return nil, er
			}
		}
		return ackMap, nil
	case enum.COMMIT:
		if ok := n.PBFTService.IsSequenceNumberCorrect(n.operationCounter, pbftMessage); !ok {
			n.PBFTService.Logger.Errorln("[send_pbft_message_to_peer]:: Error while checking sequence number", er)
			return nil, err.ErrNumberMismatch
		}
		countInt := n.countCommitMap[pbftMessage.OperationID]
		n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT before count::", countInt)
		if countInt >= 2 {
			n.Service.Logger.Infoln("[send_PBFT_Message_to_Peers]:case COMMIT after count::", pbftMessage)

			return ackMap, nil
		}
	}

	return nil, nil
}

func (n *PBFTUseCase) ReceivePBFTMessageToPeer(listenPort int, leaderPort int) (*entity.PBFTMessage, error) {
	n.Service.Logger.Debugln("value received", listenPort)
	// pbftMessageResponse := entity.PBFTMessage{}
	// var pbftMessageJSON []byte
	//waiting below as listenAndServe.
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
				n.Service.Logger.Infow("CONNECTION CLOSED BY RECEIVER")
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
					n.Service.Logger.Debugw("pbft message json for enum.COMMIT::", pbftMessageJSON)
				}
			}

			n.SendPBFTMessageToPeer(pbftMessageResponse, leaderPort)
			//return &pbftMessageResponse, nil
		}(conn)
	}

}
