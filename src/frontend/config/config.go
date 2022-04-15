package config

import (
	"os"
	"strconv"
)

type ConfigObject struct {
	HTTPPort     int
	APIServerURL string
}

var Config ConfigObject

func LoadConfig() {
	apiServerURL := os.Getenv("FRONTEND_API_SERVER_URL")
	httpPortString := os.Getenv("FRONTEND_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic(err)
	}
	Config.HTTPPort = httpPort
	Config.APIServerURL = apiServerURL
}
