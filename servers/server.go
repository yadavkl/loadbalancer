package servers

import (
	"net/http"
	"net/url"
)

type Server interface {
	//Address() string
	SetAlive(bool)
	IsAlive() bool
	GetURL() *url.URL
	GetActiveConnections() int
	Serve(rw http.ResponseWriter, req *http.Request)
}

type ServerPool interface {
	GetBackends() []Server
	GetNextValidPeer() Server
	AddBackend(Server)
	GetServerPoolSize() int
}
