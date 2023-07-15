package pproftoggle_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/saisrikark/pproftoggle"
)

func TestListen(t *testing.T) {
	t.Run("HttpServerNotSpecified", func(t *testing.T) {
		_, err := pproftoggle.NewServer(pproftoggle.ServerConfig{})
		if err == nil {
			t.Errorf("%s", err.Error())
		}
	})

	t.Run("HttpServerShouldNotStartAtPortAlreadyBeingUsed", func(t *testing.T) {
		var lc net.Listener
		var wg sync.WaitGroup
		var err error

		wg.Add(1)
		go func() {
			lc, err = net.Listen("tcp", ":"+port)
			if err != nil {
				t.Errorf("%s", err.Error())
			}
			wg.Done()
		}()
		wg.Wait()
		if lc != nil {
			defer lc.Close()
		}

		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + port},
		})
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		if err = ppfs.Listen(ctx); err == nil {
			t.Errorf("%s", err.Error())
		}
	})

	t.Run("FetchFromEndpoint", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + port},
		})
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		go func() {
			if err = ppfs.Listen(ctx); err != nil && err != http.ErrServerClosed {
				t.Errorf("%s", err.Error())
			}
		}()

		for i := 0; i < 5; i++ {
			if isPortUsed(port) {
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
		if !isPortUsed(port) {
			t.Errorf("port:[%s] not being used", port)
		}

		resp := ""
		if err := requests.URL(endpoint).BaseURL(baseURL).ToString(&resp).Fetch(context.Background()); err != nil {
			t.Errorf("%s", err.Error())
		}

		cancel()
	})
}

func TestShutdown(t *testing.T) {
	t.Run("ServerMustShutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
			HttpServer: &http.Server{Addr: ":" + port},
		})
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		go func() {
			if err := ppfs.Listen(ctx); err != nil && err != http.ErrServerClosed {
				t.Errorf("%s", err.Error())
			}
		}()

		for i := 0; i < 5; i++ {
			if isPortUsed(port) {
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
		if !isPortUsed(port) {
			t.Errorf("port:[%s] not being used", port)
		}

		resp := ""
		if err := requests.URL(endpoint).BaseURL(baseURL).ToString(&resp).Fetch(context.Background()); err != nil {
			t.Errorf("%s", err)
		}

		ppfs.Shutdown(context.Background())

		if err := requests.URL(endpoint).BaseURL(baseURL).ToString(&resp).Fetch(context.Background()); err == nil {
			t.Errorf("unexpected response\n[%s]", resp[0:30])
		}
	})
}

func TestIsRunning(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
		HttpServer: &http.Server{Addr: ":" + port},
	})
	if err != nil {
		t.Errorf("%s", err)
	}

	go func() {
		if err := ppfs.Listen(ctx); err != nil && err != http.ErrServerClosed {
			t.Errorf("%s", err.Error())
		}
	}()

	for i := 0; i < 5; i++ {
		if isPortUsed(port) {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	if !isPortUsed(port) {
		t.Errorf("port:[%s] not being used", port)
	}

	resp := ""
	if err := requests.URL(endpoint).BaseURL(baseURL).ToString(&resp).Fetch(context.Background()); err != nil {
		t.Errorf("%s", err)
	} else {
		status := ppfs.IsRunning()
		if !status {
			t.Errorf("received response but IsRunning is showing [%v]", status)
		}
	}

	if err := ppfs.Shutdown(context.Background()); err != nil {
		t.Errorf("%s", err.Error())
	}

	if err := requests.URL(endpoint).BaseURL(baseURL).ToString(&resp).Fetch(context.Background()); err == nil {
		t.Errorf("received response from endpoint even after shutting down resp\n[%s]", resp[0:10])
	}

	if isPortUsed(port) {
		t.Errorf("port [%s] is still being used after shuttiing down", port)
	}
}

func BenchmarkToggle(b *testing.B) {
	b.SetParallelism(1)

	cleanUp := func() {
		runtime.GC()
	}
	defer b.Cleanup(cleanUp)

	ppfs, err := pproftoggle.NewServer(pproftoggle.ServerConfig{
		HttpServer: &http.Server{Addr: ":" + port}})
	if err != nil {
		b.Errorf("%s", err.Error())
	}

	doTest := func(b *testing.B) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			ppfs.Listen(ctx)
		}()

		for {
			if ppfs.IsRunning() {
				break
			}
		}

		if err := ppfs.Shutdown(ctx); err != nil {
			b.Skipf(err.Error())
		}
		wg.Wait()
	}

	for i := 0; i < 200; i++ {
		b.Run(fmt.Sprintf("%d", i), doTest)
	}
}

func isPortUsed(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(localhost, port), time.Millisecond*10)
	if err != nil {
		return false
	}

	if conn != nil {
		defer conn.Close()
		return true
	}

	return false
}
