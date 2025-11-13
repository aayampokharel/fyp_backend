package repository

import "project/internals/domain/entity"

type INodeRepository interface {
	SendBlockToPeer(block entity.Block, currentPort int) (map[int]string, error)
	ReceiveBlockFromPeer(listenPort int) (*entity.Block, error)
	// SendPBFTMessageToPeer(pbftMessage entity.PBFTMessage, leaderNode int, currentMappedTCPPort int) (map[int]string, error)
	// ReceivePBFTMessageToPeer(listenPort int, leaderPort int) (*entity.PBFTMessage, error)
}
