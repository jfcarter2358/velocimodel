// user.go

package user

import (
	"auth-manager/config"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	GivenName            string `json:"given_name"`
	FamilyName           string `json:"family_name"`
	ID                   string `json:"id"`
	Roles                string `json:"roles"`
	Groups               string `json:"groups"`
	Email                string `json:"email"`
	ResetToken           string `json:"reset_token"`
	ResetTokenCreateDate string `json:"reset_token_create_date"`
}

var PageLength = 10

func GetUserCount() int {
	countObj, _ := connection.Query(fmt.Sprintf("get record %v.users | count", config.Config.DBName))
	count := int(countObj[0]["count"].(float64))
	return count
}

func GetGroupsForID(id string) string {
	queryString := fmt.Sprintf("get record %v.users | filter id = \"%v\"", config.Config.DBName, id)
	data, err := connection.Query(queryString)
	if err != nil {
		return ""
	}
	if len(data) == 0 {
		return ""
	}
	groups := data[0]["groups"].(string)
	return groups
}

func GetRolesForID(id string) string {
	queryString := fmt.Sprintf("get record %v.users | filter id = \"%v\"", config.Config.DBName, id)
	data, err := connection.Query(queryString)
	if err != nil {
		return ""
	}
	if len(data) == 0 {
		return ""
	}
	roles := data[0]["roles"].(string)
	return roles
}

func GetUserByIndex(start, end int) []User {
	users := GetAllUsers()
	sort.Slice(users, func(i, j int) bool {
		return users[i].Username < users[j].Username
	})
	if end > len(users) {
		end = len(users)
	}
	return users[start:end]
}

// Return a list of all the compasses
func GetAllUsers() []User {
	var users []User
	queryString := fmt.Sprintf("get record %v.users", config.Config.DBName)
	data, _ := connection.Query(queryString)
	dataBytes, _ := json.Marshal(data)
	_ = json.Unmarshal(dataBytes, &users)
	return users
}

// Delete a user based on the ID supplied
func DeleteUserByID(id string) error {
	queryString := fmt.Sprintf("get record %v.users | filter id = \"%v\" | delete record %v.users -", config.Config.DBName, id, config.Config.DBName)
	_, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	return nil
}

// Fetch a user based on the ID supplied
func GetUserByID(id string) (*User, error) {
	var user User
	queryString := fmt.Sprintf("get record %v.users | filter id = \"%v\"", config.Config.DBName, id)
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("User not found")
	}
	dataBytes, _ := json.Marshal(data[0])
	_ = json.Unmarshal(dataBytes, &user)
	return &user, nil
}

// Fetch a user based on the email supplied
func GetUserByEmail(email string) (*User, error) {
	var user User
	queryString := fmt.Sprintf("get record %v.users | filter email = \"%v\"", config.Config.DBName, email)
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("User not found")
	}
	dataBytes, _ := json.Marshal(data[0])
	_ = json.Unmarshal(dataBytes, &user)
	return &user, nil
}

// Fetch a user based on the username supplied
func GetUserByUsername(username string) (*User, error) {
	var user User
	queryString := fmt.Sprintf("get record %v.users | filter username = \"%v\"", config.Config.DBName, username)
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("User not found")
	}
	dataBytes, _ := json.Marshal(data[0])
	_ = json.Unmarshal(dataBytes, &user)
	return &user, nil
}

// Fetch a user based on the reset token supplied
func GetUserByResetToken(resetToken string) (*User, error) {
	var user User
	queryString := fmt.Sprintf("get record %v.users | filter reset_token = \"%v\"", config.Config.DBName, resetToken)
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("User not found")
	}
	dataBytes, _ := json.Marshal(data[0])
	_ = json.Unmarshal(dataBytes, &user)
	return &user, nil
}

// Create a new user with the data provided
func CreateNewUser(givenName, familyName, username, password, email string, roles, groups []string) (*User, error) {
	id := uuid.New().String()
	user := User{ID: id, Username: username, Password: HashAndSalt([]byte(password)), GivenName: givenName, FamilyName: familyName, Email: email, Roles: strings.Join(roles[:], ","), Groups: strings.Join(groups[:], ",")}
	userBytes, _ := json.Marshal(user)
	queryString := fmt.Sprintf("post record %v.users %v", config.Config.DBName, string(userBytes))
	_, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update an existsing user with the data provided
func UpdateUserContents(id string, updateInput User) error {
	userBytes, _ := json.Marshal(updateInput)
	var userInterface map[string]interface{}
	json.Unmarshal(userBytes, &userInterface)
	queryString := fmt.Sprintf("get record %v.users | filter id = \"%v\"", config.Config.DBName, id)
	data, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("User does not exist!")
	}
	userInterface[".id"] = data[0][".id"]
	userBytesFull, _ := json.Marshal(userInterface)

	queryString = fmt.Sprintf("put record %v.users %v", config.Config.DBName, string(userBytesFull))
	_, err = connection.Query(queryString)
	if err != nil {
		return err
	}
	return nil
}

func IsUserPartOfGroup(id, group string) bool {
	userList := GetAllUsers()
	for _, u := range userList {
		if u.ID == id {
			for _, g := range strings.Split(u.Groups, ",") {
				if g == group {
					return true
				}
			}
		}
	}
	return false
}

func IsUserValid(username, password string) (bool, string) {
	userList := GetAllUsers()
	for _, u := range userList {
		if u.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			if err != nil {
				log.Println(err)
				return false, ""
			}
			return true, u.ID
		}
	}
	return false, ""
}

// Check if the supplied username is available
func IsUsernameAvailable(username string) bool {
	userList := GetAllUsers()
	for _, u := range userList {
		if u.Username == username {
			return false
		}
	}
	return true
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func PerformLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check if the username/password combination is valid
	is_valid, user_id := IsUserValid(username, password)
	if is_valid {
		// If the username/password is valid set the token in a cookie
		token := generateSessionToken()

		c.SetCookie("token", token, 3600, "", "", false, false)
		c.SetCookie("userId", user_id, 3600, "", "", false, false)
		c.Set("is_logged_in", true)
		c.Set("user", user_id)

		c.Redirect(302, "/")
		// showIndexPage(c)

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "local_login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

func generateSessionToken() string {
	// We're using a random 16 character string as the session token
	// This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	return strconv.FormatInt(rand.Int63(), 16)
}

func Logout(c *gin.Context) {

	// var sameSiteCookie http.SameSite;

	// Clear the cookie
	// c.SetCookie("token", "", -1, "", "", sameSiteCookie, false, true)
	c.SetCookie("token", "", -1, "", "", false, true)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
