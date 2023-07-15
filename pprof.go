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
	prefix         string
	isUp           *atomic.Bool
	httpServer     *http.Server
	userHttpServer *http.Server
}

type ServerConfig struct {
	// HttpServer is the http server whose configuration is
	// used to start pprof
	HttpServer *http.Server
	// EndpointPrefix is used in order to extend the endpoint beyond /debug/pprof/
	// example if EndpointPrefix is "/India/server1/"
	// the final endpoint will begin with "/India/server1/debug/pprof"
	EndpointPrefix string
}

// newHttpServer copies all values from userHttpServer into one we can use
func newHttpServer(prefix string, userHttpServer *http.Server) *http.Server {
	prefix = prefix + pprofPrefix
	var mux = http.NewServeMux()

	mux.HandleFunc(prefix, pprof.Index)
	mux.HandleFunc(prefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"profile", pprof.Profile)
	mux.HandleFunc(prefix+"symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"trace", pprof.Trace)

	svr := &http.Server{
		Addr:                         userHttpServer.Addr,
		Handler:                      mux,
		DisableGeneralOptionsHandler: userHttpServer.DisableGeneralOptionsHandler,
		TLSConfig:                    userHttpServer.TLSConfig,
		ReadTimeout:                  userHttpServer.ReadTimeout,
		WriteTimeout:                 userHttpServer.WriteTimeout,
		IdleTimeout:                  userHttpServer.IdleTimeout,
		MaxHeaderBytes:               userHttpServer.MaxHeaderBytes,
		TLSNextProto:                 userHttpServer.TLSNextProto,
		ConnState:                    userHttpServer.ConnState,
		ErrorLog:                     userHttpServer.ErrorLog,
		BaseContext:                  userHttpServer.BaseContext,
		ConnContext:                  userHttpServer.ConnContext,
	}

	return svr
}

func NewServer(cfg ServerConfig) (*pprofServer, error) {
	if cfg.HttpServer == nil {
		return nil, errors.New("http server not configured")
	}

	return &pprofServer{
		prefix:         cfg.EndpointPrefix,
		isUp:           &atomic.Bool{},
		userHttpServer: cfg.HttpServer,
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
		ppfs.httpServer = newHttpServer(ppfs.prefix, ppfs.userHttpServer)
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
	if !ppfs.IsRunning() || ppfs.httpServer == nil {
		return nil
	}

	return ppfs.httpServer.Shutdown(context.Background())
}

func (ppfs *pprofServer) IsRunning() bool {
	return ppfs.isUp.Load()
}
