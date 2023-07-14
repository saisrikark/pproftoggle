package pproftoggle

import (
	"context"
	"errors"
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

	if cfg.HttpServer == nil {
		return nil, errors.New("http server not configured")
	}

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

func (ppfs *pprofServer) Listen(ctx context.Context) error {
	var errs chan error = make(chan error, 1)

	if ppfs.IsRunning() {
		return nil
	}

	ppfs.isUp.Store(true)
	defer ppfs.isUp.Store(false)

	go func() {
		if err := ppfs.httpServer.ListenAndServe(); err != nil {
			errs <- err
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ppfs.httpServer.Close()
		case err := <-errs:
			return err
		}
	}
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
