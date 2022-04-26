package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/action"
	"frontend/config"
	"frontend/utils"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

var Healthy = false

// Helper functions

func createMultipartFormData(c *gin.Context) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	fileHeader, err := c.FormFile("file")
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return bytes.Buffer{}, nil, err
	}
	if fw, err = w.CreateFormFile("file", fileHeader.Filename); err != nil {
		log.Printf("Encountered error: %v", err)
		return bytes.Buffer{}, nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Printf("Encountered error: %v", err)
		return bytes.Buffer{}, nil, err
	}
	w.Close()
	return b, w, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

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

func CreateFileAsset(c *gin.Context) {
	var obj map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/asset/file", config.Config.APIServerURL)
	client := &http.Client{}

	b, w, err := createMultipartFormData(c)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", requestURL, &b)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
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

func CreateGitAsset(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	var obj map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/asset/git", config.Config.APIServerURL)
	json_data, err := json.Marshal(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
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

func ModelAddAsset(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	var obj map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/model/asset", config.Config.APIServerURL)
	json_data, err := json.Marshal(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
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

func ModelDeleteAsset(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	requestURL := fmt.Sprintf("%v/api/model/asset", config.Config.APIServerURL)
	json_data, err := json.Marshal(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func DownloadAsset(c *gin.Context) {
	assetID := c.Param("id")

	requestURL := fmt.Sprintf("%v/api/asset/file/%v", config.Config.APIServerURL, assetID)
	response, err := http.Get(requestURL)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
	}

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	reader := response.Body
	defer reader.Close()
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := make(map[string]string)
	for name, values := range response.Header {
		if name == "Content-Type" {
			continue
		}
		extraHeaders[name] = values[0]
	}

	delete(extraHeaders, "Content-Type")

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func SyncGitAsset(c *gin.Context) {
	assetID := c.Param("id")
	asset, err := action.GetAssetByID(assetID)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	requestURL := fmt.Sprintf("%v/api/asset/git/sync", config.Config.APIServerURL)
	json_data, err := json.Marshal(asset)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(resp.StatusCode, gin.H{})
}

func CreateNewModel(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	var obj map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/model", config.Config.APIServerURL)
	json_data, err := json.Marshal(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
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

func DownloadModel(c *gin.Context) {
	modelID := c.Param("id")

	requestURL := fmt.Sprintf("%v/api/model/archive/%v", config.Config.APIServerURL, modelID)
	response, err := http.Get(requestURL)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
	}

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	reader := response.Body
	defer reader.Close()
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := make(map[string]string)
	for name, values := range response.Header {
		if name == "Content-Type" {
			continue
		}
		extraHeaders[name] = values[0]
	}

	delete(extraHeaders, "Content-Type")

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
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

	var obj map[string]interface{}

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
	c.JSON(resp.StatusCode, obj)
}

func DownloadRelease(c *gin.Context) {
	releaseID := c.Param("id")

	requestURL := fmt.Sprintf("%v/api/release/archive/%v", config.Config.APIServerURL, releaseID)
	response, err := http.Get(requestURL)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
	}

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	reader := response.Body
	defer reader.Close()
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := make(map[string]string)
	for name, values := range response.Header {
		if name == "Content-Type" {
			continue
		}
		extraHeaders[name] = values[0]
	}

	delete(extraHeaders, "Content-Type")

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
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
	var obj map[string]interface{}

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
	c.JSON(resp.StatusCode, obj)
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

func DownloadSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")

	requestURL := fmt.Sprintf("%v/api/snapshot/archive/%v", config.Config.APIServerURL, snapshotID)
	response, err := http.Get(requestURL)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
	}

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	reader := response.Body
	defer reader.Close()
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := make(map[string]string)
	for name, values := range response.Header {
		if name == "Content-Type" {
			continue
		}
		extraHeaders[name] = values[0]
	}

	delete(extraHeaders, "Content-Type")

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
