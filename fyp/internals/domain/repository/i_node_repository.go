package repository

import "project/internals/domain/entity"

type INodeRepository interface {
	SendBlockToPeer(block entity.Block, currentPort int) (map[int]string, error)
	ReceiveBlockFromPeer(listenPort int) (*entity.Block, error)
}
