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
	"github.com/Autumn-27/ScopeSentry-Client/pkg/portScanMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/subdomainMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/subfinderMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/urlScanMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/util"
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
	Ports                string
	WaybackurlEnabled    bool
	UrlScan              bool
	Cookie               string
	Header               []string
}

func Process(Host string, op Option) {
	normalizedHttp := ""
	if !strings.HasPrefix(Host, "http://") && !strings.HasPrefix(Host, "https://") {
		normalizedHttp = "http://" + Host
	} else {
		normalizedHttp = Host
	}
	parsedURL, err := url.Parse(normalizedHttp)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] %s: %s\n", Host, "parse url error"),
		}
		logMode.PrintLog(myLog)
		return
	}
	hostParts := strings.Split(parsedURL.Host, ":")
	hostWithoutPort := hostParts[0]
	port := parsedURL.Port()
	SubDomainResults := []types.SubdomainResult{}
	domainDnsResult := subdomainMode.SubdomainScan([]string{hostWithoutPort})
	SubDomainResults = append(SubDomainResults, domainDnsResult...)
	if op.SubfinderEnabled {
		subfinderResult := subfinderMode.SubfinderScan(hostWithoutPort)
		SubDomainResults = append(SubDomainResults, subfinderResult...)
	}

	if op.SubdomainScanEnabled {
		// 判断是否泛解析，跳过泛解析
		if !util.IsWildCard(hostWithoutPort) {
			subDomainResult := subdomainMode.SubDomainRunner(hostWithoutPort)
			SubDomainResults = append(SubDomainResults, subDomainResult...)
		}
	}
	domainList := []string{}
	if port != "" {
		domainList = append(domainList, Host)
	}
	uniqueSubDomainResults := []types.SubdomainResult{}
	seenHosts := make(map[string]struct{})
	for _, result := range SubDomainResults {
		if _, seen := seenHosts[result.Host]; seen {
			continue
		}
		seenHosts[result.Host] = struct{}{}
		domainList = append(domainList, result.Host)
		uniqueSubDomainResults = append(uniqueSubDomainResults, result)
	}

	// 子域名接管 todo

	var httpxResults []types.AssertHttp
	var httpxResultsMutex sync.Mutex
	httpxResultsHandler := func(r types.AssertHttp) {
		httpxResultsMutex.Lock()
		httpxResults = append(httpxResults, r)
		httpxResultsMutex.Unlock()
	}
	httpxMode.HttpxScan(domainList, httpxResultsHandler)

	assertOthers := []types.AssertOther{}
	if op.PortScanEnabled {
		for _, uniqueSubDomainResult := range uniqueSubDomainResults {
			assertHttpTemp, assertOtherTemp := portScanMode.PortScan(uniqueSubDomainResult.Host, op.Ports)
			httpxResults = append(httpxResults, assertHttpTemp...)
			assertOthers = append(assertOthers, assertOtherTemp...)
		}
	}

	//目录扫描 todo
	//缓存污染 todo

	var urlResults = []types.UrlResult{}
	var secretsResult = []types.SecretResults{}
	//url扫描、js信息泄露
	if op.UrlScan {
		domainUrlScanList := []string{}
		for _, httpxResult := range httpxResults {
			domainUrlScanList = append(domainUrlScanList, httpxResult.URL)
		}
		urlScanOption := urlScanMode.Option{
			Cookie:  op.Cookie,
			Headers: op.Header,
		}
		urlResults, secretsResult = urlScanMode.Run(urlScanOption, domainUrlScanList)
	}
	fmt.Println(urlResults)
	fmt.Println(secretsResult)

	uniqueurResults := []types.UrlResult{}
	seenUrls := make(map[string]struct{})
	urlList := []string{}
	for _, result := range urlResults {
		urlTemp := ""
		isHTTP := strings.HasPrefix(result.Output, "http")
		if isHTTP {
			urlTemp = result.Output
		} else {
			urlTemp = result.Input + result.Source
		}
		if _, seen := seenUrls[urlTemp]; seen {
			continue
		}
		seenUrls[urlTemp] = struct{}{}
		urlList = append(urlList, urlTemp)
		uniqueurResults = append(uniqueurResults, result)
	}
	if op.CrawlerEnabled {

	}

}
