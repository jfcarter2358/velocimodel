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

	apiServerURL := os.Getenv("ASSET_MANGER_API_SERVER_URL")
	dataPath := os.Getenv("ASSET_MANAGER_DATA_PATH")
	httpPortString := os.Getenv("ASSET_MANAGER_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic(err)
	}
	Config.HTTPPort = httpPort
	Config.DataPath = dataPath
	Config.APIServerURL = apiServerURL
}

func LoadParamsSecrets() {
	Params = loadFromServiceManager("/api/param")
	Secrets = loadFromServiceManager("/api/secret")
}

func loadFromServiceManager(path string) map[string]interface{} {
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
	return tmpObj[0]
}
