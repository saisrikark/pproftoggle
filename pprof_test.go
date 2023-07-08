package pproftoggle

import "testing"

func TestNewpprofServer(t *testing.T) {
	ppfs, _ := newpprofServer(pprofServerConfig{
		Address:    ":8080",
		UserPrefix: "",
	})
	ppfs.start()
}

func TestStart(t *testing.T) {

}

func TestStop(t *testing.T) {

}
