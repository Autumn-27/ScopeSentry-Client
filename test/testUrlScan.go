// Package main -----------------------------
// @file      : testUrlScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/15 20:37
// -------------------------------------------
package main

import "github.com/Autumn-27/ScopeSentry-Client/pkg/urlScanMode"

func main() {
	op := urlScanMode.Option{
		Headers: []string{},
		Cookie:  "",
	}
	domainList := []string{"http://127.0.0.1:666/"}
	urlScanMode.Run(op, domainList)
}
