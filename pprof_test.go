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
	PORT      = "3125"
	LOCALHOST = "127.0.0.1"
	BASE_URL  = "http://" + LOCALHOST + ":" + PORT + "/debug/pprof"
	ENDDPOINT = "/"
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
			lc, err = net.Listen("tcp", ":"+PORT)
			if err != nil {
				t.Errorf("unable to create listner at port:[%s] err:[%s]", PORT, err.Error())
			}
			wg.Done()
		}()
		wg.Wait()
		if lc != nil {
			defer lc.Close()
		}

		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + PORT},
		})
		if err != nil {
			t.Errorf("unable to initialise new server")
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		err = ppfs.Listen(ctx)
		if err == nil {
			t.Errorf("serving requests at port should not possible as port [%s] is blocked", PORT)
		} else {
			t.Logf("as expected encountered error:[%s] while trying to server request on a blocked port", err.Error())
		}
	})

	t.Run("FetchFromEndpoint", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + PORT},
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
			if isPortUsed(PORT) {
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
		if !isPortUsed(PORT) {
			t.Errorf("port:[%s] not being used", PORT)
		}

		resp := ""
		if err := requests.URL(ENDDPOINT).BaseURL(BASE_URL).ToString(&resp).Fetch(context.Background()); err != nil {
			t.Errorf("unable to fetch from /trace endpoint err:[%s]", err)
		} else {
			t.Logf("received a respose from the endpoint truncated resp\n[%s]", resp[0:30])
		}

		cancel()
	})
}

func TestShutdown(t *testing.T) {
	t.Run("ServerMustShutdown", func(t *testing.T) {
		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + PORT},
		})
		if err != nil {
			t.Errorf("unable to initialise new server")
		}

		go func() {
			err := ppfs.Listen(context.Background())
			if err != nil && err != http.ErrServerClosed {
				t.Errorf("errors while listening err:[%s]", err.Error())
			}
		}()

		for i := 0; i < 5; i++ {
			if isPortUsed(PORT) {
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
		if !isPortUsed(PORT) {
			t.Errorf("port:[%s] not being used", PORT)
		}

		resp := ""
		if err := requests.URL(ENDDPOINT).BaseURL(BASE_URL).ToString(&resp).Fetch(context.Background()); err != nil {
			t.Errorf("unable to fetch from endpoint err:[%s]", err)
		} else {
			t.Logf("received a respose from the endpoint truncated resp\n[%s]", resp[0:30])
		}

		ppfs.Shutdown()

		if err := requests.URL(ENDDPOINT).BaseURL(BASE_URL).ToString(&resp).Fetch(context.Background()); err == nil {
			t.Errorf("received a respose from the endpoint truncated resp\n[%s]", resp[0:30])
		} else {
			t.Logf("as expected unable to fetch from endpoint err:[%s]", err)
		}
	})
}

func TestIsRunning(t *testing.T) {
	ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
		HttpServer: &http.Server{Addr: ":" + PORT},
	})
	if err != nil {
		t.Errorf("unable to initialise new server")
	}

	go func() {
		err := ppfs.Listen(context.Background())
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("errors while listening err:[%s]", err.Error())
		}
	}()

	for i := 0; i < 5; i++ {
		if isPortUsed(PORT) {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	if !isPortUsed(PORT) {
		t.Errorf("port:[%s] not being used", PORT)
	}

	resp := ""
	if err := requests.URL(ENDDPOINT).BaseURL(BASE_URL).ToString(&resp).Fetch(context.Background()); err != nil {
		t.Errorf("unable to fetch from endpoint err:[%s]", err)
	} else {
		t.Logf("received a response from the endpoint truncated resp\n[%s]", resp[0:30])
		status := ppfs.IsRunning()
		if !status {
			t.Errorf("received response but IsRunning is showing [%v]", status)
		} else {
			t.Logf("as expected showing server running [%v]", status)
		}
	}

	if err := ppfs.Shutdown(); err != nil {
		t.Errorf("unexpected error while shutting down err:[%s]", err.Error())
	}

	if err := requests.URL(ENDDPOINT).BaseURL(BASE_URL).ToString(&resp).Fetch(context.Background()); err == nil {
		t.Errorf("received response from enpoint even after shutting down resp\n[%s]", resp[0:10])
	}

	if isPortUsed(PORT) {
		t.Errorf("port [%s] is still being used after shuttiing down", PORT)
	}
}

func isPortUsed(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(LOCALHOST, port), time.Millisecond*10)
	if err != nil {
		return false
	}

	if conn != nil {
		defer conn.Close()
		return true
	}

	return false
}
