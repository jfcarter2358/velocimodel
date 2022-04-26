package config

import (
	"os"
	"strconv"
)

type ConfigObject struct {
	HTTPHost          string
	HTTPPort          int
	ServiceManagerURL string
}

var Config ConfigObject
var Params map[string]interface{}
var Secrets map[string]interface{}

func LoadConfig() {
	serviceManagerURL := os.Getenv("API_SERVER_SERVICE_MANAGER_URL")
	httpHost := os.Getenv("API_SERVER_HTTP_HOST")
	httpPortString := os.Getenv("API_SERVER_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic(err)
	}
	Config.HTTPPort = httpPort
	Config.HTTPHost = httpHost
	Config.ServiceManagerURL = serviceManagerURL
}
