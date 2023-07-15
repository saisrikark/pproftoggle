package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/saisrikark/pproftoggle"
	"github.com/saisrikark/pproftoggle/rules"
)

// To use this example
// ensure dependencies are present
// run the main.go file from examples/toggler
// use "go run main.go"
// change the value of enablepprof in example.yaml
// to anything else and the server will shutdown
// toggle it again and it will start

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
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
		log.Println("received error while trying to create new toggler", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := toggler.Serve(ctx); err != nil {
			log.Println("received error while trying to serve using toggler", err)
		}
	}()
	time.Sleep(time.Minute * 10)

	cancel()
	wg.Wait()
}
