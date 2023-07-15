package pproftoggle_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/saisrikark/pproftoggle"
	"github.com/saisrikark/pproftoggle/rules"
)

func TestServe(t *testing.T) {
	_, err := pproftoggle.NewToggler(
		pproftoggle.Config{
			PollInterval: time.Second * 3,
			HttpServer: &http.Server{
				Addr: ":8080",
			},
			Rules: []pproftoggle.Rule{
				rules.EnvVarRule{
					Key:   "abcd",
					Value: "efgh",
				},
			},
		})
	if err != nil {
		fmt.Println(err)
	}

	// set := true
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 2)
	// 		if set {
	// 			fmt.Println("setting")
	// 			os.Setenv("abcd", "efgh")
	// 			set = false
	// 		} else {
	// 			fmt.Println("unsetting")
	// 			os.Unsetenv("abcd")
	// 			set = true
	// 		}
	// 	}
	// }()

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	// defer cancel()

	// if err := toggler.Serve(ctx); err != nil {
	// 	fmt.Println(err)
	// }
}
