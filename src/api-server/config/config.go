package config

import (
	"os"
	"strconv"
)

type ConfigObject struct {
	HTTPPort          int
	ServiceManagerURL string
}

var Config ConfigObject
var Params map[string]interface{}
var Secrets map[string]interface{}

func LoadConfig() {
	serviceManagerURL := os.Getenv("API_SERVER_SERVICE_MANAGER_URL")
	httpPortString := os.Getenv("API_SERVER_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic(err)
	}
	Config.HTTPPort = httpPort
	Config.ServiceManagerURL = serviceManagerURL
}
