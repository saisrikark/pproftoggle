package pproftoggle

import (
	"errors"
	"net/http"
	"net/http/pprof"
	"sync/atomic"
)

const (
	DEFAULT_PPROF_DEBUG_PREFIX = "/debug/pprof/"
)

type pprofServer struct {
	isUp       *atomic.Bool
	httpServer *http.Server
}

type pprofServerConfig struct {
	Address        string
	EndpointPrefix string
}

func newpprofServer(cfg pprofServerConfig) (pprofServer, error) {
	var prefix string
	mux := http.NewServeMux()

	prefix = cfg.EndpointPrefix + DEFAULT_PPROF_DEBUG_PREFIX
	mux.HandleFunc(prefix, pprof.Index)
	mux.HandleFunc(prefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"profile", pprof.Profile)
	mux.HandleFunc(prefix+"symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"trace", pprof.Trace)

	httpServer := &http.Server{}
	httpServer.Addr = cfg.Address
	httpServer.Handler = mux

	return pprofServer{
		httpServer: httpServer,
		isUp:       &atomic.Bool{},
	}, nil
}

func (ppfs *pprofServer) start() error {
	if ppfs.isRunning() {
		return errors.New("server already up")
	}

	ppfs.isUp.Store(true)
	defer ppfs.isUp.Store(false)

	return ppfs.httpServer.ListenAndServe()
}

func (ppfs *pprofServer) stop() error {
	if !ppfs.isRunning() {
		return errors.New("server is down")
	}

	return ppfs.httpServer.Close()
}

func (ppfs *pprofServer) isRunning() bool {
	return ppfs.isUp.Load()
}
