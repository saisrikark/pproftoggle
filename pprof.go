package pproftoggle

import (
	"errors"
	"net/http"
	"net/http/pprof"
	"sync"
	"sync/atomic"
)

const (
	DEFAULT_PPROF_DEBUG_PREFIX = "/debug/pprof/"
)

type pprofServer struct {
	mu         sync.Mutex
	isUp       *atomic.Bool
	httpServer *http.Server
}

type pprofServerConfig struct {
	HttpServer     *http.Server
	EndpointPrefix string
}

func newpprofServer(cfg pprofServerConfig) (pprofServer, error) {
	prefix := ""
	mux := http.NewServeMux()

	prefix = cfg.EndpointPrefix + DEFAULT_PPROF_DEBUG_PREFIX
	mux.HandleFunc(prefix, pprof.Index)
	mux.HandleFunc(prefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"profile", pprof.Profile)
	mux.HandleFunc(prefix+"symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"trace", pprof.Trace)

	cfg.HttpServer.Handler = mux
	return pprofServer{
		httpServer: cfg.HttpServer,
		isUp:       &atomic.Bool{},
	}, nil
}

func (ppfs *pprofServer) start() error {
	ppfs.mu.Lock()
	defer ppfs.mu.Unlock()

	if ppfs.isRunning() {
		return errors.New("server already up")
	}

	ppfs.isUp.Store(true)
	defer ppfs.isUp.Store(false)

	err := ppfs.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (ppfs *pprofServer) stop() error {
	ppfs.mu.Lock()
	defer ppfs.mu.Unlock()

	if !ppfs.isRunning() {
		return errors.New("server is down")
	}

	ppfs.isUp.Store(false)
	return ppfs.httpServer.Close()
}

func (ppfs *pprofServer) isRunning() bool {
	return ppfs.isUp.Load()
}
