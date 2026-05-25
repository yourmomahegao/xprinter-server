//go:build windows
// +build windows

package main

import (
	"xprinter/internal/platform"
)

func main() {
	platform.RunWindowsService()
}
