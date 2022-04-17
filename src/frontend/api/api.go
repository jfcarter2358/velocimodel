package api

import (
	"encoding/json"
	"fmt"
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
