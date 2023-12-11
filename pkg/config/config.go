// Package config -----------------------------
// @file      : config.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/9 21:57
// -------------------------------------------
package config

import (
	"bufio"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ScopeSentryConfig struct {
	Subdomain struct {
		ThreadNumber int `yaml:"ThreadNumber"`
	}
}

var ConfigDir string

func SetUp() {
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
	ConfigDir = filepath.Join(executableDir, "config")
	if err := os.MkdirAll(ConfigDir, os.ModePerm); err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to create config folder:", err),
		}
		logMode.PrintLog(myLog)
		return
	}
	subfinderConfigPath := filepath.Join(ConfigDir, "subfinderConfig.yaml")
	if _, err := os.Stat(subfinderConfigPath); os.IsNotExist(err) {
		content := SubfinderDefaultConfig
		if err := ioutil.WriteFile(subfinderConfigPath, content, os.ModePerm); err != nil {

			fmt.Println("Failed to create subfinderConfig.yaml file:", err)
			myLog := logMode.CustomLog{
				Status: "Error",
				Msg:    fmt.Sprintf("[Err] Failed to create subfinderConfig.yaml file:", err),
			}
			logMode.PrintLog(myLog)
			return
		}
	}

	fmt.Println("Configuration folder:", ConfigDir)
	fmt.Println("Subfinder config file path:", subfinderConfigPath)

	domainDicPath := filepath.Join(ConfigDir, "domainDic")
	if _, err := os.Stat(domainDicPath); os.IsNotExist(err) {
		content := subdomainDicDefault
		if err := ioutil.WriteFile(domainDicPath, content, os.ModePerm); err != nil {

			fmt.Println("Failed to create domainDic file:", err)
			myLog := logMode.CustomLog{
				Status: "Error",
				Msg:    fmt.Sprintf("[Err] Failed to create domainDic file:", err),
			}
			logMode.PrintLog(myLog)
			return
		}
	}

	scopeSentryConfigPath := filepath.Join(ConfigDir, "scopeSentryConfig.ymal")
	if _, err := os.Stat(scopeSentryConfigPath); os.IsNotExist(err) {
		content := scopSentryDefault
		if err := ioutil.WriteFile(scopeSentryConfigPath, content, os.ModePerm); err != nil {

			fmt.Println("Failed to create domainDic file:", err)
			myLog := logMode.CustomLog{
				Status: "Error",
				Msg:    fmt.Sprintf("[Err] Failed to create domainDic file:", err),
			}
			logMode.PrintLog(myLog)
			return
		}
	}
}

func GetDomainDic() []string {
	domainDidPath := filepath.Join(ConfigDir, "domainDic")
	// Open the file
	file, err := os.Open(domainDidPath)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Failed to create domainDic file:", err),
		}
		logMode.PrintLog(myLog)
		fmt.Println("Unable to open the domainDic file:", err)
		return nil
	}
	defer file.Close()

	// Use bufio.Scanner to read the file content line by line
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check for any errors that occurred during file reading
	if err := scanner.Err(); err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error reading the domainDic file:", err),
		}
		logMode.PrintLog(myLog)
		fmt.Println("Error reading the domainDic file:", err)
		return nil
	}

	return lines
}

func ParseConfig() ScopeSentryConfig {
	var defaultConfig = ScopeSentryConfig{
		Subdomain: struct {
			ThreadNumber int `yaml:"ThreadNumber"`
		}{
			ThreadNumber: 10,
		},
	}

	scopeSentryConfigPath := filepath.Join(ConfigDir, "scopeSentryConfig.ymal")
	yamlFile, err := ioutil.ReadFile(scopeSentryConfigPath)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error reading scopeSentryConfig.ymal file:", err),
		}
		logMode.PrintLog(myLog)
		return defaultConfig
	}
	var config ScopeSentryConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error unmarshaling YAML:", err),
		}
		logMode.PrintLog(myLog)
		fmt.Println("Error unmarshaling YAML:", err)
		return defaultConfig
	}

	return config
}
