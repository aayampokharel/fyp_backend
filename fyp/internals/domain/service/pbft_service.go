package service

import (
	"bufio"
	"encoding/json"
	"net"
	"project/internals/domain/entity"
	"project/package/enum"
	err "project/package/errors"
	logger "project/package/utils/pkg"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type PBFTService struct {
	Logger        *zap.SugaredLogger
	PbftPeerPorts []string
}

func NewPBFTService(pbftPeerPortsString string) *PBFTService {
	return &PBFTService{Logger: logger.Logger, PbftPeerPorts: strings.Split(pbftPeerPortsString, ",")}
}

func (p *PBFTService) BroadcastToOtherPeers(pbftPeerPorts []string, currentMappedTCPPort int, pbftMessageJSON []byte) (map[int]string, error) {
	ackMap := make(map[int]string)
	pbftMessageJSONLen := int32(len(pbftMessageJSON))
	for _, tcpPortValString := range p.PbftPeerPorts { //! MAKE THIS GENERAL METHOD TO BE USED IN BOTH SEND FUNCTION
		tcpPortValInt, er := strconv.Atoi(tcpPortValString)
		if er != nil {
			p.Logger.Errorw("[node_source] Error: SendBlockToPeer::", zap.Error(er))
			return nil, err.ErrIntParse
		}

		if currentMappedTCPPort != tcpPortValInt {
			conn, er := connectionWithRetries(tcpPortValString, 75)
			if er != nil {
				p.Logger.Errorw("[node_source] Error: SendBlockToPeer::", zap.Error(er))
				ackMap[tcpPortValInt] = er.Error()
				continue
			}
			tempEntity := entity.PBFTMessage{}
			json.Unmarshal(pbftMessageJSON, &tempEntity)
			p.Logger.Infoln("[node_source] SENT TO PORT : ", tcpPortValInt)
			p.Logger.Infoln("[node_source] SENT Data: ", tempEntity)
			// conn.SetDeadline(time.Now().Add(3 * time.Second))//3 sec deadline may be imp if i dont want this node to be hanging if other side  doesnot respond
			_, er = conn.Write([]byte(strconv.Itoa(int(pbftMessageJSONLen)) + "\n"))
			if er != nil {
				p.Logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(nil, er.Error())

			}
			_, er = conn.Write(pbftMessageJSON)
			if er != nil {
				p.Logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(nil, er.Error())
			}

			// ackbuf := make([]byte, 4096)
			ackStrResponse, er := bufio.NewReader(conn).ReadString('\n')
			if er != nil {
				p.Logger.Errorw("[node_source] Error: SendBlockToPeer::", er)
				return nil, err.ErrWithMoreInfo(err.ErrTcpListen, er.Error())
			}
			ackMap[tcpPortValInt] = strings.TrimSpace(ackStrResponse)
			time.Sleep(100 * time.Millisecond)
			defer func() {

				conn.Close()
				p.Logger.Infow("CONNECTION CLOSED BY SENDER")
			}()
		}
	}
	return ackMap, nil
}

func (p *PBFTService) CheckPBFTMessageStructure(operationCounter int, pbftMessage entity.PBFTMessage) error {
	if pbftMessage.OperationID < 0 {
		return err.ErrNumberMismatch
	}
	if pbftMessage.VerificationType != enum.PREPREPARE &&
		pbftMessage.VerificationType != enum.PREPARE &&
		pbftMessage.VerificationType != enum.INITIAL &&
		pbftMessage.VerificationType != enum.COMMIT {
		return err.ErrInvalidType
	}

	// if pbftMessage.QRVerificationRequestData.CertificateHash == "" {
	// 	return err.ErrInvalidType
	// }
	// if pbftMessage.QRVerificationRequestData.ClientID == "" {
	// 	return err.ErrInvalidType
	// }
	// if pbftMessage.QRVerificationRequestData.Timestamp <= 0 {
	// 	return err.ErrInvalidType
	// }
	return nil
}

func (p *PBFTService) MarshalAndBroadcast(pbftMessage entity.PBFTMessage, pbftPorts []string, currentMappedTCPPort int) (map[int]string, error) {
	pbftMessageJSON, er := json.Marshal(pbftMessage)
	if er != nil {
		return nil, err.ErrMarshaling
	}
	returnMap, er := p.BroadcastToOtherPeers(pbftPorts, currentMappedTCPPort, pbftMessageJSON)
	p.Logger.Infof("[node_source] Debug: return map::", "return map", returnMap)
	if er != nil {
		return nil, er
	}
	return returnMap, nil
}

func connectionWithRetries(tcpPortValString string, retries int) (net.Conn, error) {
	for i := 0; i < retries; i++ {
		conn, er := net.Dial("tcp", "localhost:"+tcpPortValString)
		if er == nil {
			return conn, nil
		}

	}
	return nil, err.ErrConnection
}

func (p *PBFTService) IsSequenceNumberCorrect(operationCounter *int, pbftMessage entity.PBFTMessage) bool {

	if *operationCounter != pbftMessage.OperationID && *operationCounter < pbftMessage.OperationID {
		p.Logger.Debug("operationid", pbftMessage.OperationID)
		p.Logger.Debug("operationcounter", *operationCounter)
		return false
	}
	return true
}
