package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

const DEFAULT_CONFIG_PATH = "/home/auth-manager/data/config.json"

type ConfigObject struct {
	DBUsername   string
	DBPassword   string
	DBName       string
	DBHost       string
	DBPort       int
	HTTPHost     string         `json:"http_host" env:"AUTH_MANAGER_HTTP_HOST"`
	HTTPPort     int            `json:"http_port" env:"AUTH_MANAGER_HTTP_PORT"`
	Clients      []ClientObject `json:"clients" env:"AUTH_MANAGER_CLIENTS"`
	LDAP         LDAPObject     `json:"ldap" env:"AUTH_MANAGER_LDAP"`
	Admin        AdminObject    `json:"admin" env:"AUTH_MANAGER_ADMIN"`
	URL          string         `json:"url" env:"AUTH_MANAGER_URL"`
	APIServerURL string         `json:"api_server_url" env:"AUTH_MANAGER_API_SERVER_URL"`
}

type ClientObject struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

type LDAPObject struct {
	Enabled      bool           `json:"enabled"`
	BaseDN       string         `json:"base_dn"`
	BindDN       string         `json:"bind_dn"`
	Port         int            `json:"port"`
	Host         string         `json:"host"`
	BindPassword string         `json:"bind_password"`
	Filter       string         `json:"filter"`
	Keys         LDAPKeysObject `json:"keys"`
}

type LDAPKeysObject struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Groups     string `json:"groups"`
	Roles      string `json:"roles"`
	Email      string `json:"email"`
}

type AdminObject struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Config ConfigObject
var Params map[string]interface{}
var Secrets map[string]interface{}

func LoadConfig() {
	configPath := os.Getenv("AUTH_MANAGER_CONFIG_PATH")
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
		log.Printf("%v: %v", field.Name, field.Type)
		if value != "" {
			val, present := os.LookupEnv(value)
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
