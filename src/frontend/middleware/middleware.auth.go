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
			c.Redirect(307, "/auth/login")
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
			c.Redirect(307, "/auth/login")
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Token is invalid: %v", err)
			c.Redirect(307, "/auth/login")
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
		if err == nil {
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
		if err == nil {
			log.Printf("Error on checking userinfo: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Token is invalid: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var obj map[string]interface{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Encountered error reading body: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error parsing JSON: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		groupString := obj["groups"].(string)
		groupList := strings.Split(groupString, ",")
		if !StringSliceContains(groupList, group) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func EnsureRoleAllowed(role string) gin.HandlerFunc {
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
		if err == nil {
			log.Printf("Error on checking userinfo: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("Token is invalid: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var obj map[string]interface{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Encountered error reading body: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		err = json.Unmarshal([]byte(body), &obj)
		if err != nil {
			log.Printf("Encountered error parsing JSON: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		roleString := obj["roles"].(string)
		roleList := strings.Split(roleString, ",")
		if !StringSliceContains(roleList, role) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// This middleware sets whether the user is logged in or not
// func SetUserStatus() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if token, err1 := c.Cookie("token"); err1 == nil || token != "" {
// 			if userId, err2 := c.Cookie("userId"); err2 == nil || userId != "" {
// 				c.Set("is_logged_in", true)
// 				c.Set("user", userId)
// 				c.Set("group_list", user.GetGroupsForID(userId))
// 				c.Set("role_list", user.GetRolesForID(userId))
// 			} else {
// 				c.Set("is_logged_in", false)
// 				c.Set("user", "0")
// 				c.Set("group_list", "")
// 				c.Set("role_list", "")
// 			}
// 		} else {
// 			c.Set("is_logged_in", false)
// 			c.Set("user", "0")
// 			c.Set("group_list", "")
// 			c.Set("role_list", "")
// 		}
// 	}
// }

func StringSliceContains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
