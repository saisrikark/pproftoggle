package pproftoggle_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/saisrikark/pproftoggle"
)

func TestStart(t *testing.T) {
	ppfs, _ := pproftoggle.NewServer(
		pproftoggle.ServerConfig{
			HttpServer: &http.Server{
				Addr: ":8080",
			},
		})

	go func() {
		fmt.Println("before shutting down")
		time.Sleep(time.Second * 10)
		fmt.Println("after sleeping")
		ppfs.Shutdown()
	}()

	ppfs.Listen()
}
