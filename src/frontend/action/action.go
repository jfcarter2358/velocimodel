package action

import (
	"encoding/json"
	"fmt"
	"frontend/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetAssetByID(assetID string) (map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", assetID))
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	if len(obj) == 0 {
		return nil, nil
	}

	return obj[0], nil
}

func GetAssetsAll() ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/asset", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetAssetsByIDList(assetIDs []string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/asset", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	output := make([]map[string]interface{}, 0)

	for _, val := range obj {
		if contains(assetIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetAssetsLimit(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetAssetsLimitLatest(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetModelByID(modelID string) (map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", modelID))
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	if len(obj) == 0 {
		return nil, nil
	}

	return obj[0], nil
}

func GetModelsByIDList(modelIDs []string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/model", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	output := make([]map[string]interface{}, 0)

	for _, val := range obj {
		if contains(modelIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetModelsAll() ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/model", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetModelsLimit(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetModelsLimitLatest(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetReleaseByID(releaseID string) (map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", releaseID))
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	if len(obj) == 0 {
		return nil, nil
	}

	return obj[0], nil
}

func GetReleasesAll() ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/release", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetReleasesByIDList(releaseIDs []string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/release", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	output := make([]map[string]interface{}, 0)

	for _, val := range obj {
		if contains(releaseIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetReleasesLimit(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetReleasesLimitLatest(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetSnapshotByID(snapshotID string) (map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", snapshotID))
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	if len(obj) == 0 {
		return nil, nil
	}

	return obj[0], nil
}

func GetSnapshotsAll() ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/snapshot", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetSnapshotsByIDList(snapshotIDs []string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/snapshot", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

	output := make([]map[string]interface{}, 0)

	for _, val := range obj {
		if contains(snapshotIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetSnapshotsLimit(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func GetSnapshotsLimitLatest(limit string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		return nil, err
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
