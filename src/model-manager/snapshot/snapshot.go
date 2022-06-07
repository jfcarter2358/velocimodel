package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"model-manager/config"
	"model-manager/model"
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
	Language string                 `json:"language"`
	Model    string                 `json:"model"`
	Releases []string               `json:"releases"`
	Version  int                    `json:"version"`
}

const GIT_TYPE = "git"
const RAW_TYPE = "raw"
const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"
const ORDERASC_DEFAULT = "NA"
const ORDERDSC_DEFAULT = "NA"

var allowedTypes = []string{GIT_TYPE, RAW_TYPE}

func CreateSnapshotFromModel(input model.Model) (string, error) {
	countObj, err := connection.Query(fmt.Sprintf("get record %v.snapshots | filter model = \"%v\" | count", config.Config.DB.Name, input.ID))
	if err != nil {
		return "", err
	}
	startCount := int(countObj[0]["count"].(float64))
	currentTime := time.Now().UTC()
	snapshotID := uuid.New().String()
	newSnapshot := Snapshot{
		ID:       snapshotID,
		Name:     input.Name,
		Created:  currentTime.Format("2006-01-02T15:04:05Z"),
		Updated:  currentTime.Format("2006-01-02T15:04:05Z"),
		Type:     input.Type,
		Metadata: input.Metadata,
		Assets:   input.Assets,
		Language: input.Language,
		Tags:     input.Tags,
		Model:    input.ID,
		Releases: make([]string, 0),
		Version:  startCount + 1,
	}
	queryData, _ := json.Marshal(&newSnapshot)
	queryString := fmt.Sprintf("post record %v.snapshots %v", config.Config.DB.Name, string(queryData))
	_, err = connection.Query(queryString)
	if err != nil {
		return "", err
	}

	modelObjects, err := connection.Query(fmt.Sprintf("get record %v.models | filter id = \"%v\"", config.Config.DB.Name, input.ID))
	if err != nil {
		return "", err
	}
	if len(modelObjects) == 0 {
		return "", errors.New(fmt.Sprintf("Snapshot not found with ID %v", input.ID))
	}
	realModelID := modelObjects[0][".id"].(string)
	snapshotList := modelObjects[0]["snapshots"].([]interface{})
	snapshotList = append(snapshotList, snapshotID)
	snapshotBytes, _ := json.Marshal(snapshotList)
	_, err = connection.Query(fmt.Sprintf("patch record %v.models \"%v\" {\"snapshots\":%v}", config.Config.DB.Name, realModelID, string(snapshotBytes)))
	if err != nil {
		return "", err
	}
	return snapshotID, nil
}

func RegisterSnapshot(newSnapshot Snapshot) error {
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
	queryString := fmt.Sprintf("post record %v.snapshots %v", config.Config.DB.Name, string(queryData))
	_, err := connection.Query(queryString)
	return err
}

func DeleteSnapshot(snapshotIDs []string) error {
	queryString := fmt.Sprintf("get record %v.snapshots", config.Config.DB.Name)
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
	queryString = fmt.Sprintf("delete record %v.snapshots %v", config.Config.DB.Name, string(queryData))
	_, err = connection.Query(queryString)
	return err
}

func UpdateSnapshot(newSnapshot Snapshot) error {
	if newSnapshot.ID == "" {
		err := errors.New("'id' field is required to update an snapshot")
		return err
	}
	queryString := fmt.Sprintf("get record %v.snapshots", config.Config.DB.Name)
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
			queryString := fmt.Sprintf("put record %v.snapshots %v", config.Config.DB.Name, string(queryData))
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

func GetSnapshots(limit, filter, count, orderasc, orderdsc string) ([]Snapshot, error) {
	queryString := fmt.Sprintf("get record %v.snapshots", config.Config.DB.Name)
	if filter != FILTER_DEFAULT {
		queryString += fmt.Sprintf(" | filter %v", filter)
	}
	if limit != LIMIT_DEFAULT {
		queryString += fmt.Sprintf(" | limit %v", limit)
	}
	if count != COUNT_DEFAULT {
		queryString += " | count"
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
	marshalled, _ := json.Marshal(data)
	var output []Snapshot
	json.Unmarshal(marshalled, &output)
	return output, nil
}
