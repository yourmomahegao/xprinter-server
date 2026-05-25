package handlers

import (
	"fmt"
	"xprinter/internal/services"

	"github.com/gin-gonic/gin"
)

func GetPrinters(c *gin.Context) {
	printers, printersError := services.GetPrinters()

	if printersError != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Can't get printers list",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  true,
		"message": "",
		"data":    printers,
	})
}

type PrintRequest struct {
	Printer string `form:"printer" json:"printer" binding:"required"`
	Data    string `form:"data" json:"data" binding:"required"`
	Raw     bool   `form:"raw" json:"raw"`
}

func Print(c *gin.Context) {
	var req PrintRequest
	var help string = `Usage: /api/print/
	- printer: Name of the printer (required)
	- data: Data to print in printer (required)
	- raw: Is data is raw (not required) (default: false)`

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": help,
		})
		return
	}

	printStatus, printMessage := services.Print(
		services.Printer{Name: req.Printer},
		[]byte(req.Data),
		req.Raw,
	)

	if printStatus {
		c.JSON(200, gin.H{
			"status":  true,
			"message": fmt.Sprintf("All good. %s", printMessage),
		})
	} else {
		c.JSON(502, gin.H{
			"status":  false,
			"message": fmt.Sprintf("Can't send document to printer. %s", printMessage),
		})
	}

}
