package release

import (
	"encoding/json"
	"errors"
	"fmt"
	"model-manager/config"
	"model-manager/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

type Release struct {
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

func RegisterRelease(newRelease Release) error {
	countObj, _ := connection.Query(fmt.Sprintf("get record %v.releases | count", config.Config.DBName))
	startCount := int(countObj[0]["count"].(float64))
	endCount := startCount
	if newRelease.ID == "" {
		newRelease.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newRelease.Type) {
		err := errors.New(fmt.Sprintf("Release type of %v does not exist", newRelease.Type))
		return err
	}
	currentTime := time.Now().UTC()
	newRelease.Created = currentTime.Format("2006-01-02T15:04:05Z")
	newRelease.Updated = currentTime.Format("2006-01-02T15:04:05Z")
	queryData, _ := json.Marshal(&newRelease)
	for endCount == startCount {
		queryString := fmt.Sprintf("post record %v.releases %v", config.Config.DBName, string(queryData))
		_, err := connection.Query(queryString)
		if err != nil {
			return err
		}
		countObj, _ := connection.Query(fmt.Sprintf("get record %v.releases | count", config.Config.DBName))
		endCount = int(countObj[0]["count"].(float64))
	}
	return nil
}

func DeleteRelease(releaseIDs []string) error {
	queryString := fmt.Sprintf("get record %v.releases", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(releaseIDs, datum["id"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.releases %v", config.Config.DBName, string(queryData))
	_, err = connection.Query(queryString)
	return err
}

func GetReleases(limit, filter, count string) ([]map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.releases", config.Config.DBName)
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
