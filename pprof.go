package pproftoggle

import (
	"fmt"
	"net/http"
	"net/http/pprof"
)

const (
	DEFAULT_PPROF_DEBUG_PREFIX = "/debug/pprof/"
)

type pprofServer struct {
	httpServer *http.Server
}

type pprofServerConfig struct {
	Address    string // address at which to host the server panics if empty
	UserPrefix string // prefix to the /debug/pprof endpoint default not set
}

func newpprofServer(cfg pprofServerConfig) (*pprofServer, error) {
	var prefix string
	mux := http.NewServeMux()
	httpServer := &http.Server{}

	if cfg.Address == "" {
		panic("empty address")
	}

	if cfg.UserPrefix != "" {
		prefix = cfg.UserPrefix
	} else {
		prefix = DEFAULT_PPROF_DEBUG_PREFIX
	}
	mux.HandleFunc(prefix, pprof.Index)
	mux.HandleFunc(prefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"profile", pprof.Profile)
	mux.HandleFunc(prefix+"symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"trace", pprof.Trace)

	httpServer.Addr = cfg.Address
	httpServer.Handler = mux

	return &pprofServer{
		httpServer: httpServer,
	}, nil
}

func (ppfs *pprofServer) start() error {
	if err := ppfs.httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("closing server")
		} else {
			fmt.Println("error while trying to server request")
			return err
		}
	}
	return nil
}

func (ppfs *pprofServer) stop() error {
	return ppfs.httpServer.Close()
}
