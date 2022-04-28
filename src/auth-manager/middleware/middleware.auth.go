// middleware.auth.go

package middleware

import (
	"auth-manager/user"
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
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			c.Redirect(307, "/u/login")
		}
	}
}

func EnsureLoggedInAbort() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// This middleware ensures that a request will be aborted with an error
// if the user is already logged in
func EnsureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's no error or if the token is not empty
		// the user is already logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn {
			// if token, err := c.Cookie("token"); err == nil || token != "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func EnsureGroupAllowed(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			c.Redirect(307, "/u/login")
		}
		groupString, _ := c.Get("group_list")
		log.Printf("GROUP LIST: %v", groupString)
		groupList := strings.Split(groupString.(string), ",")
		if !StringSliceContains(groupList, group) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func EnsureRoleAllowed(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			c.Redirect(307, "/u/login")
		}
		roleString, _ := c.Get("role_list")
		roleList := strings.Split(roleString.(string), ",")
		if !StringSliceContains(roleList, role) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// This middleware sets whether the user is logged in or not
func SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err1 := c.Cookie("token"); err1 == nil || token != "" {
			if userId, err2 := c.Cookie("userId"); err2 == nil || userId != "" {
				c.Set("is_logged_in", true)
				c.Set("user", userId)
				c.Set("group_list", user.GetGroupsForID(userId))
				c.Set("role_list", user.GetRolesForID(userId))
			} else {
				c.Set("is_logged_in", false)
				c.Set("user", "0")
				c.Set("group_list", "")
				c.Set("role_list", "")
			}
		} else {
			c.Set("is_logged_in", false)
			c.Set("user", "0")
			c.Set("group_list", "")
			c.Set("role_list", "")
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
