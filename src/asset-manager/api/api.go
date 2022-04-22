package api

import (
	"asset-manager/asset"
	"asset-manager/config"
	"errors"
	"fmt"
	"frontend/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"

var Healthy = false

type CreateGitObject struct {
	Repo       string `json:"repo"`
	Credential string `json:"credential"`
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
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := asset.DeleteAsset(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetAssets(c *gin.Context) {
	filter := FILTER_DEFAULT
	limit := LIMIT_DEFAULT
	count := COUNT_DEFAULT
	if val, ok := c.GetQuery("filter"); ok {
		filter = val
	}
	if val, ok := c.GetQuery("limit"); ok {
		limit = val
	}
	if val, ok := c.GetQuery("count"); ok {
		count = val
	}
	data, err := asset.GetAssets(limit, filter, count)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostAsset(c *gin.Context) {
	var input asset.Asset
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	fileID := uuid.New().String()
	input.ID = fileID
	err := asset.RegisterAsset(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

func PutAsset(c *gin.Context) {
	var input asset.Asset
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := asset.UpdateAsset(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func UploadAsset(c *gin.Context) {
	file, err := c.FormFile("file")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, filepath.Join(config.Config.DataPath, newFileName)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	fileID := uuid.New().String()
	assetData := asset.Asset{
		ID:       fileID,
		URL:      config.Params["basepath"].(string) + "/asset-manager/file/" + fileID,
		Name:     uuid.New().String(),
		Type:     "file",
		Tags:     make([]string, 0),
		Metadata: make(map[string]interface{}),
	}
	err = asset.RegisterAsset(assetData)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

func CreateGitAsset(c *gin.Context) {
	var input CreateGitObject
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	file, err := ioutil.TempFile("/tmp", "repo-")
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	defer os.Remove(file.Name())

	gitEmail := ""
	gitName := ""
	gitUser := ""
	gitPass := ""

	if val, ok := config.Params[fmt.Sprintf("git_%v_email", input.Credential)]; ok {
		gitEmail = val.(string)
	} else {
		utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, email does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
		return
	}
	if val, ok := config.Params[fmt.Sprintf("git_%v_name", input.Credential)]; ok {
		gitName = val.(string)
	} else {
		utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, name does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
		return
	}
	if val, ok := config.Params[fmt.Sprintf("git_%v_user", input.Credential)]; ok {
		gitUser = val.(string)
	} else {
		utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, user does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
		return
	}
	if val, ok := config.Params[fmt.Sprintf("git_%v_pass", input.Credential)]; ok {
		gitPass = val.(string)
	} else {
		utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, password does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
		return
	}

	log.Println("Initializing git configuration")
	_, err = exec.Command("git", "config", "--global", "user.email", fmt.Sprintf("\"%v\"", gitEmail)).Output()
	if err != nil {
		log.Println("Error setting email")
		log.Fatal(err)
	}
	_, err = exec.Command("git", "config", "--global", "user.name", fmt.Sprintf("\"%v\"", gitName)).Output()
	if err != nil {
		log.Println("Error setting name")
		log.Fatal(err)
	}

	cmd.Close()
}
