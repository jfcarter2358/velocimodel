package asset

import (
	"asset-manager/config"
	"asset-manager/utils"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
	"log"

	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

type Asset struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	URL      string                 `json:"url"`
	Created  string                 `json:"created"`
	Updated  string                 `json:"updated"`
	Type     string                 `json:"type"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Models   []string               `json:"models"`
}

const FILE_TYPE = "file"
const GIT_TYPE = "git"
const S3_TYPE = "s3"
const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"

var allowedTypes = []string{FILE_TYPE, GIT_TYPE, S3_TYPE}

func RegisterAsset(newAsset Asset) error {
	countObj, _ := connection.Query(fmt.Sprintf("get record %v.assets | count", config.Config.DBName))
	startCount := int(countObj[0]["count"].(float64))
	endCount := startCount
	if newAsset.ID == "" {
		newAsset.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newAsset.Type) {
		err := errors.New(fmt.Sprintf("Asset type of %v does not exist", newAsset.Type))
		return err
	}
	currentTime := time.Now().UTC()
	newAsset.Created = currentTime.Format("2006-01-02T15:04:05Z")
	newAsset.Updated = currentTime.Format("2006-01-02T15:04:05Z")
	queryData, _ := json.Marshal(&newAsset)
	for endCount == startCount {
		queryString := fmt.Sprintf("post record %v.assets %v", config.Config.DBName, string(queryData))
		_, err := connection.Query(queryString)
		if err != nil {
			return err
		}
		countObj, _ := connection.Query(fmt.Sprintf("get record %v.assets | count", config.Config.DBName))
		endCount = int(countObj[0]["count"].(float64))
	}
	return nil
}

func DeleteAsset(assetIDs []string) error {
	queryString := fmt.Sprintf("get record %v.assets", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(assetIDs, datum["id"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.assets %v", config.Config.DBName, string(queryData))
	_, err = connection.Query(queryString)
	return err
}

func UpdateAsset(newAsset Asset) error {
	if newAsset.ID == "" {
		err := errors.New("'id' field is required to update an asset")
		return err
	}
	queryString := fmt.Sprintf("get record %v.assets", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	for _, datum := range currentData {
		if datum["id"].(string) == newAsset.ID {
			if newAsset.Name != "" {
				if newAsset.Name != datum["name"].(string) {
					datum["name"] = newAsset.Name
				}
			}
			if newAsset.Tags != nil {
				tmpTags := make([]string, len(datum["tags"].([]interface{})))
				for idx, val := range datum["tags"].([]interface{}) {
					tmpTags[idx] = val.(string)
				}
				if !reflect.DeepEqual(newAsset.Tags, tmpTags) {
					datum["tags"] = newAsset.Tags
				}
			}
			if newAsset.Metadata != nil {
				if !reflect.DeepEqual(newAsset.Metadata, datum["metadata"].(map[string]interface{})) {
					datum["metadata"] = newAsset.Metadata
				}
			}
			if newAsset.Models != nil {
				tmpModels := make([]string, len(datum["models"].([]interface{})))
				for idx, val := range datum["models"].([]interface{}) {
					tmpModels[idx] = val.(string)
				}
				if !reflect.DeepEqual(newAsset.Models, tmpModels) {
					datum["models"] = newAsset.Models
				}
			}
			currentTime := time.Now().UTC()
			datum["updated"] = currentTime.Format("2006-01-02T15:04:05Z")
			queryData, _ := json.Marshal(&datum)
			queryString := fmt.Sprintf("put record %v.assets %v", config.Config.DBName, string(queryData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
			return nil
		}
	}
	err = RegisterAsset(newAsset)
	return err
}

func GetAssets(limit, filter, count string) ([]map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.assets", config.Config.DBName)
	if filter != FILTER_DEFAULT {
		queryString += fmt.Sprintf(" | filter %v", filter)
	}
	if limit != LIMIT_DEFAULT {
		queryString += fmt.Sprintf(" | limit %v", limit)
	}
	if count != COUNT_DEFAULT {
		queryString += " | count"
	}
	log.Printf("QUERYSTRING: %v", queryString)
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	return data, nil
}
