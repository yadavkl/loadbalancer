package servers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type backend struct {
	url          *url.URL
	alive        bool
	mux          sync.RWMutex
	connections  int
	reverseProxy *httputil.ReverseProxy
}

func NewBackend(url *url.URL, rp *httputil.ReverseProxy) *backend {
	return &backend{
		url:          url,
		alive:        true,
		reverseProxy: rp,
	}
}

func (b *backend) GetActiveConnections() int {
	b.mux.RLock()
	connections := b.connections
	b.mux.RUnlock()
	return connections
}

func (b *backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	b.mux.Unlock()
}

func (b *backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.alive
	b.mux.RUnlock()
	return alive
}
func (b *backend) GetURL() *url.URL {
	b.mux.RLock()
	url := b.url
	b.mux.RUnlock()
	return url
}

func (b *backend) Serve(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		b.mux.Lock()
		b.connections--
		b.mux.Unlock()
	}()
	b.mux.Lock()
	b.connections++
	b.mux.Unlock()
	b.reverseProxy.ServeHTTP(rw, req)
}
