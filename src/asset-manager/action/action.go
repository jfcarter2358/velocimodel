package action

import (
	"asset-manager/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func SendDelete(path string, data interface{}) error {
	if data == nil {
		data = make(map[string]interface{})
	}
	requestURL := fmt.Sprintf("%v/%v", config.Config.APIServerURL, path)
	json_data, err := json.Marshal(data)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Request failed with status code %v", resp.StatusCode))
	}
	return nil
}

func SendPost(path string, data map[string]interface{}) (map[string]interface{}, error) {
	if data == nil {
		data = make(map[string]interface{})
	}
	var obj map[string]interface{}
	requestURL := fmt.Sprintf("%v/%v", config.Config.APIServerURL, path)
	json_data, err := json.Marshal(data)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Request failed with status code %v", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	return obj, nil
}

func RecurseAddTree(tree map[string]interface{}, keys []string, fileType string) map[string]interface{} {
	if len(keys) == 1 {
		if fileType == "file" {
			tree[keys[0]] = "file"
		} else {
			if tree[keys[0]] == nil {
				tree[keys[0]] = make(map[string]interface{})
			}
		}
	} else {
		if tree[keys[0]] != nil {
			tree[keys[0]] = RecurseAddTree(tree[keys[0]].(map[string]interface{}), keys[1:], fileType)
		} else {
			tree[keys[0]] = RecurseAddTree(make(map[string]interface{}), keys[1:], fileType)
		}
	}
	return tree
}
