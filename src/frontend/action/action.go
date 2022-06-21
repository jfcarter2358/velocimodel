package action

import (
	"encoding/json"
	"fmt"
	"frontend/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func GetAssetByID(c *gin.Context, assetID string) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", assetID))
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetAssetsAll(c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/asset", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetAssetsByIDList(c *gin.Context, assetIDs []string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/asset", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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
		if Contains(assetIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetAssetsLimit(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetAssetsLimitLatest(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetModelByID(c *gin.Context, modelID string) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", modelID))
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetModelsByIDList(c *gin.Context, modelIDs []string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/model", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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
		if Contains(modelIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetModelsAll(c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/model", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetModelsLimit(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetModelsLimitLatest(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetReleaseByID(c *gin.Context, releaseID string) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", releaseID))
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetReleasesAll(c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/release", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetReleasesByIDList(c *gin.Context, releaseIDs []string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/release", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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
		if Contains(releaseIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetReleasesLimit(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetReleasesLimitLatest(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetSnapshotByID(c *gin.Context, snapshotID string) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", snapshotID))
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetSnapshotsAll(c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/snapshot", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetSnapshotsByIDList(c *gin.Context, snapshotIDs []string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/snapshot", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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
		if Contains(snapshotIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetSnapshotsLimit(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetSnapshotsLimitLatest(c *gin.Context, limit string) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetSecretsAll(c *gin.Context) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/secret", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}

	return obj[0], nil
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetUserData(c *gin.Context) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	requestURL := fmt.Sprintf("%s/oauth/userinfo", config.Config.Oauth.AuthServerInternalURL)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code was :%v, not 200", res.StatusCode)
	}
	var obj map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func GetUserByID(userID string, c *gin.Context) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", userID))
	requestURL := fmt.Sprintf("%v/api/user?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetUsersAll(c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/user", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	log.Printf("AUTH HEADER GET USERS ALL: %v", req.Header.Get("Authorization"))
	resp, err := client.Do(req)
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

func GetUsersByIDList(userIDs []string, c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/user", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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
		if Contains(userIDs, val["id"].(string)) {
			output = append(output, val)
		}
	}

	return output, nil
}

func GetUsersLimit(limit string, c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/user?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetUsersLimitLatest(limit string, c *gin.Context) ([]map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", limit)
	requestURL := fmt.Sprintf("%v/api/user?%v", config.Config.APIServerURL, params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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

func GetParamsAll(c *gin.Context) (map[string]interface{}, error) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/param", config.Config.APIServerURL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	resp, err := client.Do(req)
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
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}

	return obj[0], nil
}
