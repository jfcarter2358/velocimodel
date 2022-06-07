package api

import (
	"asset-manager/action"
	"asset-manager/asset"
	"asset-manager/config"
	"asset-manager/utils"
	"encoding/base64"
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
const ORDERASC_DEFAULT = "NA"
const ORDERDSC_DEFAULT = "NA"

var Healthy = false

type CreateGitObject struct {
	Repo       string `json:"repo"`
	Branch     string `json:"branch"`
	Credential string `json:"credential"`
}

// Health API

func GetHealth(c *gin.Context) {
	if !Healthy {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Status(http.StatusOK)
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
	orderasc := ORDERASC_DEFAULT
	orderdsc := ORDERDSC_DEFAULT
	if val, ok := c.GetQuery("filter"); ok {
		filter = val
	}
	if val, ok := c.GetQuery("limit"); ok {
		limit = val
	}
	if val, ok := c.GetQuery("count"); ok {
		count = val
	}
	if val, ok := c.GetQuery("orderasc"); ok {
		orderasc = val
	}
	if val, ok := c.GetQuery("orderdsc"); ok {
		orderdsc = val
	}
	data, err := asset.GetAssets(limit, filter, count, orderasc, orderdsc)
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

	// Grab credentials from service-manager
	if val, ok := config.Params["git_email"]; ok {
		gitEmail = val.(string)
	} else {
		utils.Error(fmt.Errorf("invalid git configuration, email does not exist for credential %v", input.Credential), c, http.StatusInternalServerError)
		return
	}
	if val, ok := config.Params["git_name"]; ok {
		gitName = val.(string)
	} else {
		utils.Error(fmt.Errorf("invalid git configuration, name does not exist for credential %v", input.Credential), c, http.StatusInternalServerError)
		return
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
	if strings.HasPrefix(input.Repo, "git@") {
		sshKey := ""
		if input.Credential != "none" {
			if val, ok := config.Secrets[fmt.Sprintf("git_%v_ssh_key", input.Credential)]; ok {
				sshKey = val.(string)
			} else {
				utils.Error(fmt.Errorf("invalid git configuration, ssh key does not exist for credential %v", input.Credential), c, http.StatusInternalServerError)
				return
			}
		}

		dec, err := base64.StdEncoding.DecodeString(sshKey)
		if err != nil {
			panic(err)
		}

		f, err := ioutil.TempFile("/tmp", "ssh-key-")
		if err != nil {
			utils.Error(fmt.Errorf("could not open ssh key file for writing"), c, http.StatusInternalServerError)
			return
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			utils.Error(fmt.Errorf("could not write ssh key"), c, http.StatusInternalServerError)
			return
		}
		if err := f.Sync(); err != nil {
			utils.Error(fmt.Errorf("could not sync ssh key file"), c, http.StatusInternalServerError)
			return
		}
		defer os.Remove(f.Name())

		parts := strings.Split(input.Repo, "@")
		domain := strings.Split(parts[1], ":")[0]

		out, err := exec.Command("ssh-keyscan", "-H", domain).CombinedOutput()
		if err != nil {
			log.Printf("Keyscan: %v", string(out))
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}
		knownHostsLoc := os.ExpandEnv("$HOME/.ssh/known_hosts")
		khf, err := os.OpenFile(knownHostsLoc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Keyscan write 1: %v", err)
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}
		defer khf.Close()
		if _, err := khf.WriteString(string(out)); err != nil {
			log.Printf("Keyscan write 2: %v", err)
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}
		err = os.Chmod(f.Name(), 0600)
		if err != nil {
			log.Printf("Permission change: %v", err)
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}

		if input.Credential == "none" {
			log.Printf("SSH Key is required")
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}

		out, err = exec.Command("/bin/bash", "-c", fmt.Sprintf("GIT_SSH_COMMAND='ssh -i %v' git clone --depth 1 -b %v %v %v", f.Name(), input.Branch, input.Repo, dir)).CombinedOutput()
		if err != nil {
			log.Printf("Git clone: %v", string(out))
			utils.Error(err, c, http.StatusInternalServerError)
			return
		}

	} else {
		gitUser := ""
		gitPass := ""
		if input.Credential != "none" {
			if val, ok := config.Secrets[fmt.Sprintf("git_%v_user", input.Credential)]; ok {
				gitUser = val.(string)
			} else {
				utils.Error(fmt.Errorf("invalid git configuration, user does not exist for credential %v", input.Credential), c, http.StatusInternalServerError)
				return
			}
			if val, ok := config.Secrets[fmt.Sprintf("git_%v_pass", input.Credential)]; ok {
				gitPass = val.(string)
			} else {
				utils.Error(fmt.Errorf("invalid git configuration, password does not exist for credential %v", input.Credential), c, http.StatusInternalServerError)
				return
			}
		}

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
	}
	out, err = exec.Command("git", "-C", dir, "log", "--format=\"%H\"", "-n", "1").CombinedOutput()
	if err != nil {
		log.Printf("Git log: %v", string(out))
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	commitID := out[1 : len(out)-2]
	out, err = exec.Command("git", "-C", dir, "show", "-s", "--date=format:'%Y-%m-%dT%H:%M:%SZ'", "--format=%cd").CombinedOutput()
	if err != nil {
		log.Printf("Git show: %v", string(out))
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	commitTimestamp := out[1 : len(out)-2]

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
			tree = action.RecurseAddTree(tree, keys, "dir")
		} else {
			tree = action.RecurseAddTree(tree, keys, "file")
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
		Name: fileID,
		Type: "git",
		Tags: make([]string, 0),
		Metadata: map[string]interface{}{
			"commit":           string(commitID),
			"branch":           input.Branch,
			"commitTimestamp":  string(commitTimestamp),
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

func SyncGitAsset(c *gin.Context) {
	var input asset.Asset
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	if err := asset.DoGitSync(input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func DownloadFileAsset(c *gin.Context) {
	assetID := c.Param("id")

	assets, err := asset.GetAssets(LIMIT_DEFAULT, fmt.Sprintf("id = \"%v\"", assetID), COUNT_DEFAULT, ORDERASC_DEFAULT, ORDERDSC_DEFAULT)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	filename := assets[0].Metadata["filename"].(string)
	extension := filepath.Ext(filename)

	localPath := filepath.Join(config.Config.DataPath, fmt.Sprintf("%v%v", assetID, extension))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	c.FileAttachment(localPath, filename)

	c.Status(http.StatusOK)
}
