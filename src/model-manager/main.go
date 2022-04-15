// main.go

package main

import (
	// "os"

	"bytes"
	"encoding/json"
	"fmt"
	"model-manager/api"
	"model-manager/ceresdb"
	"model-manager/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jfcarter2358/ceresdb-go/connection"

	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

var router *gin.Engine

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	log := logrus.New()

	config.LoadConfig()

	// Wait for api-server to become available
	statusCode := http.StatusServiceUnavailable
	requestURL := fmt.Sprintf("%v/health", config.Config.APIServerURL)
	for statusCode == http.StatusServiceUnavailable {
		resp, err := http.Get(requestURL)
		if err != nil {
			log.Printf("Error on get to %v: %v", requestURL, err)
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			statusCode = http.StatusOK
		}
	}
	time.Sleep(2 * time.Second)

	config.LoadParamsSecrets()

	routerPort := ":" + strconv.Itoa(config.Config.HTTPPort)
	connection.Initialize(config.Config.DBUsername, config.Config.DBPassword, config.Config.DBHost, config.Config.DBPort)

	if err := ceresdb.VerifyDatabase(config.Config.DBName); err != nil {
		panic(err)
	}
	if err := ceresdb.VerifyCollections(config.Config.DBName); err != nil {
		panic(err)
	}

	values := map[string]interface{}{"host": "model-manager", "port": config.Config.HTTPPort, "type": "model-manager"}
	json_data, err := json.Marshal(values)

	if err != nil {
		panic(err)
	}

	_, err = http.Post(fmt.Sprintf("%v/api/service", config.Config.APIServerURL), "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		panic(err)
	}

	log.Print("Running with port: " + strconv.Itoa(config.Config.HTTPPort))

	api.Healthy = true

	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run(routerPort)
}
