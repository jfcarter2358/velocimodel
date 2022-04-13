package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ConfigObject struct {
	DBUsername        string
	DBPassword        string
	DBName            string
	DBHost            string
	DBPort            int
	HTTPPort          int
	DataPath          string
	ServiceManagerURL string
}

var Config ConfigObject
var Params map[string]interface{}
var Secrets map[string]interface{}

func LoadConfig() {
	Params = make(map[string]interface{})
	Secrets = make(map[string]interface{})

	dbUsername := os.Getenv("ASSET_MANAGER_DB_USERNAME")
	dbPassword := os.Getenv("ASSET_MANAGER_DB_PASSWORD")
	dbHost := os.Getenv("ASSET_MANAGER_DB_HOST")
	dbName := os.Getenv("ASSET_MANAGER_DB_NAME")
	dbPortString := os.Getenv("ASSET_MANAGER_DB_PORT")
	serviceManagerURL := os.Getenv("ASSET_MANGER_SERVICE_MANAGER_URL")
	dataPath := os.Getenv("ASSET_MANAGER_DATA_PATH")
	dbPort, err := strconv.Atoi(dbPortString)
	if err != nil {
		panic(err)
	}
	httpPortString := os.Getenv("ASSET_MANAGER_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic(err)
	}
	Config.DBUsername = dbUsername
	Config.DBPassword = dbPassword
	Config.DBHost = dbHost
	Config.DBName = dbName
	Config.DBPort = dbPort
	Config.HTTPPort = httpPort
	Config.DataPath = dataPath
	Config.ServiceManagerURL = serviceManagerURL

	loadFromConfigManager("/api/param", &Params)
	loadFromConfigManager("/api/secret", &Secrets)

	log.Println(fmt.Sprintf("%v", Params))
}

func loadFromConfigManager(path string, obj *map[string]interface{}) {
	resp, err := http.Get(Config.ServiceManagerURL + path)
	if err != nil {
		panic(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))

	err = json.Unmarshal([]byte(body), obj)
	if err != nil {
		panic(err)
	}
}
