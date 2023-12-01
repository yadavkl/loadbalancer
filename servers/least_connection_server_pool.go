package servers

import (
	"sync"
)

type leastConnectionServerPool struct {
	backends []Server
	mux      sync.RWMutex
}

func (s *leastConnectionServerPool) GetBackends() []Server {
	return s.backends
}

func (s *leastConnectionServerPool) GetNextValidPeer() Server {
	var leastConnectedPeer Server
	s.mux.RLock()
	defer s.mux.RUnlock()
	for _, b := range s.backends {
		if b.IsAlive() {
			leastConnectedPeer = b
			break
		}
	}
	for _, b := range s.backends {
		if !b.IsAlive() {
			continue
		}
		if b.GetActiveConnections() < leastConnectedPeer.GetActiveConnections() {
			leastConnectedPeer = b
		}
	}
	return leastConnectedPeer
}

func (s *leastConnectionServerPool) AddBackend(server Server) {
	s.mux.Lock()
	s.backends = append(s.backends, server)
	s.mux.Unlock()
}

func (s *leastConnectionServerPool) GetServerPoolSize() int {
	return len(s.backends)
}
