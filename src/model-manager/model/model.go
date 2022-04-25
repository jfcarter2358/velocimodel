package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"model-manager/config"
	"model-manager/utils"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

type Model struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Created   string                 `json:"created"`
	Updated   string                 `json:"updated"`
	Type      string                 `json:"type"`
	Tags      []string               `json:"tags"`
	Metadata  map[string]interface{} `json:"metadata"`
	Assets    []string               `json:"assets"`
	Snapshots []string               `json:"snapshots"`
	Releases  []string               `json:"releases"`
	Language  string                 `json:"language"`
}

const GIT_TYPE = "git"
const RAW_TYPE = "raw"
const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"
const ORDERASC_DEFAULT = "NA"
const ORDERDSC_DEFAULT = "NA"

var allowedTypes = []string{RAW_TYPE}

func AddAsset(modelID, assetID string) error {
	currentTime := time.Now().UTC()
	updatedTime := currentTime.Format("2006-01-02T15:04:05Z")
	log.Printf("get record %v.models | filter id = \"%v\"", config.Config.DBName, modelID)
	modelObjects, err := connection.Query(fmt.Sprintf("get record %v.models | filter id = \"%v\"", config.Config.DBName, modelID))
	if err != nil {
		return err
	}
	if len(modelObjects) == 0 {
		return errors.New(fmt.Sprintf("Model not found with ID %v", modelID))
	}
	realModelID := modelObjects[0][".id"].(string)
	assetList := modelObjects[0]["assets"].([]interface{})
	assetList = append(assetList, assetID)
	assetBytes, _ := json.Marshal(assetList)
	log.Printf("patch record %v.models \"%v\" {\"updated\":\"%v\",\"assets\":%v}", config.Config.DBName, realModelID, updatedTime, string(assetBytes))
	_, err = connection.Query(fmt.Sprintf("patch record %v.models \"%v\" {\"updated\":\"%v\",\"assets\":%v}", config.Config.DBName, realModelID, updatedTime, string(assetBytes)))

	log.Printf("get record %v.assets | filter id = \"%v\"", config.Config.DBName, assetID)
	assetObjects, err := connection.Query(fmt.Sprintf("get record %v.assets | filter id = \"%v\"", config.Config.DBName, assetID))
	if err != nil {
		return err
	}
	if len(assetObjects) == 0 {
		return errors.New(fmt.Sprintf("Asset not found with ID %v", assetID))
	}
	realAssetID := assetObjects[0][".id"].(string)
	modelList := assetObjects[0]["models"].([]interface{})
	modelList = append(modelList, modelID)
	modelBytes, _ := json.Marshal(modelList)
	log.Printf("patch record %v.assets \"%v\" {\"updated\":\"%v\",\"models\":%v}", config.Config.DBName, realAssetID, updatedTime, string(modelBytes))
	_, err = connection.Query(fmt.Sprintf("patch record %v.assets \"%v\" {\"updated\":\"%v\",\"models\":%v}", config.Config.DBName, realAssetID, updatedTime, string(modelBytes)))
	return nil
}

func DeleteAsset(modelID, assetID string) error {
	currentTime := time.Now().UTC()
	updatedTime := currentTime.Format("2006-01-02T15:04:05Z")
	modelObjects, err := connection.Query(fmt.Sprintf("get record %v.models | filter id = \"%v\"", config.Config.DBName, modelID))
	if err != nil {
		return err
	}
	if len(modelObjects) == 0 {
		return errors.New(fmt.Sprintf("Model not found with ID %v", modelID))
	}
	realModelID := modelObjects[0][".id"].(string)
	assetList := modelObjects[0]["assets"].([]interface{})
	for idx, val := range assetList {
		if val.(string) == assetID {
			assetList = append(assetList[:idx], assetList[idx+1:]...)
			break
		}
	}
	assetBytes, _ := json.Marshal(assetList)
	_, err = connection.Query(fmt.Sprintf("patch record %v.models \"%v\" {\"updated\":\"%v\",\"assets\":%v}", config.Config.DBName, realModelID, updatedTime, string(assetBytes)))

	assetObjects, err := connection.Query(fmt.Sprintf("get record %v.assets | filter id = \"%v\"", config.Config.DBName, assetID))
	if err != nil {
		return err
	}
	if len(assetObjects) == 0 {
		return errors.New(fmt.Sprintf("Asset not found with ID %v", assetID))
	}
	realAssetID := assetObjects[0][".id"].(string)
	modelList := assetObjects[0]["models"].([]interface{})
	for idx, val := range modelList {
		if val.(string) == modelID {
			modelList = append(modelList[:idx], modelList[idx+1:]...)
			break
		}
	}
	modelBytes, _ := json.Marshal(modelList)
	_, err = connection.Query(fmt.Sprintf("patch record %v.assets \"%v\" {\"updated\":\"%v\",\"models\":%v}", config.Config.DBName, realAssetID, updatedTime, string(modelBytes)))
	return nil
}

func RegisterModel(newModel Model) error {
	countObj, _ := connection.Query(fmt.Sprintf("get record %v.models | count", config.Config.DBName))
	startCount := int(countObj[0]["count"].(float64))
	endCount := startCount
	if newModel.ID == "" {
		newModel.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newModel.Type) {
		err := errors.New(fmt.Sprintf("Model type of %v does not exist", newModel.Type))
		return err
	}
	currentTime := time.Now().UTC()
	newModel.Created = currentTime.Format("2006-01-02T15:04:05Z")
	newModel.Updated = currentTime.Format("2006-01-02T15:04:05Z")
	queryData, _ := json.Marshal(&newModel)
	for endCount == startCount {
		queryString := fmt.Sprintf("post record %v.models %v", config.Config.DBName, string(queryData))
		_, err := connection.Query(queryString)
		if err != nil {
			return err
		}
		countObj, _ := connection.Query(fmt.Sprintf("get record %v.models | count", config.Config.DBName))
		endCount = int(countObj[0]["count"].(float64))
	}
	return nil
}

func DeleteModel(modelIDs []string) error {
	queryString := fmt.Sprintf("get record %v.models", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(modelIDs, datum["id"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.models %v", config.Config.DBName, string(queryData))
	_, err = connection.Query(queryString)
	return err
}

func UpdateModel(newModel Model) error {
	if newModel.ID == "" {
		err := errors.New("'id' field is required to update an model")
		return err
	}
	queryString := fmt.Sprintf("get record %v.models", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	for _, datum := range currentData {
		if datum["id"].(string) == newModel.ID {
			if newModel.Name != "" {
				if newModel.Name != datum["name"].(string) {
					datum["name"] = newModel.Name
				}
			}
			if newModel.Tags != nil {
				tmpTags := make([]string, len(datum["tags"].([]interface{})))
				for idx, val := range datum["tags"].([]interface{}) {
					tmpTags[idx] = val.(string)
				}
				if !reflect.DeepEqual(newModel.Tags, tmpTags) {
					datum["tags"] = newModel.Tags
				}
			}
			if newModel.Metadata != nil {
				if !reflect.DeepEqual(newModel.Metadata, datum["metadata"].(map[string]interface{})) {
					datum["metadata"] = newModel.Metadata
				}
			}
			if newModel.Assets != nil {
				tmpAssets := make([]string, len(datum["assets"].([]interface{})))
				for idx, val := range datum["assets"].([]interface{}) {
					tmpAssets[idx] = val.(string)
				}
				if !reflect.DeepEqual(newModel.Assets, tmpAssets) {
					datum["assets"] = newModel.Assets
				}
			}
			if newModel.Snapshots != nil {
				tmpSnapshots := make([]string, len(datum["snapshots"].([]interface{})))
				for idx, val := range datum["snapshots"].([]interface{}) {
					tmpSnapshots[idx] = val.(string)
				}
				if !reflect.DeepEqual(newModel.Snapshots, tmpSnapshots) {
					datum["snapshots"] = newModel.Snapshots
				}
			}
			if newModel.Releases != nil {
				tmpReleases := make([]string, len(datum["releases"].([]interface{})))
				for idx, val := range datum["releases"].([]interface{}) {
					tmpReleases[idx] = val.(string)
				}
				if !reflect.DeepEqual(newModel.Releases, tmpReleases) {
					datum["releases"] = newModel.Releases
				}
			}
			currentTime := time.Now().UTC()
			datum["updated"] = currentTime.Format("2006-01-02T15:04:05Z")
			queryData, _ := json.Marshal(&datum)
			queryString := fmt.Sprintf("put record %v.models %v", config.Config.DBName, string(queryData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
			return nil
		}
	}
	err = RegisterModel(newModel)
	return err
}

func GetModels(limit, filter, count, orderasc, orderdsc string) ([]map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.models", config.Config.DBName)
	if filter != FILTER_DEFAULT {
		queryString += fmt.Sprintf(" | filter %v", filter)
	}
	if limit != LIMIT_DEFAULT {
		queryString += fmt.Sprintf(" | limit %v", limit)
	}
	if orderasc != ORDERASC_DEFAULT {
		queryString += fmt.Sprintf(" | orderasc %v", orderasc)
	}
	if orderdsc != ORDERDSC_DEFAULT {
		queryString += fmt.Sprintf(" | orderdsc %v", orderdsc)
	}
	if count != COUNT_DEFAULT {
		queryString += " | count"
	}
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	return data, nil
}
