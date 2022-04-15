package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"service-manager/config"
	"service-manager/utils"

	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

type Service struct {
	ID   string `json:"id"`
	Host string `json:"host" binding:"required"`
	Port int    `json:"port" binding:"required"`
	Type string `json:"type" binding:"required"`
}

const API_SERVER_TYPE = "api-server"
const ASSET_MANAGER_TYPE = "asset-manager"
const FRONTEND_TYPE = "frontend"
const MODEL_MANAGER_TYPE = "model-manager"
const SERVICE_MANAGER_TYPE = "service-manager"
const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"

var allowedTypes = []string{API_SERVER_TYPE, ASSET_MANAGER_TYPE, FRONTEND_TYPE, MODEL_MANAGER_TYPE, SERVICE_MANAGER_TYPE}

func RegisterService(newService Service) error {
	if newService.ID == "" {
		newService.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newService.Type) {
		err := errors.New(fmt.Sprintf("Service type of %v does not exist", newService.Type))
		return err
	}
	queryData, _ := json.Marshal(&newService)
	queryString := fmt.Sprintf("post record %v.services %v", config.Config.DBName, string(queryData))
	_, err := connection.Query(queryString)
	return err
}

func DeleteService(serviceIDs []string) error {
	queryString := fmt.Sprintf("get record %v.services", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(serviceIDs, datum["id"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.services %v", config.Config.DBName, string(queryData))
	_, err = connection.Query(queryString)
	return err
}

func UpdateService(newService Service) error {
	if newService.ID == "" {
		err := errors.New("'id' field is required to update a service")
		return err
	}
	if !utils.Contains(allowedTypes, newService.Type) {
		err := errors.New(fmt.Sprintf("Service type of %v does not exist", newService.Type))
		return err
	}
	queryData, _ := json.Marshal(&newService)
	queryString := fmt.Sprintf("post record %v.assets %v", config.Config.DBName, string(queryData))
	_, err := connection.Query(queryString)
	return err
}

func GetServices(limit, filter, count string) ([]map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.services", config.Config.DBName)
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
