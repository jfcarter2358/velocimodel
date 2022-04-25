// utils.go

package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"model-manager/config"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func Error(err error, c *gin.Context, statusCode int) {
	log.Printf("Encountered error: %v", err)
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func CollectObjects(object map[string]interface{}) (string, string, error) {
	tempDir := fmt.Sprintf("/tmp/%v", object["id"].(string))
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		return "", "", err
	}

	log.Printf("TEMPDIR: %v", tempDir)

	assetIDs := make([]string, len(object["assets"].([]interface{})))
	for idx, val := range object["assets"].([]interface{}) {
		assetIDs[idx] = val.(string)
	}

	log.Printf("ASSET IDS: %v", assetIDs)

	var obj []map[string]interface{}
	requestURL := fmt.Sprintf("%v/api/asset", config.Config.APIServerURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return "", "", err
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		return "", "", err
	}

	assets := make([]map[string]interface{}, 0)

	for _, val := range obj {
		if Contains(assetIDs, val["id"].(string)) {
			assets = append(assets, val)
		}
	}

	log.Printf("ASSETS: %v", assets)

	for _, asset := range assets {
		switch asset["type"].(string) {
		case "git":
			log.Println("GOT GIT ASSET")
			err = doGitClone(
				asset["url"].(string),
				asset["metadata"].(map[string]interface{})["branch"].(string),
				asset["metadata"].(map[string]interface{})["credential"].(string),
				tempDir,
			)
			if err != nil {
				return "", "", err
			}
		case "file":
			log.Println("GOT FILE ASSET")
			err = doFileDownload(
				filepath.Join(
					tempDir,
					asset["metadata"].(map[string]interface{})["filename"].(string),
				),
				asset["url"].(string),
			)
			if err != nil {
				return "", "", err
			}
		}
	}
	// tar + gzip
	var buf bytes.Buffer
	err = compress(tempDir, &buf)

	filename := fmt.Sprintf("%v.tar.gz", object["id"].(string))
	localPath := fmt.Sprintf("/tmp/%v", filename)

	log.Printf("FILENAME: %v", filename)
	log.Printf("LOCALPATH: %v", localPath)

	// write the .tar.gzip
	fileToWrite, err := os.OpenFile(localPath, os.O_CREATE|os.O_RDWR, os.FileMode(600))
	if err != nil {
		return "", "", err
	}
	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		return "", "", err
	}
	return localPath, filename, nil
}

func doFileDownload(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func doGitClone(repo, branch, credential, tempDir string) error {
	gitEmail := ""
	gitName := ""
	gitUser := ""
	gitPass := ""

	// Grab credentials from service-manager
	if val, ok := config.Params["git_email"]; ok {
		gitEmail = val.(string)
	} else {
		return errors.New(fmt.Sprintf("Invalid git configuration, git_email does not exist"))
	}
	if val, ok := config.Params["git_name"]; ok {
		gitName = val.(string)
	} else {
		return errors.New(fmt.Sprintf("Invalid git configuration, git_name parameter does not exist"))
	}
	if credential != "none" {
		if val, ok := config.Secrets[fmt.Sprintf("git_%v_user", credential)]; ok {
			gitUser = val.(string)
		} else {
			return errors.New(fmt.Sprintf("Invalid git configuration, user does not exist for credential %v", credential))
		}
		if val, ok := config.Secrets[fmt.Sprintf("git_%v_pass", credential)]; ok {
			gitPass = val.(string)
		} else {
			return errors.New(fmt.Sprintf("Invalid git configuration, password does not exist for credential %v", credential))
		}
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
	prefix := "https://"
	if strings.HasPrefix(repo, "http://") {
		prefix = "http://"
	} else {
		if !strings.HasPrefix(repo, "https://") {
			repo = "https://" + repo
		}
	}
	domain := repo[len(prefix):]
	repoNameArr := strings.Split(repo, "/")
	repoName := repoNameArr[len(repoNameArr)-1]
	if credential != "none" {
		out, err = exec.Command("git", "clone", "--depth", "1", "-b", branch, fmt.Sprintf("%v%v:%v@%v", prefix, gitUser, gitPass, domain), filepath.Join(tempDir, repoName)).CombinedOutput()
		if err != nil {
			log.Printf("Git clone: %v", string(out))
			return err
		}
	} else {
		out, err = exec.Command("git", "clone", "--depth", "1", "-b", branch, fmt.Sprintf("%v%v", prefix, domain), filepath.Join(tempDir, repoName)).CombinedOutput()
		if err != nil {
			log.Printf("Git clone: %v", string(out))
			return err
		}
	}
	return nil
}

func compress(src string, buf io.Writer) error {
	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(buf)

	// walk through every file in the folder
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}
	//
	return nil
}
