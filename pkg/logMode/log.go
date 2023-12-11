// Package logMode -----------------------------
// @file      : log.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2023/12/9 20:07
// -------------------------------------------
package logMode

import "fmt"

type CustomLog struct {
	Status string
	Msg    string
}

func PrintLog(l CustomLog) {
	fmt.Printf("%s: %s\n", l.Status, l.Msg)
}
