package pproftoggle_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/saisrikark/pproftoggle"
)

func TestServe(t *testing.T) {
	toggler, err := pproftoggle.NewToggler(
		pproftoggle.Config{
			PollInterval: time.Second * 3,
			HttpServer: &http.Server{
				Addr: ":8080",
			},
			Rules: []pproftoggle.Rule{
				pproftoggle.EnvVarRule{
					Key:   "abcd",
					Value: "efgh",
				},
			},
		})
	if err != nil {
		fmt.Println(err)
	}

	set := true
	go func() {
		for {
			time.Sleep(time.Second * 2)
			if set {
				fmt.Println("setting")
				os.Setenv("abcd", "efgh")
				set = false
			} else {
				fmt.Println("unsetting")
				os.Unsetenv("abcd")
				set = true
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	if err := toggler.Serve(ctx); err != nil {
		fmt.Println(err)
	}
}
