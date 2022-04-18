package api

import (
	"api-server/config"
	"api-server/utils"
	"bytes"
	"encoding/json"
	"errors"
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
	for true {
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

func sendDelete(serviceName, objectType, path string, queryParams map[string][]string, data interface{}) error {
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
		if data != nil {
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
	}
	return errors.New(fmt.Sprintf("No services of type %v responded with status 200", objectType))
}

func sendGet(serviceName, objectType, path string, queryParams map[string][]string) ([]map[string]interface{}, error) {
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
		resp, err := http.Get(requestURL)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(fmt.Sprintf("Request failed with status code %v", resp.StatusCode))
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
	return nil, errors.New(fmt.Sprintf("No services of type %v responded with status 200", objectType))
}

func sendPost(serviceName, objectType, path string, queryParams map[string][]string, data map[string]interface{}) (map[string]interface{}, error) {
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
		if data != nil {
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
				continue
			}
			err = json.Unmarshal([]byte(body), &obj)
			if err != nil {
				log.Printf("Encountered error: %v", err)
				continue
			}
			return obj, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No services of type %v responded with status 200", objectType))
}

func sendPostFile(serviceName, objectType, path string, c *gin.Context) (map[string]interface{}, error) {
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
		// Don't forget to set the content type, this will contain the boundary.
		req.Header.Set("Content-Type", w.FormDataContentType())

		resp, err := client.Do(req)

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
			continue
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			continue
		}
		return obj, nil
	}
	return nil, errors.New(fmt.Sprintf("No services of type %v responded with status 200", objectType))
}

func sendPut(serviceName, objectType, path string, queryParams map[string][]string, data map[string]interface{}) error {
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
		if data != nil {
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
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			resp, err := client.Do(req)
			if resp.StatusCode != http.StatusOK {
				return errors.New(fmt.Sprintf("Request failed with status code %v", resp.StatusCode))
			}
			return nil
		}
	}
	return errors.New(fmt.Sprintf("No services of type %v responded with status 200", objectType))
}

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

// Asset API

func DeleteAsset(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("asset-manager", "asset", "", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetAssets(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("asset-manager", "asset", "", queryParams)
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
	data, err := sendPost("asset-manager", "asset", "", queryParams, input)
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
	err := sendPut("asset-manager", "asset", "", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func UploadAsset(c *gin.Context) {
	data, err := sendPostFile("asset-manager", "asset", "upload", c)
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
	err := sendDelete("model-manager", "model", "", queryParams, input)
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
	err := sendDelete("model-manager", "model", "asset", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetModels(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("model-manager", "model", "", queryParams)
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
	data, err := sendPost("model-manager", "model", "", queryParams, input)
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
	data, err := sendPost("model-manager", "model", "asset", queryParams, input)
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
	err := sendPut("model-manager", "model", "", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteParam(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := sendDelete("service-manager", "param", "", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetParams(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("service-manager", "param", "", queryParams)
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
	data, err := sendPost("service-manager", "param", "", queryParams, input)
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
	err := sendPut("service-manager", "param", "", queryParams, input)
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
	data, err := sendGet("model-manager", "release", "", queryParams)
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
	data, err := sendPost("model-manager", "release", "", queryParams, input)
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
	data, err := sendPost("model-manager", "release", "snapshot", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
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
	err := sendDelete("service-manager", "secret", "", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetSecrets(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("service-manager", "secret", "", queryParams)
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
	data, err := sendPost("service-manager", "secret", "", queryParams, input)
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
	err := sendPut("service-manager", "secret", "", queryParams, input)
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
	err := sendDelete("service-manager", "service", "", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetServices(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	data, err := sendGet("service-manager", "service", "", queryParams)
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
	data, err := sendPost("service-manager", "service", "", queryParams, input)
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
	err := sendPut("service-manager", "service", "", queryParams, input)
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
	data, err := sendGet("model-manager", "snapshot", "", queryParams)
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
	data, err := sendPost("model-manager", "snapshot", "", queryParams, input)
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
	data, err := sendPost("model-manager", "snapshot", "model", queryParams, input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
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
