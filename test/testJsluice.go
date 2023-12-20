// Package main -----------------------------
// @file      : testJsluice.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/16 16:26
// -------------------------------------------
package main

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/jsluiceMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"os"
	"path/filepath"
)

func main() {
	executableDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("[Err] Failed to retrieve the directory of the executable file:", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Failed to retrieve the directory of the executable file:", err),
		}
		logMode.PrintLog(myLog)
		return
	}
	config.ConfigDir = filepath.Join(executableDir, "config")
	config.SetUp()
	secrets, err := config.GetSecretsRules()
	if err != nil {
		fmt.Println("Error get secretsRules:", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error get secretsRules:", err),
		}
		logMode.PrintLog(myLog)
	}

	jsluiceMode.JsluiceScan("http://127.0.0.1", "{\"key\":\"aaa\"}", secrets)
}
