package servers

import (
	"errors"
	"sync"

	"github.com/yadavkl/loadbalancer/utils"
)

type roundRobinServerPool struct {
	backends []Server
	mux      sync.RWMutex
	current  int
}

func (s *roundRobinServerPool) GetBackends() []Server {
	return s.backends
}

func (s *roundRobinServerPool) GetNextValidPeer() Server {
	for i := 0; i < s.GetServerPoolSize(); i++ {
		nextPeer := s.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil
}

func (s *roundRobinServerPool) AddBackend(server Server) {
	s.mux.Lock()
	s.backends = append(s.backends, server)
	s.mux.Unlock()
}

func (s *roundRobinServerPool) GetServerPoolSize() int {
	return len(s.backends)
}

func (s *roundRobinServerPool) Rotate() Server {
	s.mux.Lock()
	s.current = (s.current + 1) % s.GetServerPoolSize()
	s.mux.Unlock()
	return s.backends[s.current]
}

func NewServerPool(strategy utils.LBStrategy) (ServerPool, error) {
	switch strategy {
	case utils.RoundRobin:
		return &roundRobinServerPool{
			backends: make([]Server, 0),
			current:  0,
		}, nil
	case utils.LeastConnection:
		return &leastConnectionServerPool{
			backends: make([]Server, 0),
		}, nil
	default:
		return nil, errors.New("Invalid strategy")
	}
}
