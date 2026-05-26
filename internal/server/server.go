package server

import (
	"fmt"
	"log"

	"xprinter/internal/config"
	"xprinter/internal/handlers"
	"xprinter/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Run() {
	fmt.Println("-------------------------")
	fmt.Printf(" XPrinter Server v.%s \n\r", config.VERSION)
	fmt.Println("-------------------------")
	fmt.Println("")

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware())

	engine.GET("/", handlers.Index)
	engine.GET("api/version/", handlers.Version)
	engine.GET("api/printers/", handlers.GetPrinters)
	engine.POST("api/print/", handlers.Print)

	if err := engine.Run(":12011"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
