// main.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/api"
	"frontend/auth"
	"frontend/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

var router *gin.Engine

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// logrus.SetFormatter(&logrus.JSONFormatter{
	// 	FieldMap: logrus.FieldMap{
	// 		logrus.FieldKeyTime: "@timestamp",
	// 		logrus.FieldKeyMsg:  "message",
	// 	},
	// })
	log := logrus.New()

	config.LoadConfig()
	auth.LoadOauthConfig()

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

	routerPort := ":" + strconv.Itoa(config.Config.HTTPPort)

	values := map[string]interface{}{"host": config.Config.HTTPHost, "port": config.Config.HTTPPort, "type": "frontend"}
	json_data, err := json.Marshal(values)

	if err != nil {
		panic(err)
	}

	requestURL = fmt.Sprintf("%v/api/service", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", config.Config.JoinToken))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)

	if err != nil {
		panic(err)
	}

	log.Print("Running with port: " + strconv.Itoa(config.Config.HTTPPort))

	api.Healthy = true

	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.Use(ginlogrus.Logger(log), gin.Recovery())
	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run(routerPort)
}
