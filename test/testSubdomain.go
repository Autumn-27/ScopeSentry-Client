// Package test -----------------------------
// @file      : testSubdomain.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/10 20:54
// -------------------------------------------
package main

import (
	"context"
	"fmt"
	"github.com/boy-hack/ksubdomain/core/gologger"
	"github.com/boy-hack/ksubdomain/core/options"
	"github.com/boy-hack/ksubdomain/runner"
	"github.com/boy-hack/ksubdomain/runner/outputter"
	"github.com/boy-hack/ksubdomain/runner/outputter/output"
	"github.com/boy-hack/ksubdomain/runner/processbar"
)

func main() {
	process := processbar.ScreenProcess{}

	resultCallback := func(Host string, IP string) {
		// Do something with the msg in the context of the main function
		fmt.Println("Received message in main:", Host, IP)
	}

	screenPrinter, _ := output.NewScreenOutput(false, resultCallback)

	domains := []string{"rainy-autumn.top", "a.rainy-autumn.top"}
	domainChanel := make(chan string)
	go func() {
		for _, d := range domains {
			domainChanel <- d
		}
		close(domainChanel)
	}()
	opt := &options.Options{
		Rate:        options.Band2Rate("1m"),
		Domain:      domainChanel,
		DomainTotal: 2,
		Resolvers:   options.GetResolvers(""),
		Silent:      false,
		TimeOut:     10,
		Retry:       3,
		Method:      runner.VerifyType,
		DnsType:     "a",
		Writer: []outputter.Output{
			screenPrinter,
		},
		ProcessBar: &process,
		EtherInfo:  options.GetDeviceConfig(),
	}
	opt.Check()
	r, err := runner.New(opt)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	ctx := context.Background()
	r.RunEnumeration(ctx)
	r.Close()
}
