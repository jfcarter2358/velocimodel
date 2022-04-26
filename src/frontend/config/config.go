package config

import (
	"os"
	"strconv"
)

type ConfigObject struct {
	HTTPHost     string
	HTTPPort     int
	APIServerURL string
}

var Config ConfigObject

func LoadConfig() {
	apiServerURL := os.Getenv("FRONTEND_API_SERVER_URL")
	httpHost := os.Getenv("FRONTEND_HTTP_HOST")
	httpPortString := os.Getenv("FRONTEND_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic(err)
	}
	Config.HTTPPort = httpPort
	Config.HTTPHost = httpHost
	Config.APIServerURL = apiServerURL
}
