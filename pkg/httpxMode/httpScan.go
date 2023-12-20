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
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/util"
	"github.com/cloudflare/cfssl/log"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/runner"
)

func HttpxScan(Host []string, resultCallback func(r types.AssertHttp)) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelFatal) // increase the verbosity (optional)

	options := runner.Options{
		Methods:                   "GET",
		JSONOutput:                true,
		TLSProbe:                  true,
		InputTargetHost:           Host,
		Favicon:                   true,
		ExtractTitle:              true,
		TechDetect:                true,
		OutputWebSocket:           true,
		OutputIP:                  true,
		OutputCName:               false,
		ResponseHeadersInStdout:   true,
		ResponseInStdout:          true,
		Base64ResponseInStdout:    true,
		Jarm:                      true,
		OutputCDN:                 true,
		Location:                  true,
		MaxResponseBodySizeToRead: 100000,
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
			ah := httpxResultToAssertHttp(r)
			//fmt.Printf("%s %s %d\n", r.Input, r.Host, r.StatusCode)
			resultCallback(ah)
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

func httpxResultToAssertHttp(r runner.Result) types.AssertHttp {
	var ah = types.AssertHttp{
		Timestamp:    util.GetTimeNow(),
		TLSData:      r.TLSData, // You may need to set an appropriate default value based on the actual type.
		Hashes:       r.Hashes,
		CDNName:      r.CDNName,
		Port:         r.Port,
		URL:          r.URL,
		Location:     r.Location,
		Title:        r.Title,
		Type:         r.Scheme,
		Error:        r.Error,
		ResponseBody: r.ResponseBody,
		Host:         r.Host,
		FavIconMMH3:  r.FavIconMMH3,
		FaviconPath:  r.FaviconPath,
		RawHeaders:   r.RawHeaders,
		Jarm:         r.Jarm,
		Technologies: r.Technologies, // You may need to set an appropriate default value based on the actual type.
		StatusCode:   r.StatusCode,   // You may need to set an appropriate default value.
		Webcheck:     false,
	}
	return ah

}
