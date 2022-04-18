package api

import (
	"log"
	"net/http"
	"service-manager/param"
	"service-manager/secret"
	"service-manager/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"

var Healthy = true

// Health API

func GetHealth(c *gin.Context) {
	if Healthy == false {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Status(http.StatusOK)
	return
}

// Service API

func DeleteService(c *gin.Context) {
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.DeleteService(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func GetServices(c *gin.Context) {
	filter := FILTER_DEFAULT
	limit := LIMIT_DEFAULT
	count := COUNT_DEFAULT
	if val, ok := c.GetQuery("filter"); ok {
		filter = val
	}
	if val, ok := c.GetQuery("limit"); ok {
		limit = val
	}
	if val, ok := c.GetQuery("count"); ok {
		count = val
	}
	data, err := service.GetServices(limit, filter, count)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostService(c *gin.Context) {
	var input service.Service
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	serviceID := uuid.New().String()
	input.ID = serviceID
	err := service.RegisterService(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": serviceID})
}

func PutService(c *gin.Context) {
	var input service.Service
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.UpdateService(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// Param API

func DeleteParam(c *gin.Context) {
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := param.DeleteParam(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func GetParams(c *gin.Context) {
	filter := FILTER_DEFAULT
	limit := LIMIT_DEFAULT
	count := COUNT_DEFAULT
	if val, ok := c.GetQuery("filter"); ok {
		filter = val
	}
	if val, ok := c.GetQuery("limit"); ok {
		limit = val
	}
	if val, ok := c.GetQuery("count"); ok {
		count = val
	}
	data, err := param.GetParams(limit, filter, count)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostParam(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := param.RegisterParam(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func PutParam(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := param.UpdateParam(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func DeleteSecret(c *gin.Context) {
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := secret.DeleteSecret(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func GetSecrets(c *gin.Context) {
	filter := FILTER_DEFAULT
	limit := LIMIT_DEFAULT
	count := COUNT_DEFAULT
	if val, ok := c.GetQuery("filter"); ok {
		filter = val
	}
	if val, ok := c.GetQuery("limit"); ok {
		limit = val
	}
	if val, ok := c.GetQuery("count"); ok {
		count = val
	}
	data, err := secret.GetSecrets(limit, filter, count)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostSecret(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := secret.RegisterSecret(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func PutSecret(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := secret.UpdateSecret(input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
