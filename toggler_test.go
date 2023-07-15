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

func TestServe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	envKey := "ENABLE_PPROF"
	envVal := "true"
	pollInterval := time.Second

	toggler, err := pproftoggle.NewToggler(
		pproftoggle.Config{
			PollInterval: pollInterval,
			HttpServer: &http.Server{
				Addr: ":8080",
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
			time.Sleep(2 * pollInterval)

			// check if IsUp is showing expected status
			if shouldBeRunning && !toggler.IsUp(ctx) {
				t.Error("not running when it should be")
			} else if !shouldBeRunning && toggler.IsUp(ctx) {
				t.Error("running when it shouldn't be")
			}

			if err := requests.URL(ENDDPOINT).BaseURL(BASE_URL).Fetch(ctx); err != nil && shouldBeRunning {
				t.Errorf("%s", err.Error())
			} else if err == nil && !shouldBeRunning {
				t.Error("unexpected response")
			}
		})
	}
}
