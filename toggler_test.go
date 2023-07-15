package pproftoggle_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/saisrikark/pproftoggle"
	"github.com/saisrikark/pproftoggle/rules"
)

var (
	localhost               = "127.0.0.1"
	port                    = "8080"
	baseURL                 = "http://" + localhost + ":" + port + "/debug/pprof"
	endpoint                = "/"
	envKey                  = "ENABLE_PPROF"
	envVal                  = "true"
	pollInterval            = time.Second
	acceptablePollDeviation = time.Millisecond * 200
)

func TestServe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	toggler, err := pproftoggle.NewToggler(
		pproftoggle.Config{
			PollInterval: pollInterval,
			HttpServer: &http.Server{
				Addr: ":" + port,
			},
			Rules: []pproftoggle.Rule{
				rules.EnvVarRule{
					Key:   envKey,
					Value: envVal,
				},
			},
		})
	if err != nil {
		t.Errorf("%s", err.Error())
		t.FailNow()
	}

	go func() {
		if err := toggler.Serve(ctx); err != nil {
			t.Errorf("%s", err.Error())
		}
	}()

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			shouldBeRunning := (i%2 == 0)
			if !shouldBeRunning {
				if err := os.Unsetenv(envKey); err != nil {
					t.Errorf("%s", err.Error())
				}
			} else {
				if err := os.Setenv(envKey, envVal); err != nil {
					t.Errorf("%s", err.Error())
				}
			}

			// expecting a toggle after the poll interval elapses
			// waiting for slightly longer to avoid race conditions
			time.Sleep(pollInterval + acceptablePollDeviation)

			// check if IsUp is showing expected status
			if shouldBeRunning && !toggler.IsUp(ctx) {
				t.Error("not running when it should be")
			} else if !shouldBeRunning && toggler.IsUp(ctx) {
				t.Error("running when it shouldn't be")
			}

			if err := requests.URL(endpoint).BaseURL(baseURL).Fetch(ctx); err != nil && shouldBeRunning {
				t.Errorf("%s", err.Error())
			} else if err == nil && !shouldBeRunning {
				t.Error("unexpected response")
			}
		})
	}
}
