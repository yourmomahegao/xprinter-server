//go:build windows
// +build windows

package services

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Printer struct {
	Name string `json:"name"`
}

var (
	winspool = windows.NewLazySystemDLL("winspool.drv")

	// printers
	enumPrintersProc = winspool.NewProc("EnumPrintersW")

	// print
	openPrinter      = winspool.NewProc("OpenPrinterW")
	closePrinter     = winspool.NewProc("ClosePrinter")
	startDocPrinter  = winspool.NewProc("StartDocPrinterW")
	endDocPrinter    = winspool.NewProc("EndDocPrinter")
	writePrinter     = winspool.NewProc("WritePrinter")
	startPagePrinter = winspool.NewProc("StartPagePrinter")
	endPagePrinter   = winspool.NewProc("EndPagePrinter")
)

const (
	PRINTER_ENUM_LOCAL       = 0x00000002
	PRINTER_ENUM_CONNECTIONS = 0x00000004
)

type PRINTER_INFO_4 struct {
	pPrinterName *uint16
	pServerName  *uint16
	Attributes   uint32
}

type DOC_INFO_1 struct {
	pDocName    *uint16
	pOutputFile *uint16
	pDatatype   *uint16
}

func utf16Ptr(s string) *uint16 {
	p, _ := windows.UTF16PtrFromString(s)
	return p
}

func GetPrinters() (printers []Printer, err error) {

	flags := PRINTER_ENUM_LOCAL | PRINTER_ENUM_CONNECTIONS

	var needed uint32
	var returned uint32

	// Buffer size
	r1, _, err := enumPrintersProc.Call(
		uintptr(flags),
		0,
		4,
		0,
		0,
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)

	if r1 == 0 && needed == 0 {
		return nil, err
	}

	buffer := make([]byte, needed)

	// Actual call
	r1, _, err = enumPrintersProc.Call(
		uintptr(flags),
		0,
		4,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(needed),
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)

	if r1 == 0 {
		return nil, err
	}

	items := (*[1 << 20]PRINTER_INFO_4)(
		unsafe.Pointer(&buffer[0]),
	)[:returned:returned]

	for _, item := range items {

		name := windows.UTF16PtrToString(item.pPrinterName)

		if name == "" {
			continue
		}

		printers = append(printers, Printer{
			Name: name,
		})
	}

	return printers, nil
}

func Print(printerModel Printer, data []byte, raw bool) (bool, string) {

	if len(data) == 0 {
		return false, "Print data is empty"
	}

	var handle windows.Handle

	r1, _, err := openPrinter.Call(
		uintptr(unsafe.Pointer(utf16Ptr(printerModel.Name))),
		uintptr(unsafe.Pointer(&handle)),
		0,
	)

	if r1 == 0 {
		return false, fmt.Sprintf(
			"OpenPrinter failed: %v",
			err,
		)
	}

	defer closePrinter.Call(uintptr(handle))

	dataType := ""

	if raw {
		dataType = "RAW"
	}

	docInfo := DOC_INFO_1{
		pDocName:  utf16Ptr("Go Print Job"),
		pDatatype: utf16Ptr(dataType),
	}

	r1, _, err = startDocPrinter.Call(
		uintptr(handle),
		1,
		uintptr(unsafe.Pointer(&docInfo)),
	)

	if r1 == 0 {
		return false, fmt.Sprintf(
			"StartDocPrinter failed: %v",
			err,
		)
	}

	defer endDocPrinter.Call(uintptr(handle))

	// RAW printers usually do not want pages
	if !raw {

		r1, _, err = startPagePrinter.Call(
			uintptr(handle),
		)

		if r1 == 0 {
			return false, fmt.Sprintf(
				"StartPagePrinter failed: %v",
				err,
			)
		}

		defer endPagePrinter.Call(uintptr(handle))
	}

	var written uint32

	r1, _, err = writePrinter.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)

	if r1 == 0 {
		return false, fmt.Sprintf(
			"WritePrinter failed: %v",
			err,
		)
	}

	return true, fmt.Sprintf(
		"Document sent to printer (%d bytes)",
		written,
	)
}
