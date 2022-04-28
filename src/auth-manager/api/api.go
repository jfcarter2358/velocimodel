package api

import (
	"auth-manager/generates"
	"auth-manager/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var PrivKey = []byte("velocimodelsign")

// get all users
func UserGetAllHandler(c *gin.Context) {
	users := user.GetAllUsers()
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// get user by id
func UserGetByIdHandler(c *gin.Context) {
	userID := c.Param("id")
	userData, err := user.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"user": userData})
}

// create user
func UserCreateHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	givenName := c.PostForm("given_name")
	familyName := c.PostForm("family_name")
	email := c.PostForm("email")
	roles := c.PostForm("roles")
	groups := c.PostForm("groups")
	_, err := user.CreateNewUser(givenName, familyName, username, password, email, strings.Split(roles, ","), strings.Split(groups, ","))
	if err != nil {
		log.Println("Unable to create user")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Redirect(302, "/ui/edit")
}

// update user
func UserUpdateHandler(c *gin.Context) {
	userID := c.Param("id")
	username := c.PostForm("username")
	password := c.PostForm("password")
	givenName := c.PostForm("given_name")
	familyName := c.PostForm("family_name")
	email := c.PostForm("email")
	roles := c.PostForm("roles")
	groups := c.PostForm("groups")

	usr, err := user.GetUserByID(userID)

	if password != usr.Password {
		usr.Password = HashAndSalt([]byte(password))
	}

	usr.Username = username
	usr.GivenName = givenName
	usr.FamilyName = familyName
	usr.Email = email
	usr.Roles = roles
	usr.Groups = groups

	err = user.UpdateUserContents(userID, *usr)
	if err != nil {
		log.Println("Unable to update user")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Redirect(302, "/ui/edit")
}

// delete user
func UserDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := user.DeleteUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// helper function
func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func UserInfo(c *gin.Context) {
	if len(c.Request.Header["Authorization"]) > 0 {
		if val, ok := c.Request.Header["Authorization"]; ok {
			authToken := strings.Split(val[0], " ")[1]

			// Parse and verify jwt access token
			token, err := jwt.ParseWithClaims(authToken, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("parse error")
				}
				return PrivKey, nil
			})
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(500)
				return
			}

			claims, ok := token.Claims.(*generates.JWTAccessClaims)
			if !ok || !token.Valid {
				log.Println("invalid token")
				c.AbortWithStatus(500)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"given_name":  claims.GivenName,
				"family_name": claims.FamilyName,
				"email":       claims.Email,
				"groups":      claims.Groups,
				"roles":       claims.Roles,
				"exp":         claims.ExpiresAt,
				"sub":         claims.Subject,
				"aud":         claims.Audience,
			})
			return
		}
	}
	c.AbortWithStatus(500)
}
