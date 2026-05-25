package server

import (
	"fmt"
	"log"

	"xprinter/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Run() {
	fmt.Println("-------------------------")
	fmt.Println(" XPrinter Server v.1.0.0 ")
	fmt.Println("-------------------------")
	fmt.Println("")

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.GET("api/get_printers/", handlers.GetPrinters)
	engine.POST("api/print/", handlers.Print)

	if err := engine.Run(":12011"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
