package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"frontend/config"
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

func sendDelete(objectType, path string, queryParams map[string][]string, data []string) error {
	requestURL := fmt.Sprintf("%v/api/%v", config.Config.APIServerURL, objectType)
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
	}
	return nil
}

func sendGet(objectType, path string, queryParams map[string][]string) ([]map[string]interface{}, error) {
	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/%v", config.Config.APIServerURL, objectType)
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
		tmpObj := make(map[string]interface{})
		err = json.Unmarshal([]byte(body), &tmpObj)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
		log.Printf("Recovered from error, unmarshalled JSON into map")
		obj = []map[string]interface{}{tmpObj}
	}
	return obj, nil
}

func sendPost(objectType, path string, queryParams map[string][]string, data map[string]interface{}) (map[string]interface{}, error) {
	var obj map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/%v", config.Config.APIServerURL, objectType)
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
			return nil, err
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error: %v", err)
			return nil, err
		}
	}
	return obj, nil
}

func sendPostFile(objectType, path string, c *gin.Context) (map[string]interface{}, error) {
	// try each service of the correct type we want to talk to
	var obj map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/%v", config.Config.APIServerURL, objectType)
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
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	return obj, nil
}

func sendPut(objectType, path string, queryParams map[string][]string, data map[string]interface{}) error {
	requestURL := fmt.Sprintf("%v/api/%v", config.Config.APIServerURL, objectType)
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
	}
	return nil
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

func DoDelete(c *gin.Context) {
	path := c.Param("path")
	queryParams := c.Request.URL.Query()
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := sendDelete(path, "", queryParams, input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func DoGet(c *gin.Context) {
	path := c.Param("path")
	queryParams := c.Request.URL.Query()
	data, err := sendGet(path, "", queryParams)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(data) > 0 {
		c.JSON(http.StatusOK, data)
	}
	c.Status(http.StatusNotFound)
}

func DoPost(c *gin.Context) {
	path := c.Param("path")
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := sendPost(path, "", queryParams, input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func DoPut(c *gin.Context) {
	path := c.Param("path")
	queryParams := c.Request.URL.Query()
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := sendPut(path, "", queryParams, input)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func DoUpload(c *gin.Context) {
	path := c.Param("path")
	data, err := sendPostFile(path, "upload", c)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, data)
}
