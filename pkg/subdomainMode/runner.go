// Package subdomainMode -----------------------------
// @file      : runner.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/11 19:15
// -------------------------------------------
package subdomainMode

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/util"
	"path"
)

func SubDomainRunner(Host string) []types.SubdomainResult {
	dicFilePath := path.Join(config.ConfigDir, "domainDic")

	subdomainDict, err := util.ReadFileLiness(dicFilePath)
	if err != nil {
		fmt.Printf("[Err] open file :%s error:%s", dicFilePath, err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] open file :%s error:%s", dicFilePath, err),
		}
		logMode.PrintLog(myLog)
	}
	dicDomainList := []string{}
	for _, value := range subdomainDict {
		dicDomainList = append(dicDomainList, value+"."+Host)
	}

	subDoaminResult := SubdomainScan(dicDomainList)

	return subDoaminResult

}
