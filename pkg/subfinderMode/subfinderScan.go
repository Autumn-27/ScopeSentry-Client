// Package subfinderMode -----------------------------
// @file      : subfinderScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/9 22:36
// -------------------------------------------
package subfinderMode

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/subdomainMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"log"
	"path/filepath"
	"strings"
)

func SubfinderScan(Host string) []types.SubdomainResult {

	subfinderOpts := &runner.Options{
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		// ResultCallback: func(s *resolve.HostEntry) {
		// callback function executed after each unique subdomain is found
		// },
		ProviderConfig: filepath.Join(config.ConfigDir, "subfinderConfig.yaml"),
		// and other config related options
	}

	// disable timestamps in logs / configure logger
	log.SetFlags(0)

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] failed to create subfinder runner: ", err),
		}
		logMode.PrintLog(myLog)
		log.Fatalf("failed to create subfinder runner: %v", err)
	}

	output := &bytes.Buffer{}
	// To run subdomain enumeration on a single domain
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), Host, []io.Writer{output}); err != nil {

		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] failed to enumerate single domain:", err),
		}
		logMode.PrintLog(myLog)
		log.Fatalf("failed to enumerate single domain: %v", err)

	}

	// To run subdomain enumeration on a list of domains from file/reader
	// file, err := os.Open("domains.txt")
	// if err != nil {
	// 	log.Fatalf("failed to open domains file: %v", err)
	// }
	// defer file.Close()
	// if err = subfinder.EnumerateMultipleDomainsWithCtx(context.Background(), file, []io.Writer{output}); err != nil {
	// 	log.Fatalf("failed to enumerate subdomains from file: %v", err)
	// }

	// print the output
	log.Println(output.String())
	outputString := output.String()
	lines := []string{}
	if outputString != "" {

		scanner := bufio.NewScanner(strings.NewReader(outputString))

		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}

		if err := scanner.Err(); err != nil {

			myLog := logMode.CustomLog{
				Status: "Error",
				Msg:    fmt.Sprintf("[Err] scanner subfinder result:", err),
			}
			logMode.PrintLog(myLog)
		}

	}
	subfinderResult := subdomainMode.SubdomainScan(lines)

	return subfinderResult
}
