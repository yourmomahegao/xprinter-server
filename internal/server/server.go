package server

import (
	"log"

	"xprinter/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.GET("api/get_printers/", handlers.GetPrinters)
	engine.POST("api/print/", handlers.Print)

	if err := engine.Run(":12011"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
