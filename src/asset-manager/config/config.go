package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ConfigObject struct {
	DBUsername   string
	DBPassword   string
	DBName       string
	DBHost       string
	DBPort       int
	HTTPPort     int
	DataPath     string
	APIServerURL string
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
	apiServerURL := os.Getenv("ASSET_MANGER_API_SERVER_URL")
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
	Config.APIServerURL = apiServerURL
}

func LoadParamsSecrets() {
	loadFromServiceManager("/api/param", &Params)
	loadFromServiceManager("/api/secret", &Secrets)
}

func loadFromServiceManager(path string, obj *map[string]interface{}) {
	tmpObj := make([]map[string]interface{}, 0)
	resp, err := http.Get(Config.APIServerURL + path)
	if err != nil {
		panic(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(body), &tmpObj)
	if err != nil {
		panic(err)
	}
	obj = &tmpObj[0]
}
