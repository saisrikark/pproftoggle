package pproftoggle

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"log"

	"github.com/pkg/errors"
)

type toggler struct {
	pollInterval time.Duration
	ppfs         pprofServer
	rules        []Rule
	logger       *log.Logger
	canToggle    *atomic.Bool
	shouldBeUp   *atomic.Bool
	toggleChan   chan (bool)
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
	// Logger
	Logger *log.Logger
}

// NewToggler returns an instance of the toggler
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

	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "pproftoggle ", 0)
	}

	pprofServer, err := newpprofServer(pprofServerConfig{
		Address:        cfg.Address,
		EndpointPrefix: cfg.EndpointPrefix,
	})
	if err != nil {
		return nil, err
	}

	canToggle := &atomic.Bool{}
	canToggle.Swap(true)

	return &toggler{
		pollInterval: cfg.PollInterval,
		ppfs:         pprofServer,
		rules:        cfg.Rules,
		logger:       cfg.Logger,
		canToggle:    canToggle,
		shouldBeUp:   &atomic.Bool{},
		toggleChan:   make(chan bool, 1),
	}, nil
}

// Serve continuously polls for the given conditions
// and starts the pprof server if conditions match
// this is a blocking operation that return when ctx is cancelled
// or when an error is hit
func (pt *toggler) Serve(ctx context.Context) error {
	exit := false
	timer := time.NewTimer(pt.pollInterval)
	errs := make(chan error, 1)

	asyncStart := func() {
		pt.logger.Println("starting pprof server")
		if err := pt.ppfs.start(); err != nil && err != http.ErrServerClosed {
			errs <- errors.Wrap(err, "unable to start pprof server")
		}
	}

	asyncStop := func() {
		pt.logger.Println("stopping pprof server")
		if err := pt.ppfs.stop(); err != nil {
			errs <- errors.Wrap(err, "unable to stop pprof server")
		}
	}

	stop := func() error {
		pt.logger.Println("stopping pprof server")
		if err := pt.ppfs.stop(); err != nil {
			return errors.Wrap(err, "unable to stop pprof server")
		}
		return nil
	}

	for {
		select {
		case err := <-errs:
			if pt.IsUp(ctx) {
				stop()
			}
			return err
		case <-ctx.Done():
			if pt.IsUp(ctx) {
				stop()
			}
			exit = true
		case <-timer.C:
			pt.logger.Println("performing rule check")
			status, err := getStatus(pt.rules)
			if err != nil {
				pt.logger.Printf("error while trying to fetch status %s", err.Error())
				return err
			}
			if pt.canToggle.Load() {
				if status.hasMatched && !pt.IsUp(ctx) {
					go asyncStart()
				} else if !status.hasMatched && pt.IsUp(ctx) {
					go asyncStop()
				}
			}
			timer = time.NewTimer(pt.pollInterval)
		case <-pt.toggleChan:
			if !pt.canToggle.Load() {
				if pt.shouldBeUp.Load() && !pt.IsUp(ctx) {
					asyncStart()
				} else if !pt.shouldBeUp.Load() && !pt.IsUp(ctx) {
					asyncStop()
				}
			}
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
	pt.canToggle.Store(false)
	pt.shouldBeUp.Store(true)

	if pt.IsUp(ctx) {
		return nil
	}

	pt.toggleChan <- true
	return nil
}

// ForceStop brings down the pprof server
// once used polling will no longer work
func (pt *toggler) ForceStop(ctx context.Context) error {
	pt.canToggle.Store(false)
	pt.shouldBeUp.Store(false)
	pt.toggleChan <- true
	return nil
}

// Toggle either brings up or shuts down the pprof server
// depending on the current state
// once used polling will no longer work
func (pt *toggler) Toggle(ctx context.Context) error {
	pt.canToggle.Store(false)
	pt.shouldBeUp.Store(!pt.IsUp(ctx))
	pt.toggleChan <- true
	return nil
}

// IsUp returns the running status of the server hosting pprof endpoints
func (pt *toggler) IsUp(ctx context.Context) bool {
	return pt.ppfs.isRunning()
}
