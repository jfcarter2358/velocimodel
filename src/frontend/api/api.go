package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/action"
	"frontend/config"
	"frontend/utils"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var Healthy = false

// Health API

func GetHealth(c *gin.Context) {
	if Healthy == false {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Status(http.StatusOK)
	return
}

func GetAssets(c *gin.Context) {
	var obj []map[string]interface{}

	queryParams := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%v/api/asset", config.Config.APIServerURL)
	if len(queryParams) > 0 {
		requestURL += "?"
		params := url.Values{}
		for key, val := range queryParams {
			params.Add(key, val[0])
		}
		requestURL += params.Encode()
	}
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetAsset(c *gin.Context) {
	assetID := c.Param("id")
	var obj []map[string]interface{}

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", assetID))
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	if len(obj) > 0 {
		c.JSON(http.StatusOK, obj[0])
	}
	c.JSON(http.StatusOK, obj)
}

func UpdateAsset(c *gin.Context) {
	assetID := c.Param("id")
	requestURL := fmt.Sprintf("%v/api/asset", config.Config.APIServerURL)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data["id"] = assetID
	json_data, err := json.Marshal(data)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(resp.StatusCode)
}

func GetModels(c *gin.Context) {
	var obj []map[string]interface{}

	queryParams := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%v/api/model", config.Config.APIServerURL)
	if len(queryParams) > 0 {
		requestURL += "?"
		params := url.Values{}
		for key, val := range queryParams {
			params.Add(key, val[0])
		}
		requestURL += params.Encode()
	}
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetModel(c *gin.Context) {
	modelID := c.Param("id")
	var obj []map[string]interface{}

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", modelID))
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	if len(obj) > 0 {
		c.JSON(http.StatusOK, obj[0])
	}
	c.JSON(http.StatusOK, obj)
}

func UpdateModel(c *gin.Context) {
	modelID := c.Param("id")
	requestURL := fmt.Sprintf("%v/api/model", config.Config.APIServerURL)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data["id"] = modelID
	json_data, err := json.Marshal(data)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(resp.StatusCode)
}

func GetReleases(c *gin.Context) {
	var obj []map[string]interface{}

	queryParams := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%v/api/release", config.Config.APIServerURL)
	if len(queryParams) > 0 {
		requestURL += "?"
		params := url.Values{}
		for key, val := range queryParams {
			params.Add(key, val[0])
		}
		requestURL += params.Encode()
	}
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetRelease(c *gin.Context) {
	releaseID := c.Param("id")
	var obj []map[string]interface{}

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", releaseID))
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	if len(obj) > 0 {
		c.JSON(http.StatusOK, obj[0])
	}
	c.JSON(http.StatusOK, obj)
}

func CreateRelease(c *gin.Context) {
	snapshotID := c.Param("id")

	snapshot, err := action.GetSnapshotByID(snapshotID)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	requestURL := fmt.Sprintf("%v/api/release/snapshot", config.Config.APIServerURL)

	json_data, err := json.Marshal(snapshot)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(resp.StatusCode)
}

func GetSnapshots(c *gin.Context) {
	var obj []map[string]interface{}

	queryParams := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%v/api/snapshot", config.Config.APIServerURL)
	if len(queryParams) > 0 {
		requestURL += "?"
		params := url.Values{}
		for key, val := range queryParams {
			params.Add(key, val[0])
		}
		requestURL += params.Encode()
	}
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")
	var obj []map[string]interface{}

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", snapshotID))
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		utils.Error(err, c, resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	if len(obj) > 0 {
		c.JSON(http.StatusOK, obj[0])
	}
	c.JSON(http.StatusOK, obj)
}

func CreateSnapshot(c *gin.Context) {
	modelID := c.Param("id")

	model, err := action.GetModelByID(modelID)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	requestURL := fmt.Sprintf("%v/api/snapshot/model", config.Config.APIServerURL)

	json_data, err := json.Marshal(model)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(resp.StatusCode)
}

func UpdateSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")
	requestURL := fmt.Sprintf("%v/api/snapshot", config.Config.APIServerURL)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data["id"] = snapshotID
	json_data, err := json.Marshal(data)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(resp.StatusCode)
}
