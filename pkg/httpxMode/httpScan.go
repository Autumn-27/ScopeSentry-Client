// Package httpx_mode -----------------------------
// @file      : httpx-scan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/7 21:29
// -------------------------------------------
package httpxMode

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/cloudflare/cfssl/log"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/runner"
)

func HttpxScan(Host string, resultCallback func(runner.Result)) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelFatal) // increase the verbosity (optional)

	options := runner.Options{
		Methods:         "GET",
		JSONOutput:      true,
		TLSProbe:        true,
		InputTargetHost: goflags.StringSlice{Host},
		//InputFile: "./targetDomains.txt", // path to file containing the target domains list
		OnResult: func(r runner.Result) {
			// handle error
			if r.Err != nil {
				httpxLog := logMode.CustomLog{
					Status: "Error",
					Msg:    fmt.Sprintf("[Err] %s: %s\n", r.Input, r.Err),
				}
				logMode.PrintLog(httpxLog)
				//fmt.Printf("[Err] %s: %s\n", r.Input, r.Err)
				return
			}
			//fmt.Printf("%s %s %d\n", r.Input, r.Host, r.StatusCode)
			resultCallback(r)
		},
	}

	if err := options.ValidateOptions(); err != nil {
		log.Fatal(err)
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()
}
