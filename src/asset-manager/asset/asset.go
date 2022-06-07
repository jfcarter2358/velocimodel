package asset

import (
	"asset-manager/action"
	"asset-manager/config"
	"asset-manager/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

type Asset struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	URL      string                 `json:"url"`
	Created  string                 `json:"created"`
	Updated  string                 `json:"updated"`
	Type     string                 `json:"type"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Models   []string               `json:"models"`
}

const FILE_TYPE = "file"
const GIT_TYPE = "git"
const S3_TYPE = "s3"
const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"
const ORDERASC_DEFAULT = "NA"
const ORDERDSC_DEFAULT = "NA"

var allowedTypes = []string{FILE_TYPE, GIT_TYPE, S3_TYPE}

func RegisterAsset(newAsset Asset) error {
	if newAsset.ID == "" {
		newAsset.ID = uuid.New().String()
	}
	if !utils.Contains(allowedTypes, newAsset.Type) {
		err := fmt.Errorf("asset type of %v does not exist", newAsset.Type)
		return err
	}
	currentTime := time.Now().UTC()
	newAsset.Created = currentTime.Format("2006-01-02T15:04:05Z")
	newAsset.Updated = currentTime.Format("2006-01-02T15:04:05Z")
	queryData, _ := json.Marshal(&newAsset)
	queryString := fmt.Sprintf("post record %v.assets %v", config.Config.DB.Name, string(queryData))
	_, err := connection.Query(queryString)
	return err
}

func DeleteAsset(assetIDs []string) error {
	queryString := fmt.Sprintf("get record %v.assets", config.Config.DB.Name)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(assetIDs, datum["id"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.assets %v", config.Config.DB.Name, string(queryData))
	_, err = connection.Query(queryString)
	return err
}

func UpdateAsset(newAsset Asset) error {
	if newAsset.ID == "" {
		err := errors.New("'id' field is required to update an asset")
		return err
	}
	queryString := fmt.Sprintf("get record %v.assets", config.Config.DB.Name)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	for _, datum := range currentData {
		if datum["id"].(string) == newAsset.ID {
			if newAsset.Name != "" {
				if newAsset.Name != datum["name"].(string) {
					datum["name"] = newAsset.Name
				}
			}
			if newAsset.Tags != nil {
				tmpTags := make([]string, len(datum["tags"].([]interface{})))
				for idx, val := range datum["tags"].([]interface{}) {
					tmpTags[idx] = val.(string)
				}
				if !reflect.DeepEqual(newAsset.Tags, tmpTags) {
					datum["tags"] = newAsset.Tags
				}
			}
			if newAsset.Metadata != nil {
				if !reflect.DeepEqual(newAsset.Metadata, datum["metadata"].(map[string]interface{})) {
					datum["metadata"] = newAsset.Metadata
				}
			}
			if newAsset.Models != nil {
				tmpModels := make([]string, len(datum["models"].([]interface{})))
				for idx, val := range datum["models"].([]interface{}) {
					tmpModels[idx] = val.(string)
				}
				if !reflect.DeepEqual(newAsset.Models, tmpModels) {
					datum["models"] = newAsset.Models
				}
			}
			currentTime := time.Now().UTC()
			datum["updated"] = currentTime.Format("2006-01-02T15:04:05Z")
			queryData, _ := json.Marshal(&datum)
			queryString := fmt.Sprintf("put record %v.assets %v", config.Config.DB.Name, string(queryData))
			_, err := connection.Query(queryString)
			if err != nil {
				return err
			}
			return nil
		}
	}
	err = RegisterAsset(newAsset)
	return err
}

func GetAssets(limit, filter, count, orderasc, orderdsc string) ([]Asset, error) {
	queryString := fmt.Sprintf("get record %v.assets", config.Config.DB.Name)
	if filter != FILTER_DEFAULT {
		queryString += fmt.Sprintf(" | filter %v", filter)
	}
	if limit != LIMIT_DEFAULT {
		queryString += fmt.Sprintf(" | limit %v", limit)
	}
	if count != COUNT_DEFAULT {
		queryString += " | count"
	}
	if orderdsc != ORDERDSC_DEFAULT {
		queryString += fmt.Sprintf(" | orderdsc %v", orderdsc)
	}
	if count != COUNT_DEFAULT {
		queryString += " | count"
	}
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	marshalled, _ := json.Marshal(data)
	var output []Asset
	json.Unmarshal(marshalled, &output)
	return output, nil
}

func DoGitSync(input Asset) error {
	dir, err := ioutil.TempDir("/tmp", "repo-")
	if err != nil {
		return err
	}
	defer os.Remove(dir)

	gitEmail := ""
	gitName := ""

	// Grab credentials from service-manager
	if val, ok := config.Params["git_email"]; ok {
		gitEmail = val.(string)
	} else {
		return errors.New("invalid git configuration, git_email does not exist")
	}
	if val, ok := config.Params["git_name"]; ok {
		gitName = val.(string)
	} else {
		return errors.New("invalid git configuration, git_name does not exist")
	}

	// Configure git
	log.Println("Initializing git configuration")
	out, err := exec.Command("git", "config", "--global", "user.email", fmt.Sprintf("\"%v\"", gitEmail)).CombinedOutput()
	if err != nil {
		log.Printf("Git config email: %v", string(out))
		return err
	}
	out, err = exec.Command("git", "config", "--global", "user.name", fmt.Sprintf("\"%v\"", gitName)).CombinedOutput()
	if err != nil {
		log.Printf("Git config name: %v", string(out))
		return err
	}

	// Shallow clone the git repo
	if strings.HasPrefix(input.URL, "git@") {
		sshKey := ""
		if input.Metadata["credential"].(string) != "none" {
			if val, ok := config.Secrets[fmt.Sprintf("git_%v_ssh_key", input.Metadata["credential"].(string))]; ok {
				sshKey = val.(string)
			} else {
				return fmt.Errorf("invalid git configuration, ssh key does not exist for credential %v", input.Metadata["credential"].(string))
			}
		}

		dec, err := base64.StdEncoding.DecodeString(sshKey)
		if err != nil {
			panic(err)
		}

		f, err := ioutil.TempFile("/tmp", "ssh-key-")
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			return err
		}
		if err := f.Sync(); err != nil {
			return err
		}
		defer os.Remove(f.Name())

		parts := strings.Split(input.URL, "@")
		domain := strings.Split(parts[1], ":")[0]

		out, err := exec.Command("ssh-keyscan", "-H", domain).CombinedOutput()
		if err != nil {
			log.Printf("Keyscan: %v", string(out))
			return err
		}
		knownHostsLoc := os.ExpandEnv("$HOME/.ssh/known_hosts")
		khf, err := os.OpenFile(knownHostsLoc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Keyscan write 1: %v", err)
			return err
		}
		defer khf.Close()
		if _, err := khf.WriteString(string(out)); err != nil {
			log.Printf("Keyscan write 2: %v", err)
			return err
		}
		err = os.Chmod(f.Name(), 0600)
		if err != nil {
			log.Printf("Permission change: %v", err)
			return err
		}

		if input.Metadata["credential"].(string) == "none" {
			log.Printf("SSH Key is required")
			return errors.New("ssh key is required")
		}

		out, err = exec.Command("/bin/bash", "-c", fmt.Sprintf("GIT_SSH_COMMAND='ssh -i %v' git clone --depth 1 -b %v %v %v", f.Name(), input.Metadata["branch"].(string), input.URL, dir)).CombinedOutput()
		if err != nil {
			log.Printf("Git clone: %v", string(out))
			return err
		}
	} else {
		gitUser := ""
		gitPass := ""
		if input.Metadata["credential"].(string) != "none" {
			if val, ok := config.Secrets[fmt.Sprintf("git_%v_user", input.Metadata["credential"].(string))]; ok {
				gitUser = val.(string)
			} else {
				return fmt.Errorf("invalid git configuration, user does not exist for credential %v", input.Metadata["credential"].(string))
			}
			if val, ok := config.Secrets[fmt.Sprintf("git_%v_pass", input.Metadata["credential"].(string))]; ok {
				gitPass = val.(string)
			} else {
				return fmt.Errorf("invalid git configuration, password does not exist for credential %v", input.Metadata["credential"].(string))
			}
		}

		prefix := "https://"
		if strings.HasPrefix(input.URL, "http://") {
			prefix = "http://"
		} else {
			if !strings.HasPrefix(input.URL, "https://") {
				input.URL = "https://" + input.URL
			}
		}
		domain := input.URL[len(prefix):]
		if input.Metadata["credential"].(string) != "none" {
			out, err = exec.Command("git", "clone", "--depth", "1", "-b", input.Metadata["branch"].(string), fmt.Sprintf("%v%v:%v@%v", prefix, gitUser, gitPass, domain), dir).CombinedOutput()
			if err != nil {
				log.Printf("Git clone: %v", string(out))
				return err
			}
		} else {
			out, err = exec.Command("git", "clone", "--depth", "1", "-b", input.Metadata["branch"].(string), fmt.Sprintf("%v%v", prefix, domain), dir).CombinedOutput()
			if err != nil {
				log.Printf("Git clone: %v", string(out))
				return err
			}
		}
	}
	out, err = exec.Command("git", "-C", dir, "log", "--format=\"%H\"", "-n", "1").CombinedOutput()
	if err != nil {
		log.Printf("Git log: %v", string(out))
		return err
	}
	commitID := out[1 : len(out)-2]
	out, err = exec.Command("git", "-C", dir, "show", "-s", "--date=format:'%Y-%m-%dT%H:%M:%SZ'", "--format=%cd").CombinedOutput()
	if err != nil {
		log.Printf("Git show: %v", string(out))
		return err
	}
	commitTimestamp := out[1 : len(out)-2]

	if input.Metadata["commit"].(string) != string(commitID) {
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
			return err
		}

		// Create the actual asset
		currentTime := time.Now().UTC()
		assetID := uuid.New().String()
		assetData := Asset{
			ID:   assetID,
			URL:  input.URL,
			Name: input.Name,
			Type: "git",
			Tags: make([]string, 0),
			Metadata: map[string]interface{}{
				"commit":           string(commitID),
				"branch":           input.Metadata["branch"],
				"commitTimestamp":  string(commitTimestamp),
				"refreshTimestamp": currentTime.Format("2006-01-02T15:04:05Z"),
				"structure":        tree,
				"credential":       input.Metadata["credential"],
			},
			Models: input.Models,
		}
		err = RegisterAsset(assetData)
		if err != nil {
			return err
		}
		for _, modelID := range input.Models {
			data := map[string]interface{}{
				"model": modelID,
				"asset": input.ID,
			}
			err := action.SendDelete("api/model/asset", data)
			if err != nil {
				return err
			}
			data = map[string]interface{}{
				"model": modelID,
				"asset": assetID,
			}
			_, err = action.SendPost("api/model/asset", data)
			if err != nil {
				return err
			}
		}
		input.Models = make([]string, 0)
		input.Metadata["refreshTimestamp"] = currentTime.Format("2006-01-02T15:04:05Z")
		err = UpdateAsset(input)
		if err != nil {
			return err
		}
	} else {
		currentTime := time.Now().UTC()
		input.Metadata["refreshTimestamp"] = currentTime.Format("2006-01-02T15:04:05Z")
		err := UpdateAsset(input)
		if err != nil {
			return err
		}
	}
	return nil
}
