// Package crawlerMode -----------------------------
// @file      : crawlerScan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/17 22:48
// -------------------------------------------
package crawlerMode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/util"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

type Request struct {
	Method  string `json:"Method"`
	URL     string `json:"URL"`
	B64Body string `json:"b64_body,omitempty"`
}

func CrawlerScan(targets []string) []types.CrawlerResult {
	config.SetUp()
	radModePath := config.CrawlerPath
	radPath := config.CrawlerExecPath
	command := radPath
	timeRandom := util.GetTimeNow()
	strRandom := util.GenerateRandomString(8)
	targetFileName := util.CalculateMD5(timeRandom + strRandom)
	targetPath := filepath.Join(radModePath, "target", targetFileName)
	resultPath := filepath.Join(radModePath, "result", targetFileName)
	radConfigPath := filepath.Join(radModePath, "rad_config.yml")
	fileContent := ""
	for _, target := range targets {
		fileContent += target + "\n"
	}
	flag := util.WriteContentFile(targetPath, fileContent)
	if !flag {
		fmt.Printf("Write target file error")
		return nil
	}
	fmt.Println(targetPath)
	fmt.Println(command)
	fmt.Println(resultPath)
	args := []string{"--url-file", targetPath, "--json", resultPath, "--config", radConfigPath}
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Printf("执行命令时出错：%v\n", err)
		return nil
	}
	fmt.Println(string(output))

	// Read the content of the file
	resultContent, err := ioutil.ReadFile(resultPath)
	if err != nil {
		fmt.Println(err)
	}

	var requests []Request

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(resultContent, &requests)
	if err != nil {
		fmt.Println(err)
	}
	var CrawlerResults []types.CrawlerResult
	// Print the parsed JSON data
	for _, req := range requests {
		fmt.Printf("Method: %s\n", req.Method)
		fmt.Printf("URL: %s\n", req.URL)
		body := ""
		if req.B64Body != "" {
			decodedBytes, err := base64.StdEncoding.DecodeString(req.B64Body)
			if err != nil {
				fmt.Println(err)
			}
			body = string(decodedBytes)
		}
		fmt.Println("-----------------------------")
		crawlerResult := types.CrawlerResult{
			Url:    req.URL,
			Method: req.Method,
			Body:   body,
		}
		CrawlerResults = append(CrawlerResults, crawlerResult)
	}
	util.DeleteFile(targetPath)
	util.DeleteFile(resultPath)
	return CrawlerResults

}
