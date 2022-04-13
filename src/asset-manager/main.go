// main.go

package main

import (
	// "os"
	"asset-manager/asset"
	"asset-manager/ceresdb"
	"asset-manager/config"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

var router *gin.Engine

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	config.LoadConfig()
	routerPort := ":" + strconv.Itoa(config.Config.HTTPPort)
	connection.Initialize(config.Config.DBUsername, config.Config.DBPassword, config.Config.DBHost, config.Config.DBPort)

	if err := ceresdb.VerifyDatabase(config.Config.DBName); err != nil {
		panic(err)
	}
	if err := ceresdb.VerifyCollections(config.Config.DBName); err != nil {
		panic(err)
	}

	asset.LoadAssets()

	os.MkdirAll(config.Config.DataPath, 0755)

	log.Print("Running with port: " + strconv.Itoa(config.Config.HTTPPort))

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run(routerPort)
}
