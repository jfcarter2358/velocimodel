// middleware.auth.go

package middleware

import (
	"api-server/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// "log"
)

// This middleware ensures that a request will be aborted with an error
// if the user is not logged in
func EnsureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		var token string
		authString := c.Request.Header.Get("Authorization")
		if authString == "" {
			tempToken, err := c.Cookie("access_token")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			token = tempToken
		} else {
			token = strings.Split(authString, " ")[1]
		}
		if token == config.Config.JoinToken {
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if res.StatusCode != http.StatusOK {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func EnsureGroupAllowed(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		var token string
		authString := c.Request.Header.Get("Authorization")
		if authString == "" {
			tempToken, err := c.Cookie("access_token")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			token = tempToken
		} else {
			token = strings.Split(authString, " ")[1]
		}
		if token == config.Config.JoinToken {
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if res.StatusCode != http.StatusOK {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var obj map[string]interface{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		groups := make([]string, len(obj["groups"].([]interface{})))
		for idx, val := range obj["groups"].([]interface{}) {
			groups[idx] = val.(string)
		}
		if !StringSliceContains(groups, group) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func EnsureRoleAllowed(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		var token string
		authString := c.Request.Header.Get("Authorization")
		if authString == "" {
			tempToken, err := c.Cookie("access_token")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			token = tempToken
		} else {
			token = strings.Split(authString, " ")[1]
		}
		if token == config.Config.JoinToken {
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if res.StatusCode != http.StatusOK {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var obj map[string]interface{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		roles := make([]string, len(obj["roles"].([]interface{})))
		for idx, val := range obj["roles"].([]interface{}) {
			roles[idx] = val.(string)
		}
		if !StringSliceContains(roles, role) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func StringSliceContains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
