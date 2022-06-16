// Package proxy establishes a reverse proxy with sidecar backend
package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
)

var (
	// Port is the proxy server port
	Port = ""
	// Backend is the fully qualified  backend url ip:port
	Backend = ""
)

// Proxy struct containing information about single host proxy
type Proxy struct {
	backend *url.URL
	proxy   *httputil.ReverseProxy
}

// NewReverseProxy initialize backend with single host proxy server
func NewReverseProxy() (*Proxy, error) {
	backendServer, err := url.Parse(Backend)
	if err != nil {
		log.Panic(err.Error())
		return nil, errors.Wrapf(err, "URL Parsing failed for: %s", Backend)
	}

	if backendServer.Port() == "" {
		log.Panicf("protocol and port are required in %v", Backend)
	}

	log.Infof("Proxy server is running on: %v", Port)
	log.Infof("Backend server is running on: %v", Backend)

	revProxy := httputil.NewSingleHostReverseProxy(backendServer)
	revProxy.Transport = &ochttp.Transport{}

	return &Proxy{
		backend: backendServer,
		proxy:   revProxy,
	}, nil
}

// ServeHTTPProxy HTTP proxy
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
