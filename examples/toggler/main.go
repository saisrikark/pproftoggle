package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/saisrikark/pproftoggle"
	"github.com/saisrikark/pproftoggle/rules"
)

func main() {
	toggler, err := pproftoggle.NewToggler(
		pproftoggle.Config{
			PollInterval: time.Second * 1,
			HttpServer: &http.Server{
				Addr: ":8080",
			},
			Rules: []pproftoggle.Rule{
				rules.EnvVarRule{
					Key:   "ENABLE_PPROF",
					Value: "true",
				},
				rules.SimpleYamlRule{
					Key:   "enablepprof",
					Value: "true",
					Path:  "example.yaml",
				},
			},
		})
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		toggler.Serve(ctx)
	}()

	time.Sleep(time.Minute * 10)

	cancel()
	wg.Wait()
}
