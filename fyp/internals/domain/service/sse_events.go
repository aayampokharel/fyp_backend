package service

import (
	"project/internals/domain/entity"
	logger "project/package/utils/pkg"
	"sync"

	"go.uber.org/zap"
)

type SSEManager struct {
	clients map[string]chan<- entity.Institution
	mu      sync.RWMutex
	Logger  *zap.SugaredLogger
}

func NewSSEManager(channelMap map[string]chan<- entity.Institution) *SSEManager {
	return &SSEManager{clients: channelMap, mu: sync.RWMutex{}, Logger: logger.Logger}
}

func (m *SSEManager) AddClient(adminTokenID string) (ch chan entity.Institution) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logger.Info("[add_client] Info: addClient", adminTokenID)
	ch = make(chan entity.Institution)
	m.clients[adminTokenID] = ch
	return ch
}
func (m *SSEManager) RemoveClient(adminTokenID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	close(m.clients[adminTokenID])
	delete(m.clients, adminTokenID)
}

func (m *SSEManager) Broadcast(data entity.Institution) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ch := range m.clients {
		select {
		case ch <- data:
		default:
		}
	}
}
