package pproftoggle

import (
	"net/http"
	"net/http/pprof"
	"sync/atomic"
)

const (
	pprofPrefix = "/debug/pprof/"
)

type pprofServer struct {
	isUp       *atomic.Bool
	httpServer *http.Server
}

type ServerConfig struct {
	HttpServer     *http.Server
	EndpointPrefix string
}

func NewServer(cfg ServerConfig) (*pprofServer, error) {
	var prefix = cfg.EndpointPrefix + pprofPrefix
	var mux = http.NewServeMux()

	mux.HandleFunc(prefix, pprof.Index)
	mux.HandleFunc(prefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"profile", pprof.Profile)
	mux.HandleFunc(prefix+"symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"trace", pprof.Trace)
	cfg.HttpServer.Handler = mux

	return &pprofServer{
		httpServer: cfg.HttpServer,
		isUp:       &atomic.Bool{},
	}, nil
}

func (ppfs *pprofServer) Listen() error {
	if ppfs.IsRunning() {
		return nil
	}

	ppfs.isUp.Store(true)
	defer ppfs.isUp.Store(false)

	return ppfs.httpServer.ListenAndServe()
}

func (ppfs *pprofServer) Shutdown() error {
	if !ppfs.IsRunning() {
		return nil
	}

	return ppfs.httpServer.Close()
}

func (ppfs *pprofServer) IsRunning() bool {
	return ppfs.isUp.Load()
}
