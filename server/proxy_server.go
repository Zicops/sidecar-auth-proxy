// Package server implementation for reverse proxy
package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/zicops/sidecar-auth-proxy/handlers/authz"
	"github.com/zicops/sidecar-auth-proxy/proxy"

	graceful "gopkg.in/tylerb/graceful.v1"
)

//ProxyServer struct to handle start/shut of proxy
type ProxyServer struct {
	httpServer *graceful.Server
}
type checkerHandler func(http.Handler) http.Handler

var (
	runningServer *ProxyServer
	//AuthZHandler  very typical of auth checks multiple handlers can be served
	AuthZHandler checkerHandler = authz.Check
)

func getHandler(mode string, h http.Handler) http.Handler {
	handler := h
	modes := strings.Split(mode, "|")
	for i := len(modes) - 1; i >= 0; i-- {
		switch modes[i] {
		case "authz":
			handler = AuthZHandler(handler)
		default:
			log.Panicf("required -mode not given or illegal value")
		}
	}
	return handler
}

// ProxyServerStart handles the incoming request validation and rountripping to backend
func ProxyServerStart(ctx context.Context, port string, modes *string, proxyProtocol *proxy.Proxy) {
	log.Infof("Proxy Server Initiate: Successfully")
	handler := getHandler(*modes, proxyProtocol)

	http.Handle("/", handler)
	addr := fmt.Sprintf(":%v", port)

	// init server
	srv := &graceful.Server{
		Timeout: 30 * time.Second,
		//BeforeShutdown:    beforeShutDown,
		ShutdownInitiated: shutdownInitiated,
		Server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  time.Duration(30) * time.Second,
			WriteTimeout: time.Duration(30) * time.Second,
		},
	}

	// start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("ProxyServer: Failed to start server : %s", err.Error())
		}
	}()
	// wait for the context to be canceled
	<-ctx.Done()

	log.Infof("Sidecar Auth Proxy: context done is called %v", time.Now())
	srv.Shutdown(ctx)
	log.Infof("Sidecar Auth Proxy: Server shutdown at %v", time.Now())
}

// ProxyServerShutDown shutdown on command
func ProxyServerShutDown(ctx context.Context) {
	runningServer.httpServer.Shutdown(ctx)
}

func shutdownInitiated() {
	log.Infof("Sidecar Auth Proxy: Shutting down server at %v", time.Now())
}
