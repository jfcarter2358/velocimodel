package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"model-manager/config"
	"model-manager/utils"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

type Snapshot struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Created  string                 `json:"created"`
	Updated  string                 `json:"updated"`
	Type     string                 `json:"type"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Assets   []string               `json:"assets"`
	Language string                 `json:"langauge"`
}

const GIT_TYPE = "git"
const RAW_TYPE = "raw"
const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"

var allowedTypes = []string{GIT_TYPE, RAW_TYPE}

func RegisterSnapshot(newSnapshot Snapshot) error {
	countObj, _ := connection.Query(fmt.Sprintf("get record %v.snapshots | count", config.Config.DBName))
	startCount := int(countObj[0]["count"].(float64))
	endCount := startCount
	if newSnapshot.ID == "" {
		newSnapshot.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newSnapshot.Type) {
		err := errors.New(fmt.Sprintf("Snapshot type of %v does not exist", newSnapshot.Type))
		return err
	}
	currentTime := time.Now().UTC()
	newSnapshot.Created = currentTime.Format("2006-01-02T15:04:05Z")
	newSnapshot.Updated = currentTime.Format("2006-01-02T15:04:05Z")
	queryData, _ := json.Marshal(&newSnapshot)
	for endCount == startCount {
		queryString := fmt.Sprintf("post record %v.snapshots %v", config.Config.DBName, string(queryData))
		_, err := connection.Query(queryString)
		if err != nil {
			return err
		}
		countObj, _ := connection.Query(fmt.Sprintf("get record %v.snapshots | count", config.Config.DBName))
		endCount = int(countObj[0]["count"].(float64))
	}
	return nil
}

func DeleteSnapshot(snapshotIDs []string) error {
	queryString := fmt.Sprintf("get record %v.snapshots", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(snapshotIDs, datum["id"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.snapshots %v", config.Config.DBName, string(queryData))
	_, err = connection.Query(queryString)
	return err
}

func UpdateSnapshot(newSnapshot Snapshot) error {
	if newSnapshot.ID == "" {
		err := errors.New("'id' field is required to update an snapshot")
		return err
	}
	queryString := fmt.Sprintf("get record %v.snapshots", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	for _, datum := range currentData {
		if datum["id"].(string) == newSnapshot.ID {
			if newSnapshot.Name != "" {
				if newSnapshot.Name != datum["name"].(string) {
					datum["name"] = newSnapshot.Name
				}
			}
			if newSnapshot.Tags != nil {
				tmpTags := make([]string, len(datum["tags"].([]interface{})))
				for idx, val := range datum["tags"].([]interface{}) {
					tmpTags[idx] = val.(string)
				}
				if !reflect.DeepEqual(newSnapshot.Tags, tmpTags) {
					datum["tags"] = newSnapshot.Tags
				}
			}
			currentTime := time.Now().UTC()
			datum["updated"] = currentTime.Format("2006-01-02T15:04:05Z")
			queryData, _ := json.Marshal(&datum)
			queryString := fmt.Sprintf("put record %v.snapshots %v", config.Config.DBName, string(queryData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
			return nil
		}
	}
	err = RegisterSnapshot(newSnapshot)
	return err
}

func GetSnapshots(limit, filter, count string) ([]map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.snapshots", config.Config.DBName)
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
	return data, nil
}
