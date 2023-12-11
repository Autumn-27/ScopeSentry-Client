// Package subdomainMode -----------------------------
// @file      : subdomainScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/9 23:22
// -------------------------------------------
package subdomainMode

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/boy-hack/ksubdomain/core/gologger"
	"github.com/boy-hack/ksubdomain/core/options"
	"github.com/boy-hack/ksubdomain/runner"
	"github.com/boy-hack/ksubdomain/runner/outputter"
	"github.com/boy-hack/ksubdomain/runner/outputter/output"
	"github.com/boy-hack/ksubdomain/runner/processbar"
)

func subdomainScan(Host string) {

	var scopeSentryConfig config.ScopeSentryConfig
	scopeSentryConfig = config.ParseConfig()
	var threadNumber = scopeSentryConfig.Subdomain.ThreadNumber
	fmt.Println(threadNumber)

}

func main() {
	process := processbar.ScreenProcess{}
	screenPrinter, _ := output.NewScreenOutput(false)

	domains := []string{"www.hacking8.com", "x.hacking8.com"}
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
		DomainTotal: 2,
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
}
