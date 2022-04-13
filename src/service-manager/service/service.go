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

var Services = make([]Service, 0)

var ASSET_MANAGER_TYPE = "asset-manager"
var MODEL_MANAGER_TYPE = "model-manager"
var API_SERVER_TYPE = "api-server"
var FRONTEND_TYPE = "frontend"
var allowedTypes = []string{ASSET_MANAGER_TYPE, MODEL_MANAGER_TYPE, API_SERVER_TYPE, FRONTEND_TYPE}

func RegisterService(newService Service) error {
	if newService.ID == "" {
		newService.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newService.Type) {
		err := errors.New(fmt.Sprintf("Service type of %v does not exist", newService.Type))
		return err
	}
	Services = append(Services, newService)
	return nil
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
	for idx, svc := range Services {
		if svc.ID == newService.ID {
			if !utils.Contains(allowedTypes, newService.Type) {
				err := errors.New(fmt.Sprintf("Service type of %v does not exist", newService.Type))
				return err
			}
			Services[idx] = newService
		}
	}
	return nil
}

func SyncServices() error {
	for _, svc := range Services {
		queryString := fmt.Sprintf("get record %v.services | filter id = \"%v\"", config.Config.DBName, svc.ID)
		data, err := connection.Query(queryString)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			serviceData, _ := json.Marshal(&svc)
			queryString := fmt.Sprintf("post record %v.services %v", config.Config.DBName, string(serviceData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
		} else {
			tempData := map[string]interface{}{
				".id":  data[0][".id"].(string),
				"id":   svc.ID,
				"host": svc.Host,
				"port": svc.Port,
				"type": svc.Type,
			}
			serviceData, _ := json.Marshal(&tempData)
			queryString := fmt.Sprintf("put record %v.services %v", config.Config.DBName, string(serviceData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func LoadServices() error {
	queryString := fmt.Sprintf("get record %v.services", config.Config.DBName)
	data, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	Services = make([]Service, 0)
	for _, datum := range data {
		newService := Service{
			ID:   datum["id"].(string),
			Host: datum["host"].(string),
			Port: int(datum["port"].(float64)),
			Type: datum["type"].(string),
		}
		Services = append(Services, newService)
	}
	return nil
}
