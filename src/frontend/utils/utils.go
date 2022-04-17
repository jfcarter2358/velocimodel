// utils.go

package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Error(err error, c *gin.Context, statusCode int) {
	log.Printf("Encountered error: %v", err)
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
