// Package urlScanMode -----------------------------
// @file      : urlScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/15 19:43
// -------------------------------------------
package urlScanMode

import (
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/jsluiceMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/util"
	"github.com/jaeles-project/gospider/core"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Option struct {
	Cookie  string
	Headers []string
}

func isMatchingFilter(fs []*regexp.Regexp, d []byte) bool {
	for _, r := range fs {
		if r.Match(d) {
			return true
		}
	}
	return false
}

func Run(op Option, siteList []string) ([]types.UrlResult, []types.SecretResults) {

	core.Logger.SetLevel(logrus.PanicLevel)

	// Check again to make sure at least one site in slice
	if len(siteList) == 0 {
		core.Logger.Info("No site in list. Please check your site input again")
		os.Exit(1)
	}

	threads := 1
	sitemap := true
	linkfinder := true
	robots := true
	otherSource := false
	includeSubs := true
	includeOtherSourceResult := true

	seenUrls := make(map[string]struct{})
	urlInfos := []types.UrlResult{}
	var DisallowedURLFilters []*regexp.Regexp
	disallowedRegex := `(?i)\.(png|apng|bmp|gif|ico|cur|jpg|jpeg|jfif|pjp|pjpeg|svg|tif|tiff|webp|xbm|3gp|aac|flac|mpg|mpeg|mp3|mp4|m4a|m4v|m4p|oga|ogg|ogv|mov|wav|webm|eot|woff|woff2|ttf|otf|css)(?:\?|#|$)`
	DisallowedURLFilters = append(DisallowedURLFilters, regexp.MustCompile(disallowedRegex))

	urlScanResultHandler := func(msg string) {
		urlInfo := types.UrlResult{}
		err := json.Unmarshal([]byte(msg), &urlInfo)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		urlInfo.Time = util.GetTimeNow()
		var url string
		if strings.HasPrefix(urlInfo.Output, "http") {
			url = urlInfo.Output
		} else {
			url = urlInfo.Input + urlInfo.Output
		}
		if !isMatchingFilter(DisallowedURLFilters, []byte(url)) {
			if _, seen := seenUrls[url]; !seen {
				seenUrls[url] = struct{}{}
				urlInfos = append(urlInfos, urlInfo)
			}
		}
	}
	secrets, err := config.GetSecretsRules()
	if err != nil {
		fmt.Println("Error get secretsRules:", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error get secretsRules:", err),
		}
		logMode.PrintLog(myLog)
	}
	secretsResult := []types.SecretResults{}
	respBodyHandler := func(url string, msg string) {

		if !isMatchingFilter(DisallowedURLFilters, []byte(url)) {
			secretsResultTemp := jsluiceMode.JsluiceScan(url, msg, secrets)
			secretsResult = append(secretsResult, secretsResultTemp...)

		}
	}

	var wg sync.WaitGroup
	inputChan := make(chan string, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rawSite := range inputChan {
				site, err := url.Parse(rawSite)
				if err != nil {
					logrus.Errorf("Failed to parse %s: %s", rawSite, err)
					continue
				}

				var siteWg sync.WaitGroup
				crawler := core.NewCrawler(site, op.Cookie, op.Headers, urlScanResultHandler, respBodyHandler)
				siteWg.Add(1)
				go func() {
					defer siteWg.Done()
					crawler.Start(linkfinder)
				}()

				// Brute force Sitemap path
				if sitemap {
					siteWg.Add(1)
					go core.ParseSiteMap(site, crawler, crawler.C, &siteWg)
				}

				// Find Robots.txt
				if robots {
					siteWg.Add(1)
					go core.ParseRobots(site, crawler, crawler.C, &siteWg)
				}

				if otherSource {
					siteWg.Add(1)
					go func() {
						defer siteWg.Done()
						urls := core.OtherSources(site.Hostname(), includeSubs)
						for _, url := range urls {
							url = strings.TrimSpace(url)
							if len(url) == 0 {
								continue
							}

							outputFormat := fmt.Sprintf("[other-sources] - %s", url)
							if includeOtherSourceResult {
								if crawler.JsonOutput {
									sout := core.SpiderOutput{
										Input:      crawler.Input,
										Source:     "other-sources",
										OutputType: "url",
										Output:     url,
									}
									if data, err := jsoniter.MarshalToString(sout); err == nil {
										outputFormat = data
									}
								} else if crawler.Quiet {
									outputFormat = url
								}
								fmt.Println(outputFormat)

								if crawler.Output != nil {
									crawler.Output.WriteToFile(outputFormat)
								}
							}

							_ = crawler.C.Visit(url)
						}
					}()
				}
				siteWg.Wait()
				crawler.C.Wait()
				crawler.LinkFinderCollector.Wait()
			}
		}()
	}

	for _, site := range siteList {
		inputChan <- site
	}
	close(inputChan)
	wg.Wait()
	//core.Logger.Info("Done.")
	return urlInfos, secretsResult

}
