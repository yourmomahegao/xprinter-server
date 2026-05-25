//go:build !windows
// +build !windows

package main

import (
	"xprinter/internal/server"
)

func main() {
	server.Run()
}
