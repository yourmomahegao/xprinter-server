package handlers

import (
	"xprinter/internal/config"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Service working.",
	})
}

func Version(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": true,
		"data":   config.VERSION,
	})
}
