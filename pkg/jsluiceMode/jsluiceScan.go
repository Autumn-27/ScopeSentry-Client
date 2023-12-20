// Package jsluice -----------------------------
// @file      : jsluiceScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/16 16:26
// -------------------------------------------
package jsluiceMode

import (
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/BishopFox/jsluice"
	"regexp"
)

type SecretData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Secret struct {
	Kind     string            `json:"kind"`
	Data     SecretData        `json:"data"`
	Severity string            `json:"severity"`
	Context  map[string]string `json:"context"`
}

func JsluiceScan(url string, resp string, secrets config.Secrets) []types.SecretResults {

	analyzer := jsluice.NewAnalyzer([]byte(resp))

	analyzer.AddSecretMatcher(
		// The first value in the jsluice.SecretMatcher struct is a
		// tree-sitter query to run on the JavaScript source.
		jsluice.SecretMatcher{"(pair) @match", func(n *jsluice.Node) *jsluice.Secret {
			key := n.ChildByFieldName("key").DecodedString()
			value := n.ChildByFieldName("value").DecodedString()
			var id string
			for _, rule := range secrets.Rules {
				if rule.Enabled {
					pattern := rule.Pattern
					regexpPattern, err := regexp.Compile(pattern)
					if err != nil {
						fmt.Printf("Error compiling regex: %v\n", err)
						myLog := logMode.CustomLog{
							Status: "Error",
							Msg:    fmt.Sprintf("[Err] Error JsluiceScan compiling regex:", err),
						}
						logMode.PrintLog(myLog)
						fmt.Println("Error JsluiceScan compiling regex:", err)
						return nil
					}
					mathStr := ""
					if rule.Type == "value" {
						mathStr = value
					} else {
						mathStr = key
					}
					match := regexpPattern.MatchString(mathStr)
					if match {
						id = rule.ID
						return &jsluice.Secret{
							Kind: id,
							Data: map[string]string{
								"key":   key,
								"value": value,
							},
							Severity: jsluice.SeverityLow,
							Context:  n.Parent().AsMap(),
						}
					}
				}
			}
			return nil
		}},
	)
	secretsResult := []types.SecretResults{}
	for _, match := range analyzer.GetSecrets() {
		j, err := json.MarshalIndent(match, "", "  ")
		if err != nil {
			continue
		}
		var secret Secret
		if err := json.Unmarshal(j, &secret); err != nil {
			continue
		}
		secretsResultTemp := types.SecretResults{
			Url:   url, // Replace with your actual URL
			Kind:  secret.Kind,
			Key:   secret.Data.Key,
			Value: secret.Data.Value,
		}
		secretsResult = append(secretsResult, secretsResultTemp)
	}
	return secretsResult
}
