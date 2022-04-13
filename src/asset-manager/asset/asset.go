package asset

import (
	"asset-manager/config"
	"asset-manager/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

type Asset struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	URL      string                 `json:"url" binding:"required"`
	Created  string                 `json:"created"`
	Updated  string                 `json:"updated"`
	Type     string                 `json:"type" binding:"required"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
}

var Assets = make([]Asset, 0)

var FILE_TYPE = "file"
var GIT_TYPE = "git"
var S3_TYPE = "s3"
var allowedTypes = []string{FILE_TYPE, GIT_TYPE, S3_TYPE}

func RegisterAsset(newAsset Asset) error {
	if newAsset.ID == "" {
		newAsset.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newAsset.Type) {
		err := errors.New(fmt.Sprintf("Assset type of %v does not exist", newAsset.Type))
		return err
	}
	Assets = append(Assets, newAsset)
	log.Println(fmt.Sprintf("%v", Assets))
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
	for idx, svc := range Assets {
		if svc.ID == newAsset.ID {
			if !utils.Contains(allowedTypes, newAsset.Type) {
				err := errors.New(fmt.Sprintf("Assset type of %v does not exist", newAsset.Type))
				return err
			}
			Assets[idx] = newAsset
		}
	}
	return nil
}

func SyncAssets() error {
	for _, ast := range Assets {
		queryString := fmt.Sprintf("get record %v.assets | filter id = \"%v\"", config.Config.DBName, ast.ID)
		data, err := connection.Query(queryString)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			log.Println("POST ASSET")
			assetData, _ := json.Marshal(&ast)
			queryString := fmt.Sprintf("post record %v.assets %v", config.Config.DBName, string(assetData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
		} else {
			log.Println("PUT ASSET")
			tempData := map[string]interface{}{
				".id":      data[0][".id"].(string),
				"id":       ast.ID,
				"name":     ast.Name,
				"url":      ast.URL,
				"created":  ast.Created,
				"updated":  ast.Updated,
				"type":     ast.Type,
				"tags":     ast.Tags,
				"metadata": ast.Metadata,
			}
			assetData, _ := json.Marshal(&tempData)
			queryString := fmt.Sprintf("put record %v.assets %v", config.Config.DBName, string(assetData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func LoadAssets() error {
	queryString := fmt.Sprintf("get record %v.assets", config.Config.DBName)
	data, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	Assets = make([]Asset, 0)
	for _, datum := range data {
		interfaceTags := datum["tags"].([]interface{})
		stringTags := make([]string, len(interfaceTags))
		for idx, val := range interfaceTags {
			stringTags[idx] = val.(string)
		}
		newService := Asset{
			ID:       datum["id"].(string),
			Name:     datum["name"].(string),
			URL:      datum["url"].(string),
			Created:  datum["created"].(string),
			Updated:  datum["updated"].(string),
			Type:     datum["type"].(string),
			Tags:     stringTags,
			Metadata: datum["metadata"].(map[string]interface{}),
		}
		Assets = append(Assets, newService)
	}
	return nil
}
