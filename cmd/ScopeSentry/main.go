// Package ScopeSentry -----------------------------
// @file      : main.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/6 17:24
// -------------------------------------------
package main

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/runner"
)

func main() {
	config.SetUp()
	fmt.Println(config.ConfigDir)
	runnerOption := runner.Option{
		SubfinderEnabled:     true,
		SubdomainScanEnabled: true,
		PortScanEnabled:      true,
		DirScanEnabled:       true,
		CrawlerEnabled:       true,
	}
	Host := "rainy-autumn.top"
	runner.Process(Host, runnerOption)
	// setup the scan config (mirrors command line options)
	//fxConfig := scan.Config{
	//	DefaultTimeout: time.Duration(2) * time.Second,
	//	FastMode:       false,
	//	Verbose:        false,
	//	UDP:            false,
	//}
	//
	//// create a target list to scan
	//ip, _ := netip.ParseAddr("146.148.61.165")
	//target := plugins.Target{
	//	Address: netip.AddrPortFrom(ip, 443),
	//	Host:    "praetorian.com",
	//}
	//targets := make([]plugins.Target, 1)
	//targets = append(targets, target)
	//
	//// run the scan
	//results, err := scan.ScanTargets(targets, fxConfig)
	//if err != nil {
	//	log.Fatalf("error: %s\n", err)
	//}
	//
	//// process the results
	//for _, result := range results {
	//	fmt.Printf("%s:%d (%s/%s)\n", result.Host, result.Port, result.Transport, result.Protocol)
	//}
}
