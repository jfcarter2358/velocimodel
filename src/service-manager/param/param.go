package param

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"service-manager/config"
	"service-manager/utils"

	"github.com/jfcarter2358/ceresdb-go/connection"
)

const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"

func LoadParams(dbHost, dbName string, dbPort int) {
	jsonFile, err := os.Open(config.Config.ParamsPath)
	if err != nil {
		panic(err)
	}

	var params map[string]interface{}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &params)
	if err != nil {
		panic(err)
	}
	// Populate with env variable config
	params["db_host"] = dbHost
	params["db_port"] = dbPort
	params["db_name"] = dbName
	err = UpdateParam(params)
	if err != nil {
		panic(err)
	}
}

func DeleteParam(paramNames []string) error {
	queryString := fmt.Sprintf("get record %v.params", config.Config.DB.Name)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(paramNames, datum["name"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.params %v", config.Config.DB.Name, queryData)
	_, err = connection.Query(queryString)
	return err
}

func GetParams(limit, filter, count string) (map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.params", config.Config.DB.Name)
	if filter != FILTER_DEFAULT {
		queryString += fmt.Sprintf(" | filter %v", filter)
	}
	if limit != LIMIT_DEFAULT {
		queryString += fmt.Sprintf(" | limit %v", limit)
	}
	if count != COUNT_DEFAULT {
		queryString += " | count"
	}
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	output := make(map[string]interface{})
	for _, datum := range data {
		output[datum["name"].(string)] = datum["value"]
	}
	return output, nil
}

func RegisterParam(input map[string]interface{}) error {
	queryList := make([]map[string]interface{}, 0)
	for key, val := range input {
		tmpObject := map[string]interface{}{"name": key, "value": val}
		queryList = append(queryList, tmpObject)
	}

	queryData, _ := json.Marshal(&queryList)
	queryString := fmt.Sprintf("post record %v.params %v", config.Config.DB.Name, string(queryData))
	_, err := connection.Query(queryString)
	return err
}

func UpdateParam(input map[string]interface{}) error {
	queryString := fmt.Sprintf("get record %v.params", config.Config.DB.Name)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	newData := make([]map[string]interface{}, 0)
	for key, val := range input {
		exists := false
		for idx, data := range currentData {
			if data["name"] == key {
				data["value"] = val
				currentData[idx] = data
				exists = true
				break
			}
		}
		if !exists {
			tmpObject := map[string]interface{}{"name": key, "value": val}
			newData = append(newData, tmpObject)
		}
	}
	if len(currentData) > 0 {
		queryData, _ := json.Marshal(&currentData)
		queryString = fmt.Sprintf("put record %v.params %v", config.Config.DB.Name, string(queryData))
		_, err = connection.Query(queryString)
		if err != nil {
			return err
		}
	}
	if len(newData) > 0 {
		queryData, _ := json.Marshal(&newData)
		queryString = fmt.Sprintf("post record %v.params %v", config.Config.DB.Name, string(queryData))
		_, err = connection.Query(queryString)
		if err != nil {
			return err
		}
	}
	return nil
}
