// Package runner -----------------------------
// @file      : runner.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/9 20:20
// -------------------------------------------
package runner

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/httpxMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/subfinderMode"
	httpxRunner "github.com/projectdiscovery/httpx/runner"
	"net/url"
	"strings"
	"sync"
)

type Option struct {
	SubfinderEnabled     bool
	SubdomainScanEnabled bool
	PortScanEnabled      bool
	DirScanEnabled       bool
	CrawlerEnabled       bool
}

func Process(Host string, op Option) {
	normalizedHttp := ""
	if !strings.HasPrefix(Host, "http://") && !strings.HasPrefix(Host, "https://") {
		// 如果不以 "http://" 或 "https://" 开头，则在其前面添加 "http://"
		normalizedHttp = "http://" + Host
	} else {
		normalizedHttp = Host
	}
	parsedURL, err := url.Parse(normalizedHttp)
	if err != nil {
		httpxLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] %s: %s\n", Host, "parse url error"),
		}
		logMode.PrintLog(httpxLog)
		return
	}
	hostParts := strings.Split(parsedURL.Host, ":")
	hostWithoutPort := hostParts[0]
	port := parsedURL.Port()
	fmt.Println(hostWithoutPort)
	fmt.Println(port)
	subfinderDomsString := ""
	if port != "" {
		subfinderDomsString = Host + "\n"
		subfinderDomsString += hostWithoutPort + "\n"
	} else {
		subfinderDomsString = Host + "\n"
	}

	if op.SubfinderEnabled != false {
		subfinderDomsString += subfinderMode.SubfinderScan(hostWithoutPort)
	}

	if op.SubdomainScanEnabled != false {

	}

	var httpxResults []httpxRunner.Result
	var httpxResultsMutex sync.Mutex
	httpxResultsHandler := func(r httpxRunner.Result) {
		fmt.Printf("Result in process: %s %s %d\n", r.Input, r.Host, r.StatusCode)

		httpxResultsMutex.Lock()
		httpxResults = append(httpxResults, r)
		httpxResultsMutex.Unlock()
	}

	httpxMode.HttpxScan("rainy-autumn.top", httpxResultsHandler)
	fmt.Println("所有结果:", httpxResults)
	// 待解析的URL字符串
	//rawURL := "https://example.com"
	//
	//// 解析URL
	//parsedURL, err := url.Parse(rawURL)
	//if err != nil {
	//	fmt.Println("Error parsing URL:", err)
	//	return
	//}
	//// 打印URL的各个部分
	//fmt.Println(parsedURL.Port())
	//fmt.Println("Scheme:", parsedURL.Scheme)
	//fmt.Println("Host:", parsedURL.Host)
	//fmt.Println("Path:", parsedURL.Path)
	//fmt.Println("RawQuery:", parsedURL.RawQuery)
	//
	//// 解析查询参数
	//queryParams, err := url.ParseQuery(parsedURL.RawQuery)
	//if err != nil {
	//	fmt.Println("Error parsing query parameters:", err)
	//	return
	//}
	//
	//// 打印查询参数
	//fmt.Println("Name:", queryParams.Get("name"))
	//fmt.Println("Age:", queryParams.Get("age"))
}
