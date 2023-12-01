package servers

import (
	"context"
	"net"
	"net/url"
	"time"
)

func IsBackendAlive(ctx context.Context, aliveChannel chan bool, u *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}

func HealthCheck(ctx context.Context, s ServerPool) {
	aliveChannel := make(chan bool)
	for _, b := range s.GetBackends() {
		b := b
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		go IsBackendAlive(requestCtx, aliveChannel, b.GetURL())

		select {
		case <-ctx.Done():
			return
		case alive := <-aliveChannel:
			b.SetAlive(alive)
		}
	}
}

func LaunchHealthCheck(ctx context.Context, s ServerPool) {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			go HealthCheck(ctx, s)
		case <-ctx.Done():
			return
		}
	}
}
