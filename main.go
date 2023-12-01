package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yadavkl/loadbalancer/lb"
	"github.com/yadavkl/loadbalancer/servers"
	"github.com/yadavkl/loadbalancer/utils"
)

func main() {
	logger := utils.InitLogger()
	defer logger.Sync()
	config, err := utils.GetLBConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverPool, err := servers.NewServerPool(utils.GetLBStrategy(config.Strategy))
	if err != nil {
		logger.Fatal(err.Error())
	}

	loadbalancer := lb.NewLoadBalancer(serverPool)

	for _, s := range config.Backends {
		endpoint, err := url.Parse(s)
		if err != nil {
			logger.Fatal(err.Error())
		}
		rp := httputil.NewSingleHostReverseProxy(endpoint)
		backendServer := servers.NewBackend(endpoint, rp)
		rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			backendServer.SetAlive(false)
			if !lb.AllowRetry(r) {
				return
			}
			loadbalancer.Serve(w, r.WithContext(context.WithValue(r.Context(), lb.RETRY_ATTEMPTED, true)))
		}
		serverPool.AddBackend(backendServer)
	}
	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", config.Port),
		Handler: http.HandlerFunc(loadbalancer.Serve),
	}
	fmt.Println(server.Addr)
	go servers.LaunchHealthCheck(ctx, serverPool)

	go func() {
		<-ctx.Done()
		shutDownCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
		if err := server.Shutdown(shutDownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal(err.Error())
	}
}
