// Package portScanMode -----------------------------
// @file      : portScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/13 18:13
// -------------------------------------------
package portScanMode

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/httpxMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/praetorian-inc/fingerprintx/pkg/plugins"
	"github.com/praetorian-inc/fingerprintx/pkg/scan"
	"net/netip"
	"strconv"
	"sync"
	"time"
)

func PortScan(Domains string, Ports string) ([]types.AssertHttp, []types.AssertOther) {
	fmt.Printf("PortScan start")
	var PortAlives []types.PortAlive
	CallBack := func(alive []types.PortAlive) {
		PortAlives = alive
	}
	var err = NaabuScan([]string{Domains}, Ports, CallBack)
	if err != nil {
		fmt.Printf(err.Error())
	}
	var targets []plugins.Target
	for _, value := range PortAlives {
		ip, _ := netip.ParseAddr(value.IP)
		target := plugins.Target{
			Address: netip.AddrPortFrom(ip, uint16(value.Port)),
			Host:    value.Host,
		}
		targets = append(targets, target)
	}

	// setup the scan config (mirrors command line options)
	fxConfig := scan.Config{
		DefaultTimeout: time.Duration(2) * time.Second,
		FastMode:       false,
		Verbose:        false,
		UDP:            false,
	}

	// run the scan
	results, err := scan.ScanTargets(targets, fxConfig)
	if err != nil {
		fmt.Println("error: %s\n", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] PortScan error: %s\n", err),
		}
		logMode.PrintLog(myLog)
	}

	var httpxResults []types.AssertHttp
	var httpxResultsMutex sync.Mutex
	httpxResultsHandler := func(r types.AssertHttp) {
		//fmt.Printf("Result in process: %s %s %d\n", r.Host, r.StatusCode)
		httpxResultsMutex.Lock()
		httpxResults = append(httpxResults, r)
		httpxResultsMutex.Unlock()
	}
	urlList := []string{}
	assertOthers := []types.AssertOther{}
	// process the results
	for _, result := range results {
		if result.Protocol == "http" || result.Protocol == "https" {
			portStr := strconv.Itoa(result.Port)
			url := result.Protocol + "://" + result.Host + ":" + portStr
			urlList = append(urlList, url)
			httpxMode.HttpxScan(urlList, httpxResultsHandler)
		} else {
			assertedOther := types.AssertOther{
				Host:      result.Host,
				IP:        result.IP,
				Port:      result.Port,
				Protocol:  result.Protocol,
				TLS:       result.TLS,
				Transport: result.Transport,
				Version:   result.Version,
				Raw:       result.Raw,
				// Add additional fields specific to AssertOther if needed
			}
			assertOthers = append(assertOthers, assertedOther)
		}

	}
	assertOthersTemp := []types.AssertOther{}
	for _, portAlive := range PortAlives {
		found := false
		for _, assertOther := range assertOthers {
			if portAlive.Host == assertOther.Host && portAlive.Port == assertOther.Port {
				found = true
				break
			}
		}
		if found == false {
			assertedOther := types.AssertOther{
				Host:      portAlive.Host,
				IP:        portAlive.IP,
				Port:      portAlive.Port,
				Protocol:  "",
				TLS:       false,
				Transport: "",
				Version:   "",
			}
			assertOthersTemp = append(assertOthersTemp, assertedOther)
		}
	}
	assertOthers = append(assertOthers, assertOthersTemp...)

	return httpxResults, assertOthers

}
