// middleware.auth.go

package middleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
		token, err := c.Cookie("access_token")
		if err != nil {
			c.Redirect(307, "/ui/401")
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error on checking userinfo: %v", err)
			c.Redirect(307, "/ui/401")
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Token is invalid: %v", err)
			c.Redirect(307, "/ui/401")
			return
		}
	}
}

func EnsureLoggedInAbort() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		token, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error on checking userinfo: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Token is invalid: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

// This middleware ensures that a request will be aborted with an error
// if the user is already logged in
func EnsureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		token, err := c.Cookie("access_token")
		if err != nil {
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			return
		}
		if res.StatusCode != http.StatusOK {
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func EnsureGroupAllowed(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		token, err := c.Cookie("access_token")
		if err != nil {
			c.Redirect(307, "/ui/401")
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error on checking userinfo: %v", err)
			c.Redirect(307, "/ui/401")
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Token is invalid: %v", err)
			c.Redirect(307, "/ui/401")
			return
		}
		var obj map[string]interface{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Encountered error reading body: %v", err)
			c.Redirect(307, "/ui/401")
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error parsing JSON: %v", err)
			c.Redirect(307, "/ui/401")
		}
		groupString := obj["groups"].(string)
		log.Printf("GROUP STRING: %v", groupString)
		groupList := strings.Split(groupString, ",")
		if !StringSliceContains(groupList, group) {
			c.Redirect(307, "/ui/401")
		}
	}
}

func EnsureRoleAllowed(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		token, err := c.Cookie("access_token")
		if err != nil {
			c.Redirect(307, "/ui/401")
			return
		}
		client := http.Client{}
		requestURL := "http://auth-manager:9005/oauth/userinfo"
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error on checking userinfo: %v", err)
			c.Redirect(307, "/ui/401")
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Token is invalid: %v", err)
			c.Redirect(307, "/ui/401")
			return
		}
		var obj map[string]interface{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Encountered error reading body: %v", err)
			c.Redirect(307, "/ui/401")
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error parsing JSON: %v", err)
			c.Redirect(307, "/ui/401")
		}
		roleString := obj["roles"].(string)
		roleList := strings.Split(roleString, ",")
		if !StringSliceContains(roleList, role) {
			c.Redirect(307, "/ui/401")
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
