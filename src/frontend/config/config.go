package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
)

const DEFAULT_CONFIG_PATH = "/home/frontend/data/config.json"
const ENV_PREFIX = "FRONTEND_"

type ConfigObject struct {
	HTTPHost         string      `json:"http_host" env:"HTTP_HOST"`
	HTTPPort         int         `json:"http_port" env:"HTTP_PORT"`
	ExternalHTTPHost string      `json:"external_http_host" env:"EXTERNAL_HTTP_HOST"`
	APIServerURL     string      `json:"api_server_url" env:"API_SERVER_URL"`
	JoinToken        string      `json:"join_token" env:"JOIN_TOKEN"`
	Oauth            OauthObject `json:"oauth" env:"OAUTH"`
}

type OauthObject struct {
	ClientID              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	AuthServerInternalURL string `json:"auth_manager_internal_url"`
	AuthServerExternalURL string `json:"auth_manager_external_url"`
}

var Config ConfigObject

func LoadConfig() {
	configPath := os.Getenv(ENV_PREFIX + "CONFIG_PATH")
	if configPath == "" {
		configPath = DEFAULT_CONFIG_PATH
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Println("Unable to read json file")
		panic(err)
	}

	log.Printf("Successfully Opened %v", configPath)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &Config)

	v := reflect.ValueOf(Config)
	t := reflect.TypeOf(Config)

	for i := 0; i < v.NumField(); i++ {
		field, found := t.FieldByName(v.Type().Field(i).Name)
		if !found {
			continue
		}

		value := field.Tag.Get("env")
		if value != "" {
			val, present := os.LookupEnv(ENV_PREFIX + value)
			if present {
				w := reflect.ValueOf(&Config).Elem().FieldByName(t.Field(i).Name)
				x := getAttr(&Config, t.Field(i).Name).Kind().String()
				if w.IsValid() {
					switch x {
					case "int", "int64":
						i, err := strconv.ParseInt(val, 10, 64)
						if err == nil {
							w.SetInt(i)
						}
					case "int8":
						i, err := strconv.ParseInt(val, 10, 8)
						if err == nil {
							w.SetInt(i)
						}
					case "int16":
						i, err := strconv.ParseInt(val, 10, 16)
						if err == nil {
							w.SetInt(i)
						}
					case "int32":
						i, err := strconv.ParseInt(val, 10, 32)
						if err == nil {
							w.SetInt(i)
						}
					case "string":
						w.SetString(val)
					case "float32":
						i, err := strconv.ParseFloat(val, 32)
						if err == nil {
							w.SetFloat(i)
						}
					case "float", "float64":
						i, err := strconv.ParseFloat(val, 64)
						if err == nil {
							w.SetFloat(i)
						}
					case "bool":
						i, err := strconv.ParseBool(val)
						if err == nil {
							w.SetBool(i)
						}
					default:
						objValue := reflect.New(field.Type)
						objInterface := objValue.Interface()
						err := json.Unmarshal([]byte(val), objInterface)
						obj := reflect.ValueOf(objInterface)
						if err == nil {
							w.Set(reflect.Indirect(obj).Convert(field.Type))
						} else {
							log.Println(err)
						}
					}
				}
			}
		}
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}

func getAttr(obj interface{}, fieldName string) reflect.Value {
	pointToStruct := reflect.ValueOf(obj) // addressable
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}
	curField := curStruct.FieldByName(fieldName) // type: reflect.Value
	if !curField.IsValid() {
		panic("not found:" + fieldName)
	}
	return curField
}
