package pproftoggle_test

import (
	"context"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/saisrikark/pproftoggle"
)

var (
	TestingPort = "3125"
	BASE_URL    = "http://localhost:" + TestingPort + "/debug/pprof"
)

func TestListen(t *testing.T) {
	t.Run("HttpServerNotSpecified", func(t *testing.T) {
		_, err := pproftoggle.NewServer(pproftoggle.ServerConfig{})
		if err == nil {
			t.Errorf("should lead to an error")
		} else {
			t.Logf("as expected could not create a server with empty config err:[%s]", err.Error())
		}
	})

	t.Run("HttpServerShouldNotStartAtPortAlreadyBeingUsed", func(t *testing.T) {
		var lc net.Listener
		var wg sync.WaitGroup
		var err error

		wg.Add(1)
		go func() {
			lc, err = net.Listen("tcp", ":"+TestingPort)
			if err != nil {
				t.Errorf("unable to create listner at port:[%s] err:[%s]", TestingPort, err.Error())
			}
			wg.Done()
		}()
		wg.Wait()
		if lc != nil {
			defer lc.Close()
		}

		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + TestingPort},
		})
		if err != nil {
			t.Errorf("unable to initialise new server")
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		err = ppfs.Listen(ctx)
		if err == nil {
			t.Errorf("serving requests at port should not possible as port [%s] is blocked", TestingPort)
		} else {
			t.Logf("as expected encountered error:[%s] while trying to server request on a blocked port", err.Error())
		}
	})

	t.Run("FetchFromEndpoint", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + TestingPort},
		})
		if err != nil {
			t.Errorf("unable to initialise new server")
		}

		go func() {
			if err = ppfs.Listen(ctx); err != nil {
				t.Errorf("unable to listen err:[%s]", err.Error())
			}
		}()

		for i := 0; i < 5; i++ {
			time.Sleep(time.Millisecond * 10)
			if ppfs.IsRunning() {
				break
			}
		}

		resp := ""
		if err := requests.URL("/").BaseURL(BASE_URL).ToString(&resp).Fetch(context.Background()); err != nil {
			t.Errorf("unable to fetch from /trace endpoint err:[%s]", err)
		} else {
			t.Logf("received a respose from the endpoint truncated resp\n[%s]", resp[0:30])
		}

		cancel()
	})
}

func TestShutdown(t *testing.T) {

}
