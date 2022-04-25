package api

import (
	"asset-manager/asset"
	"asset-manager/config"
	"asset-manager/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	Branch     string `json:"branch"`
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

func CreateFileAsset(c *gin.Context) {
	file, err := c.FormFile("file")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	fileID := uuid.New().String()

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := fileID + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, filepath.Join(config.Config.DataPath, newFileName)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	assetData := asset.Asset{
		ID:       fileID,
		URL:      config.Params["apiserver"].(string) + "/api/asset/file/" + fileID,
		Name:     fileID,
		Type:     "file",
		Tags:     make([]string, 0),
		Metadata: map[string]interface{}{"filename": filepath.Base(file.Filename)},
		Models:   make([]string, 0),
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
	dir, err := ioutil.TempDir("/tmp", "repo-")
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	defer os.Remove(dir)

	gitEmail := ""
	gitName := ""
	gitUser := ""
	gitPass := ""

	// Grab credentials from service-manager
	if val, ok := config.Params["git_email"]; ok {
		gitEmail = val.(string)
	} else {
		utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, email does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
		return
	}
	if val, ok := config.Params["git_name"]; ok {
		gitName = val.(string)
	} else {
		utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, name does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
		return
	}
	if input.Credential != "none" {
		if val, ok := config.Secrets[fmt.Sprintf("git_%v_user", input.Credential)]; ok {
			gitUser = val.(string)
		} else {
			utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, user does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
			return
		}
		if val, ok := config.Secrets[fmt.Sprintf("git_%v_pass", input.Credential)]; ok {
			gitPass = val.(string)
		} else {
			utils.Error(errors.New(fmt.Sprintf("Invalid git configuration, password does not exist for credential %v", input.Credential)), c, http.StatusInternalServerError)
			return
		}
	}

	// Configure git
	log.Println("Initializing git configuration")
	out, err := exec.Command("git", "config", "--global", "user.email", fmt.Sprintf("\"%v\"", gitEmail)).CombinedOutput()
	if err != nil {
		log.Printf("Git config email: %v", string(out))
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	out, err = exec.Command("git", "config", "--global", "user.name", fmt.Sprintf("\"%v\"", gitName)).CombinedOutput()
	if err != nil {
		log.Printf("Git config name: %v", string(out))
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	// Shallow clone the git repo
	prefix := "https://"
	if strings.HasPrefix(input.Repo, "http://") {
		prefix = "http://"
	} else {
		if !strings.HasPrefix(input.Repo, "https://") {
			input.Repo = "https://" + input.Repo
		}
	}
	domain := input.Repo[len(prefix):]
	if input.Credential != "none" {
		out, err = exec.Command("git", "clone", "--depth", "1", "-b", input.Branch, fmt.Sprintf("%v%v:%v@%v", prefix, gitUser, gitPass, domain), dir).CombinedOutput()
		if err != nil {
			log.Printf("Git clone: %v", string(out))
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}
	} else {
		out, err = exec.Command("git", "clone", "--depth", "1", "-b", input.Branch, fmt.Sprintf("%v%v", prefix, domain), dir).CombinedOutput()
		if err != nil {
			log.Printf("Git clone: %v", string(out))
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}
	}
	out, err = exec.Command("git", "-C", dir, "log", "--format=\"%H\"", "-n", "1").CombinedOutput()
	if err != nil {
		log.Printf("Git log: %v", string(out))
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	commitID := out
	out, err = exec.Command("git", "-C", dir, "show", "-s", "--date=format:'%Y-%m-%dT%H:%M:%SZ'", "--format=%cd").CombinedOutput()
	if err != nil {
		log.Printf("Git show: %v", string(out))
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	commitTimestamp := out

	// Get the structure of the repo
	tree := make(map[string]interface{})
	visit := func(path string, info os.FileInfo, err error) error {
		keys := strings.Split(path, "/")
		if len(keys) <= 3 {
			return nil
		}
		keys = keys[3:]
		if keys[0] == ".git" {
			return nil
		}
		if info.IsDir() {
			tree = recurseAddTree(tree, keys, "dir")
		} else {
			tree = recurseAddTree(tree, keys, "file")
		}
		return nil
	}

	err = filepath.Walk(dir, visit)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	// Create the actual asset
	currentTime := time.Now().UTC()
	fileID := uuid.New().String()
	assetData := asset.Asset{
		ID:   fileID,
		URL:  input.Repo,
		Name: uuid.New().String(),
		Type: "git",
		Tags: make([]string, 0),
		Metadata: map[string]interface{}{
			"commit":           commitID,
			"branch":           input.Branch,
			"commitTimestamp":  commitTimestamp,
			"refreshTimestamp": currentTime.Format("2006-01-02T15:04:05Z"),
			"structure":        tree,
			"credential":       input.Credential,
		},
		Models: make([]string, 0),
	}
	err = asset.RegisterAsset(assetData)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

func DownloadFileAsset(c *gin.Context) {
	assetID := c.Param("id")

	assets, err := asset.GetAssets(LIMIT_DEFAULT, fmt.Sprintf("id = \"%v\"", assetID), COUNT_DEFAULT)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	filename := assets[0]["metadata"].(map[string]interface{})["filename"].(string)
	extension := filepath.Ext(filename)

	localPath := filepath.Join(config.Config.DataPath, fmt.Sprintf("%v%v", assetID, extension))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	c.FileAttachment(localPath, filename)

	c.Status(http.StatusOK)
}

func recurseAddTree(tree map[string]interface{}, keys []string, fileType string) map[string]interface{} {
	if len(keys) == 1 {
		if fileType == "file" {
			tree[keys[0]] = "file"
		} else {
			if tree[keys[0]] == nil {
				tree[keys[0]] = make(map[string]interface{})
			}
		}
	} else {
		if tree[keys[0]] != nil {
			tree[keys[0]] = recurseAddTree(tree[keys[0]].(map[string]interface{}), keys[1:], fileType)
		} else {
			tree[keys[0]] = recurseAddTree(make(map[string]interface{}), keys[1:], fileType)
		}
	}
	return tree
}
