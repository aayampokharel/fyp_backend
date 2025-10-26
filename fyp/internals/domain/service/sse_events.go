package service

import (
	"project/internals/domain/entity"
	"sync"
)

type SSEManager struct {
	clients map[string]chan entity.Institution
	mu      sync.RWMutex
}

func (m *SSEManager) AddClient(adminID string) <-chan entity.Institution {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan entity.Institution)
	m.clients[adminID] = ch
	return ch
}
func (m *SSEManager) RemoveClient(adminID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	close(m.clients[adminID])
	delete(m.clients, adminID)
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
