package pproftoggle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saisrikark/pproftoggle"
)

func TestServe(t *testing.T) {
	toggler, err := pproftoggle.NewToggler(pproftoggle.Config{
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	if err := toggler.Serve(ctx); err != nil {
		fmt.Println(err)
	}
}
