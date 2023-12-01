package lb

import (
	"net/http"

	"github.com/yadavkl/loadbalancer/servers"
)

const (
	RETRY_ATTEMPTED int = 0
)

type LoadBalancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

type loadBalancer struct {
	serverPool servers.ServerPool
}

func NewLoadBalancer(serverPool servers.ServerPool) *loadBalancer {
	return &loadBalancer{
		serverPool: serverPool,
	}
}

func (lb *loadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	peer := lb.serverPool.GetNextValidPeer()
	if peer != nil {
		peer.Serve(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func AllowRetry(r *http.Request) bool {
	if _, ok := r.Context().Value(RETRY_ATTEMPTED).(bool); ok {
		return false
	}
	return true
}
