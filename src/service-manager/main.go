// main.go

package main

import (
	// "os"

	"service-manager/api"
	"service-manager/ceresdb"
	"service-manager/config"
	"service-manager/param"
	"service-manager/secret"
	"service-manager/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	routerPort := ":" + strconv.Itoa(config.Config.HTTPPort)
	connection.Initialize(config.Config.DBUsername, config.Config.DBPassword, config.Config.DBHost, config.Config.DBPort)

	if err := ceresdb.VerifyDatabase(config.Config.DBName); err != nil {
		panic(err)
	}
	if err := ceresdb.VerifyCollections(config.Config.DBName); err != nil {
		panic(err)
	}

	param.LoadParams(config.Config.DBHost, config.Config.DBName, config.Config.DBPort)
	secret.LoadSecrets(config.Config.DBUsername, config.Config.DBPassword)

	serviceID := uuid.New().String()

	selfService := service.Service{
		ID:   serviceID,
		Host: "service-manager",
		Port: config.Config.HTTPPort,
		Type: "service-manager",
	}

	err := service.RegisterService(selfService)
	if err != nil {
		panic(err)
	}

	api.Healthy = true

	log.Print("Running with port: " + strconv.Itoa(config.Config.HTTPPort))

	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run(routerPort)
}
