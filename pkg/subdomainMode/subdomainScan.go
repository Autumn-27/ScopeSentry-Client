// Package subdomainMode -----------------------------
// @file      : subdomainScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/9 23:22
// -------------------------------------------
package subdomainMode

import (
	"context"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/util"
	"github.com/boy-hack/ksubdomain/core/gologger"
	"github.com/boy-hack/ksubdomain/core/options"
	"github.com/boy-hack/ksubdomain/runner"
	"github.com/boy-hack/ksubdomain/runner/outputter"
	"github.com/boy-hack/ksubdomain/runner/outputter/output"
	"github.com/boy-hack/ksubdomain/runner/processbar"
	"strings"
)

func SubdomainScan(Host []string) []types.SubdomainResult {
	//var scopeSentryConfig config.ScopeSentryConfig
	//scopeSentryConfig := config.ParseConfig()
	//var threadNumber = scopeSentryConfig.Subdomain.ThreadNumber
	subDomainResult := []types.SubdomainResult{}
	process := processbar.ScreenProcess{}
	resultCallback := func(Domains []string) {
		// Do something with the msg in the context of the main function
		//fmt.Println("Received message in main:", Domains)

		_domain := types.SubdomainResult{}
		_domain.Host = Domains[0]
		_domain.Type = "A"
		for i := 1; i < len(Domains); i++ {
			containsSpace := strings.Contains(Domains[i], " ")
			if containsSpace {
				result := strings.SplitN(Domains[i], " ", 2)
				_domain.Type = result[0]
				_domain.Value = append(_domain.Value, result[1])
			} else {
				//_domain.Value = append(_domain.Value, Domains[i])
				_domain.IP = append(_domain.IP, Domains[i])
			}
		}
		time := util.GetTimeNow()
		_domain.Time = time
		subDomainResult = append(subDomainResult, _domain)
	}

	screenPrinter, _ := output.NewScreenOutput(false, resultCallback)

	domains := Host
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
		DomainTotal: len(domains),
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
	return subDomainResult

}
