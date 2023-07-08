package pproftoggle

import (
	"context"
	"fmt"
	"time"
)

type toggler struct {
	pollInterval time.Duration
	ppfs         pprofServer
	rules        []Rule
}

type Config struct {
	// Address at which to host the pprof server
	// defaults to ":8080"
	Address string
	// EndpointPrefix is used to extend the path to access pprof
	// by default it is served at /debug/pprof/...
	// if given as "/extra" endpoint it extended to /extra/debug/pprof/...
	EndpointPrefix string
	// PollInterval is duration between the end of one poll and the beginning of another
	// minimum value is defaulted to 1s
	PollInterval time.Duration
	// Rules is a list of conditions based on which
	// we decide whether to start or stop the pprof server
	// executed in same order specified
	Rules []Rule
}

// NewTogger returns an instance of the toggler
// used to perform pprof toggle operations
// pollInterval is the wait time between
// the end of a match and the beggining of a new one
// paincs if no rules are configured
func NewToggler(cfg Config) (*toggler, error) {

	if cfg.PollInterval < time.Second {
		cfg.PollInterval = time.Second
	}

	if cfg.Rules == nil || len(cfg.Rules) == 0 {
		panic("no rules configured")
	}

	pprofServer, err := newpprofServer(pprofServerConfig{
		Address:        cfg.Address,
		EndpointPrefix: cfg.EndpointPrefix,
	})
	if err != nil {
		return nil, err
	}

	return &toggler{
		pollInterval: cfg.PollInterval,
		ppfs:         pprofServer,
		rules:        cfg.Rules,
	}, nil
}

// Serve continuously polls for the given conditions
// and starts the pprof server if conditions match
// this is a blocking operation that return when ctx is cancelled
// or when an error is hit
func (pt *toggler) Serve(ctx context.Context) error {
	exit := false
	timer := time.NewTimer(pt.pollInterval)
	for {
		select {
		case <-timer.C:
			fmt.Println("SRIKAR HERE ", time.Now())
			timer = time.NewTimer(pt.pollInterval)
			// TODO check if condition matches
			// if matches and server is not up
			// bring it up
			// if doesn't match and server is up
			// bring it down
		case <-ctx.Done():
			fmt.Println("DONE", time.Now())
			exit = true
		}
		if exit {
			break
		}
	}

	return nil
}

// ForceStart brings up the pprof server is not already up
// once used polling will no longer work
func (pt *toggler) ForceStart(ctx context.Context) error {
	if pt.IsUp(ctx) {
		return nil
	}

	return nil
}

// ForceStop brings down the pprof server
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
	return pt.ppfs.isRunning()
}
