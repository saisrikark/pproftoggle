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
	isUp           *atomic.Bool
	httpServer     *http.Server
	userHttpServer *http.Server
}

type ServerConfig struct {
	HttpServer     *http.Server
	EndpointPrefix string
}

// newHttpServer copies all values from userHttpServer into one we can use
func newHttpServer(userHttpServer *http.Server) *http.Server {
	var prefix = pprofPrefix
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
		userHttpServer: cfg.HttpServer,
		isUp:           &atomic.Bool{},
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
		ppfs.httpServer = newHttpServer(ppfs.userHttpServer)
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
