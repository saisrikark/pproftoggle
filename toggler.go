package pproftoggle

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"
)

type toggler struct {
	isRunning    atomic.Bool
	prefix       string
	pollInterval time.Duration
	httpServer   *http.Server
}

// NewTogger returns an instance of the toggler
// used to perform pprof toggle operations
func NewToggler(pollInterval time.Duration) *toggler {
	httpServer := &http.Server{
		Addr: ":8080",
	}
	return &toggler{
		prefix:       "",
		pollInterval: pollInterval,
		httpServer:   httpServer,
	}
}

// Serve continuously polls for the given conditions
// and toggles the pprof server if conditions match
func (pt *toggler) Serve(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Start brings up the pprof server
// once used polling will no longer work
func (pt *toggler) ForceStart(ctx context.Context) error {
	return nil
}

// Stop brings down the pprof server
// once used polling will no longer work
func (pt *toggler) ForceStop(ctx context.Context) error {
	return nil
}

// Toggle either brings up or shuts down the pprof server depending on
// the current state
// once used polling will no longer work
func (pt *toggler) Toggle(ctx context.Context) error {
	return nil
}

// IsUp returns the running status of the server hosting pprof endpoints
func (pt *toggler) IsUp(ctx context.Context) bool {
	return pt.isRunning.Load()
}
