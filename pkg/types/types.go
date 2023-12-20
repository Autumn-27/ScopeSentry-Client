// Package _type -----------------------------
// @file      : type.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/11 10:05
// -------------------------------------------
package types

import (
	"encoding/json"
	"github.com/projectdiscovery/tlsx/pkg/tlsx/clients"
)

type SubdomainResult struct {
	Host  string
	Type  string
	Value []string
	IP    []string
	Time  string
}

type AssertHttp struct {
	Timestamp     string                 `json:"timestamp,omitempty" csv:"timestamp"`
	TLSData       *clients.Response      `json:"tls,omitempty" csv:"tls"`
	Hashes        map[string]interface{} `json:"hash,omitempty" csv:"hash"`
	CDNName       string                 `json:"cdn_name,omitempty" csv:"cdn_name"`
	Port          string                 `json:"port,omitempty" csv:"port"`
	URL           string                 `json:"url,omitempty" csv:"url"`
	Location      string                 `json:"location,omitempty" csv:"location"`
	Title         string                 `json:"title,omitempty" csv:"title"`
	Type          string                 `json:"Type,omitempty" csv:"Type"`
	Error         string                 `json:"error,omitempty" csv:"error"`
	ResponseBody  string                 `json:"body,omitempty" csv:"body"`
	Host          string                 `json:"host,omitempty" csv:"host"`
	FavIconMMH3   string                 `json:"favicon,omitempty" csv:"favicon"`
	FaviconPath   string                 `json:"favicon_path,omitempty" csv:"favicon_path"`
	RawHeaders    string                 `json:"raw_header,omitempty" csv:"raw_header"`
	Jarm          string                 `json:"jarm,omitempty" csv:"jarm"`
	Technologies  []string               `json:"tech,omitempty" csv:"tech"`
	StatusCode    int                    `json:"status_code,omitempty" csv:"status_code"`
	ContentLength int                    `json:"content_length,omitempty" csv:"content_length"`
	CDN           bool                   `json:"cdn,omitempty" csv:"cdn"`
	Webcheck      bool                   `json:"cdn,omitempty" csv:"cdn"`
}

type PortAlive struct {
	Host string `json:"Host,omitempty"`
	IP   string `json:"Host,omitempty"`
	Port int    `json:"Port,omitempty"`
}

type AssertOther struct {
	Host      string          `json:"host,omitempty"`
	IP        string          `json:"ip"`
	Port      int             `json:"port"`
	Protocol  string          `json:"protocol"`
	TLS       bool            `json:"tls"`
	Transport string          `json:"transport"`
	Version   string          `json:"version,omitempty"`
	Raw       json.RawMessage `json:"metadata"`
}

type UrlResult struct {
	Input      string `json:"input"`
	Source     string `json:"source"`
	OutputType string `json:"type"`
	Output     string `json:"output"`
	StatusCode int    `json:"status"`
	Length     int    `json:"length"`
	Time       string `json:"time"`
}

type SecretResults struct {
	Url   string
	Kind  string
	Key   string
	Value string
}

type CrawlerResult struct {
	Url    string
	Method string
	Body   string
}
