// utils.go

package utils

import (
	"fmt"
	"frontend/logging"

	"github.com/gin-gonic/gin"
)

func Error(err error, c *gin.Context, statusCode int) {
	logging.Logger.Error(fmt.Sprintf("Encountered error: %v", err))
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
