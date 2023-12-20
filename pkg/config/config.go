// Package config -----------------------------
// @file      : config.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/9 21:57
// -------------------------------------------
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type ScopeSentryConfig struct {
	Subdomain struct {
		ThreadNumber int `yaml:"ThreadNumber"`
	} `yaml:"Subdomain"`
	Time struct {
		TimeZoneName string `yaml:"TimeZoneName"`
	} `yaml:"Time"`
}

type Secrets struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
	Pattern string `json:"pattern"`
}

var ConfigDir string
var ExtPath string
var CrawlerPath string
var CrawlerExecPath string

func SetUp() bool {
	executableDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("[Err] Failed to retrieve the directory of the executable file:", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Failed to retrieve the directory of the executable file:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	ConfigDir = filepath.Join(executableDir, "config")
	if err := os.MkdirAll(ConfigDir, os.ModePerm); err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to create config folder:", err),
		}
		logMode.PrintLog(myLog)
		return false
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
			return false
		}
	}

	//fmt.Println("Configuration folder:", ConfigDir)
	//fmt.Println("Subfinder config file path:", subfinderConfigPath)

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
			return false
		}
	}

	scopeSentryConfigPath := filepath.Join(ConfigDir, "scopeSentryConfig.yaml")
	if _, err := os.Stat(scopeSentryConfigPath); os.IsNotExist(err) {
		content := scopSentryDefault
		if err := ioutil.WriteFile(scopeSentryConfigPath, content, os.ModePerm); err != nil {

			fmt.Println("Failed to create domainDic file:", err)
			myLog := logMode.CustomLog{
				Status: "Error",
				Msg:    fmt.Sprintf("[Err] Failed to create domainDic file:", err),
			}
			logMode.PrintLog(myLog)
			return false
		}
	}

	SecretsRulesPath := filepath.Join(ConfigDir, "SecretsRules.yaml")
	if _, err := os.Stat(SecretsRulesPath); os.IsNotExist(err) {
		flag := generateScret(SecretsRulesPath)
		if !flag {
			return false
		}
	}
	flag := checkCrawler()
	if !flag {
		return false
	}

	return true
}

func checkCrawler() bool {
	executableDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("[Err] Failed to retrieve the directory of the executable file:", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Failed to retrieve the directory of the executable file:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	ExtPath = filepath.Join(executableDir, "ext")
	if err := os.MkdirAll(ExtPath, os.ModePerm); err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to create ext folder:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	radPath := filepath.Join(ExtPath, "rad")
	if err := os.MkdirAll(radPath, os.ModePerm); err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to create radPath folder:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	targetPath := filepath.Join(radPath, "target")
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to create targetPath folder:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	resultPath := filepath.Join(radPath, "result")
	if err := os.MkdirAll(resultPath, os.ModePerm); err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to create resultPath folder:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}

	osType := runtime.GOOS
	// 判断操作系统类型
	var path string
	switch osType {
	case "windows":
		path = "rad.exe"
	case "linux":
		path = "rad"
	}
	CrawlerPath = filepath.Join(ExtPath, "rad")
	CrawlerExecPath = filepath.Join(CrawlerPath, path)
	if _, err := os.Stat(CrawlerExecPath); os.IsNotExist(err) {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] The crawler tool rad does not exist:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	return true
}

func generateScret(SecretsRulesPath string) bool {
	// 解析 JSON 字符串到结构体
	var secrets Secrets
	err := json.Unmarshal([]byte(SecretsRulesJson), &secrets)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error parsing SecretsRulesJson:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}

	// 将结构体转换为 YAML 格式
	yamlData, err := yaml.Marshal(&secrets)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error converting to SecretsRules YAML:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}

	// 写入 YAML 文件
	err = ioutil.WriteFile(SecretsRulesPath, yamlData, 0644)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error writing to SecretsRules YAML:", err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	return true
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
		Time: struct {
			TimeZoneName string `yaml:"TimeZoneName"`
		}{
			TimeZoneName: "Asia/Shanghai",
		},
	}

	scopeSentryConfigPath := filepath.Join(ConfigDir, "scopeSentryConfig.yaml")
	yamlFile, err := ioutil.ReadFile(scopeSentryConfigPath)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error reading scopeSentryConfig.yaml file:", err),
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
		return defaultConfig
	}

	return config
}

func GetSecretsRules() (Secrets, error) {

	var secrets Secrets
	// 读取YAML文件内容
	SecretsRulesPath := filepath.Join(ConfigDir, "SecretsRules.yaml")
	var yamlFile, err = ioutil.ReadFile(SecretsRulesPath)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error reading SecretsRules.yaml file:", err),
		}
		logMode.PrintLog(myLog)
		return secrets, err
	}

	// 创建一个Config对象，用于存储解析后的配置

	// 使用yaml库解析YAML内容到Config对象
	err = yaml.Unmarshal(yamlFile, &secrets)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error unmarshaling secrets YAML:", err),
		}
		logMode.PrintLog(myLog)
		return secrets, err
	}

	return secrets, nil
}
