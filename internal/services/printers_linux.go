//go:build !windows
// +build !windows

package services

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/google/uuid"
)

type Printer struct {
	Name string `json:"name"`
}

func GetPrinters() (printers []Printer, err error) {
	cmd := exec.Command("lpstat", "-p")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile(`printer\s+([^\s]+)`)
	matches := r.FindAllStringSubmatch(string(out), -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		newPrinter := Printer{
			Name: match[1],
		}

		printers = append(printers, newPrinter)
	}

	return printers, nil
}

func Print(printer Printer, data []byte, raw bool) (bool, string) {
	tempDir := os.TempDir()
	tempFile := tempDir + "/xprinter-printdata-" + uuid.New().String() + ".prn"

	f, err := os.Create(tempFile)
	if err != nil {
		return false, "Can't create file with print data"
	}

	if _, err := f.Write(data); err != nil {
		return false, "Can't write to print file"
	}

	f.Close()

	args := []string{"-d", printer.Name}

	if raw {
		args = append(args, "-o", "raw")
	}

	args = append(args, tempFile)

	cmd := exec.Command("lp", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, string(out)
	}

	return true, "Document sent to printer"
}
