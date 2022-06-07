package api

import (
	"api-server/config"
	"api-server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"
const SERVICE_CHECK_DELAY = 5

var Healthy = false
var Services = make([]map[string]interface{}, 0)

// Helper functions

func GetServicesLoop() {
	queryURL := fmt.Sprintf("%v/api/service", config.Config.ServiceManagerURL)
	for {
		resp, err := http.Get(queryURL)
		if err != nil {
			log.Printf("Encountered error: %v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Encountered error: %v", err)
		}
		err = json.Unmarshal([]byte(body), &Services)
		if err != nil {
			log.Printf("Encountered error: %v", err)
		}
		time.Sleep(SERVICE_CHECK_DELAY * time.Second)
	}
}

func sendDelete(serviceName, objectType, path string, queryParams map[string][]string, data interface{}, headers http.Header) error {
	if data == nil {
		data = make(map[string]interface{})
	}
	// try each service of the correct type we want to talk to
	for _, service := range Services {
		if service["type"].(string) != serviceName {
			continue
		}
		requestURL := fmt.Sprintf("http://%v:%v/api/%v", service["host"].(string), int(service["port"].(float64)), objectType)
		if path != "" {
			requestURL += fmt.Sprintf("/%v", path)
		}
		if len(queryParams) > 0 {
			requestURL += "?"
			params := url.Values{}
			for key, val := range queryParams {
				params.Add(key, val[0])
			}
			requestURL += params.Encode()
		}
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
			return fmt.Errorf("request failed with status code %v", resp.StatusCode)
		}
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("no services of type %v responded with status 200", serviceName)
}

func sendGet(serviceName, objectType, path string, queryParams map[string][]string, headers http.Header) ([]map[string]interface{}, error) {
	// try each service of the correct type we want to talk to
	for _, service := range Services {
		if service["type"].(string) != serviceName {
			continue
		}
		var obj []map[string]interface{}
		requestURL := fmt.Sprintf("http://%v:%v/api/%v", service["host"].(string), int(service["port"].(float64)), objectType)
		if path != "" {
			requestURL += fmt.Sprintf("/%v", path)
		}
		if len(queryParams) > 0 {
			requestURL += "?"
			params := url.Values{}
			for key, val := range queryParams {
				params.Add(key, val[0])
			}
			requestURL += params.Encode()
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		req.Header = headers
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed with status code %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			tmpObj := make(map[string]interface{})
			err = json.Unmarshal([]byte(body), &tmpObj)
			if err != nil {
				log.Printf("Encountered error: %v", err)
				continue
			}
			log.Printf("Recovered from error, unmarshalled JSON into map")
			obj = []map[string]interface{}{tmpObj}
		}
		return obj, nil
	}
	return nil, fmt.Errorf("no services of type %v responded with status 200", serviceName)
}

func sendGetNoAPI(serviceName, objectType, path string, queryParams map[string][]string, headers http.Header) ([]map[string]interface{}, error) {
	// try each service of the correct type we want to talk to
	for _, service := range Services {
		if service["type"].(string) != serviceName {
			continue
		}
		var obj []map[string]interface{}
		requestURL := fmt.Sprintf("http://%v:%v/%v", service["host"].(string), int(service["port"].(float64)), objectType)
		if path != "" {
			requestURL += fmt.Sprintf("/%v", path)
		}
		if len(queryParams) > 0 {
			requestURL += "?"
			params := url.Values{}
			for key, val := range queryParams {
				params.Add(key, val[0])
			}
			requestURL += params.Encode()
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		req.Header = headers
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed with status code %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			tmpObj := make(map[string]interface{})
			err = json.Unmarshal([]byte(body), &tmpObj)
			if err != nil {
				log.Printf("Encountered error: %v", err)
				continue
			}
			log.Printf("Recovered from error, unmarshalled JSON into map")
			obj = []map[string]interface{}{tmpObj}
		}
		return obj, nil
	}
	return nil, fmt.Errorf("no services of type %v responded with status 200", serviceName)
}

func sendGetRaw(serviceName, objectType, path string, queryParams map[string][]string, headers http.Header) (*http.Response, error) {
	// try each service of the correct type we want to talk to
	for _, service := range Services {
		if service["type"].(string) != serviceName {
			continue
		}
		requestURL := fmt.Sprintf("http://%v:%v/api/%v", service["host"].(string), int(service["port"].(float64)), objectType)
		if path != "" {
			requestURL += fmt.Sprintf("/%v", path)
		}
		if len(queryParams) > 0 {
			requestURL += "?"
			params := url.Values{}
			for key, val := range queryParams {
				params.Add(key, val[0])
			}
			requestURL += params.Encode()
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		req.Header = headers
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("no services of type %v responded with status 200", serviceName)
}

func sendPost(serviceName, objectType, path string, queryParams map[string][]string, data map[string]interface{}, headers http.Header) (map[string]interface{}, error) {
	if data == nil {
		data = make(map[string]interface{})
	}
	// try each service of the correct type we want to talk to
	for _, service := range Services {
		if service["type"].(string) != serviceName {
			continue
		}
		var obj map[string]interface{}
		requestURL := fmt.Sprintf("http://%v:%v/api/%v", service["host"].(string), int(service["port"].(float64)), objectType)
		if path != "" {
			requestURL += fmt.Sprintf("/%v", path)
		}
		if len(queryParams) > 0 {
			requestURL += "?"
			params := url.Values{}
			for key, val := range queryParams {
				params.Add(key, val[0])
			}
			requestURL += params.Encode()
		}
		json_data, err := json.Marshal(data)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(json_data))
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		req.Header = headers
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed with status code %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		return obj, nil
	}
	return nil, fmt.Errorf("no services of type %v responded with status 200", serviceName)
}

func sendPostFile(serviceName, objectType, path string, c *gin.Context, headers http.Header) (map[string]interface{}, error) {
	// try each service of the correct type we want to talk to
	for _, service := range Services {
		if service["type"].(string) != serviceName {
			continue
		}
		var obj map[string]interface{}
		requestURL := fmt.Sprintf("http://%v:%v/api/%v", service["host"].(string), int(service["port"].(float64)), objectType)
		if path != "" {
			requestURL += fmt.Sprintf("/%v", path)
		}

		client := &http.Client{}

		b, w, err := createMultipartFormData(c)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}

		req, err := http.NewRequest("POST", requestURL, &b)
		if err != nil {
			return nil, err
		}
		req.Header = headers
		// Don't forget to set the content type, this will contain the boundary.
		req.Header.Set("Content-Type", w.FormDataContentType())

		resp, err := client.Do(req)

		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed with status code %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		return obj, nil
	}
	return nil, fmt.Errorf("no services of type %v responded with status 200", serviceName)
}

func sendPut(serviceName, objectType, path string, queryParams map[string][]string, data map[string]interface{}, headers http.Header) error {
	if data == nil {
		data = make(map[string]interface{})
	}
	// try each service of the correct type we want to talk to
	for _, service := range Services {
		if service["type"].(string) != serviceName {
			continue
		}
		requestURL := fmt.Sprintf("http://%v:%v/api/%v", service["host"].(string), int(service["port"].(float64)), objectType)
		if path != "" {
			requestURL += fmt.Sprintf("/%v", path)
		}
		if len(queryParams) > 0 {
			requestURL += "?"
			params := url.Values{}
			for key, val := range queryParams {
				params.Add(key, val[0])
			}
			requestURL += params.Encode()
		}
		json_data, err := json.Marshal(data)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return err
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(json_data))
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return err
		}
		req.Header = headers
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("request failed with status code %v", resp.StatusCode)
		}
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("no services of type %v responded with status 200", serviceName)
}

func createMultipartFormData(c *gin.Context) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return bytes.Buffer{}, nil, err
	}
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
	if !Healthy {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Status(http.StatusOK)
}

func GetStatuses(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGetNoAPI("service-manager", "status", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

// Asset API

func DeleteAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("asset-manager", "asset", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetAssets(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("asset-manager", "asset", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("asset-manager", "asset", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PutAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendPut("asset-manager", "asset", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func CreateFileAsset(c *gin.Context) {
	data, err := sendPostFile("asset-manager", "asset", "file", c, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func CreateGitAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("asset-manager", "asset", "git", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func SyncGitAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("asset-manager", "asset", "git/sync", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func DeleteModel(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("model-manager", "model", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteModelAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("model-manager", "model", "asset", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func DownloadFileAsset(c *gin.Context) {
	assetID := c.Param("id")
	queryParams := c.Request.URL.Query()

	response, err := sendGetRaw("asset-manager", "asset", "file/"+assetID, queryParams, c.Request.Header)
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

func GetModels(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("model-manager", "model", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostModel(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("model-manager", "model", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostModelAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("model-manager", "model", "asset", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PutModel(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendPut("model-manager", "model", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func DownloadModel(c *gin.Context) {
	modelID := c.Param("id")
	queryParams := c.Request.URL.Query()

	response, err := sendGetRaw("model-manager", "model", "archive/"+modelID, queryParams, c.Request.Header)
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

func DeleteParam(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("service-manager", "param", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetParams(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("service-manager", "param", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostParam(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("service-manager", "param", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PutParam(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendPut("service-manager", "param", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// func DeleteRelease(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()
// 	var input []string
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	err := sendDelete("model-manager", "release", "", queryParams, input)
// 	if err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	c.Status(http.StatusOK)
// }

func GetReleases(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("model-manager", "release", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostRelease(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("model-manager", "release", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostReleaseSnapshot(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("model-manager", "release", "snapshot", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func DownloadRelease(c *gin.Context) {
	releaseID := c.Param("id")
	queryParams := c.Request.URL.Query()

	response, err := sendGetRaw("model-manager", "release", "archive/"+releaseID, queryParams, c.Request.Header)
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

// func PutRelease(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()
// 	var input map[string]interface{}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	err := sendPut("model-manager", "release", "", queryParams, input)
// 	if err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	c.Status(http.StatusOK)
// }

func DeleteSecret(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("service-manager", "secret", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetSecrets(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("service-manager", "secret", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostSecret(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("service-manager", "secret", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PutSecret(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendPut("service-manager", "secret", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteService(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("service-manager", "service", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetServices(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("service-manager", "service", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostService(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("service-manager", "service", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PutService(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendPut("service-manager", "service", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// func DeleteSnapshot(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()
// 	var input []string
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	err := sendDelete("model-manager", "snapshot", "", queryParams, input)
// 	if err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	c.Status(http.StatusOK)
// }

func GetSnapshots(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("model-manager", "snapshot", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostSnapshot(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("model-manager", "snapshot", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostSnapshotModel(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("model-manager", "snapshot", "model", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func DownloadSnapshot(c *gin.Context) {
	snapshotID := c.Param("id")
	queryParams := c.Request.URL.Query()

	response, err := sendGetRaw("model-manager", "snapshot", "archive/"+snapshotID, queryParams, c.Request.Header)
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

// func PutSnapshot(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()
// 	var input map[string]interface{}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	err := sendPut("model-manager", "snapshot", "", queryParams, input)
// 	if err != nil {
// 		utils.Error(err, c, http.StatusInternalServerError)
// 		return
// 	}
// 	c.Status(http.StatusOK)
// }

func DeleteUser(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("auth-manager", "user", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("auth-manager", "user", "", queryParams, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostUser(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	data, err := sendPost("auth-manager", "user", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PutUser(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendPut("auth-manager", "user", "", queryParams, input, c.Request.Header)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
