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

func LoadParams() {
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
	err = UpdateParams(params)
	if err != nil {
		panic(err)
	}
}

func DeleteParams(paramNames []string) error {
	queryString := fmt.Sprintf("get record %v.config", config.Config.DBName)
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
	queryString = fmt.Sprintf("delete record %v.config %v", config.Config.DBName, queryData)
	_, err = connection.Query(queryString)
	return err
}

func GetParams() (map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.config", config.Config.DBName)
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

func SetParams(input map[string]interface{}) error {
	queryList := make([]map[string]interface{}, 0)
	for key, val := range input {
		tmpObject := map[string]interface{}{"name": key, "value": val}
		queryList = append(queryList, tmpObject)
	}

	queryData, _ := json.Marshal(&queryList)
	queryString := fmt.Sprintf("post record %v.config %v", config.Config.DBName, string(queryData))
	_, err := connection.Query(queryString)
	return err
}

func UpdateParams(input map[string]interface{}) error {
	queryString := fmt.Sprintf("get record %v.config", config.Config.DBName)
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

	queryData, _ := json.Marshal(&currentData)
	queryString = fmt.Sprintf("put record %v.config %v", config.Config.DBName, string(queryData))
	_, err = connection.Query(queryString)
	if err != nil {
		return err
	}
	if len(newData) > 0 {
		queryData, _ := json.Marshal(&newData)
		queryString = fmt.Sprintf("post record %v.config %v", config.Config.DBName, string(queryData))
		_, err = connection.Query(queryString)
		if err != nil {
			return err
		}
	}
	return nil
}
