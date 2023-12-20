// Package util -----------------------------
// @file      : util.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/11 10:13
// -------------------------------------------
package util

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/config"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/logMode"
	"github.com/Autumn-27/ScopeSentry-Client/pkg/types"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetTimeNow() string {
	utilConfig := config.ParseConfig()
	TimeZoneName := utilConfig.Time.TimeZoneName
	// 获取当前时间
	timeLocation, err := time.LoadLocation(TimeZoneName)
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("[Err] Error loading time zone:", err),
		}
		logMode.PrintLog(myLog)
		return ""
	}
	currentTime := time.Now()
	var easternTime = currentTime.In(timeLocation)
	return easternTime.Format(time.RFC3339)
}

func ReadFileLiness(fileName string) ([]string, error) {
	result := []string{}
	f, err := os.Open(fileName)
	if err != nil {
		return result, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

func StartTime(name string) {
	fmt.Println("[*]start %s time: %s", name, GetTimeNow())
}

func EndTime(name string) {
	fmt.Println("[*]end %s time: %s", name, GetTimeNow())
}

func GenerateRandomLetters(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func IsWildCard(domain string) bool {
	for i := 0; i < 2; i++ {
		subdomain := GenerateRandomLetters(6) + "." + domain
		_, err := net.LookupIP(subdomain)
		if err != nil {
			continue
		}
		return true
	}
	return false
}

func WriteStructArrayToFile(filename string, Result []types.AssertHttp) error {
	// 将结构体数组编码为 JSON 格式
	jsonData, err := json.MarshalIndent(Result, "", "  ")
	if err != nil {
		return err
	}

	// 将 JSON 数据写入文件
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("结构体数组写入文件 %s 成功\n", filename)
	return nil
}
func GetFileExtension(url string) string {
	// 从URL中提取文件路径
	filePath := filepath.Base(url)

	// 获取文件的后缀名
	fileExtension := filepath.Ext(filePath)

	// 去除后缀名前的点号
	fileExtension = strings.TrimPrefix(fileExtension, ".")

	return fileExtension
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	// 定义字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 构建随机字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func WriteContentFile(filPath string, fileContent string) bool {
	// 将字符串写入文件
	err := ioutil.WriteFile(filPath, []byte(fileContent), 0644)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to create filPath:", filPath, err),
		}
		logMode.PrintLog(myLog)
		return false
	}
	return true

}

func CalculateMD5(input string) string {
	// Convert the input string to bytes
	data := []byte(input)

	// Calculate the MD5 hash
	hash := md5.Sum(data)

	// Convert the hash to a hex string
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

func DeleteFile(filePath string) {
	// 调用Remove函数删除文件
	err := os.Remove(filePath)
	if err != nil {
		myLog := logMode.CustomLog{
			Status: "Error",
			Msg:    fmt.Sprintf("Failed to DeleteFile:", filePath, err),
		}
		logMode.PrintLog(myLog)
		return
	}

}
