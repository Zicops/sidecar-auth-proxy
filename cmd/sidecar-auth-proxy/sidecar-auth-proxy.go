// Package main for authorization service proxy
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/zicops/sidecar-auth-proxy/proxy"
	"github.com/zicops/sidecar-auth-proxy/server"
)

func main() {

	log.Infof("Starting sidecar-auth-proxy as sidecar container.")
	flag.StringVar(&proxy.Port, "port", "", "Expose a port to accept HTTP/1.x connections")
	flag.StringVar(&proxy.Backend, "backend", "", "Change the application server address to which to proxies the requests")
	mode := flag.String("mode", "", "What to proxy authn, authz etc. They may be piped together e.g authn|authz")

	// parse flags to get port, backend and mode
	flag.Parse()
	if flag.NFlag() < 3 {
		flag.PrintDefaults()
		log.Panicf("Expected 3 arguments, got %v", os.Args[1:])
	}

	// Channels for os signals
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// initializing single host proxy
	proxyProtocol, err := proxy.NewReverseProxy()

	if err != nil {
		message := "Auth Proxy: Proxy server init failed"
		log.Panicf("%s: %s", message, err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	//Start the server and establish communication with backend
	go server.ProxyServerStart(ctx, proxy.Port, mode, proxyProtocol)
	log.Infof("Proxy Server Started.")
	sig := <-sigC
	log.Infof("Received %d, shutting down", sig)

	defer cancel()
	server.ProxyServerShutDown(ctx)
}
