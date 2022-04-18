// pages.go

package page

import (
	"encoding/json"
	"fmt"
	"frontend/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func RedirectIndexPage(c *gin.Context) {
	c.Redirect(301, "/ui/dashboard")
}

func ShowDashboardPage(c *gin.Context) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "updated")
	params.Add("limit", "10")
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Render the dashboard.html page
	render(c, gin.H{
		"models": obj},
		"dashboard.html")
}

func ShowAssetPage(c *gin.Context) {
	assetID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", assetID))
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	}

	// Render the models.html page
	render(c, gin.H{
		"asset": obj[0]},
		"asset.html")
}

func ShowAssetsPage(c *gin.Context) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "name")
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Render the models.html page
	render(c, gin.H{
		"assets": obj},
		"assets.html")
}

func ShowAssetEditPage(c *gin.Context) {
	assetID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", assetID))
	requestURL := fmt.Sprintf("%v/api/asset?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	} else {
		delete(obj[0], ".id")
		delete(obj[0], "created")
		delete(obj[0], "id")
		delete(obj[0], "url")
		delete(obj[0], "type")
		delete(obj[0], "updated")
	}

	jsonString, _ := json.Marshal(obj[0])

	// Render the asset.html page
	render(c, gin.H{
		"asset_json": string(jsonString),
		"asset":      obj[0]},
		"asset-edit.html")
}

func ShowModelPage(c *gin.Context) {
	modelID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", modelID))
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	}

	// Render the models.html page
	render(c, gin.H{
		"model": obj[0]},
		"model.html")
}

func ShowModelsPage(c *gin.Context) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "name")
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Render the models.html page
	render(c, gin.H{
		"models": obj},
		"models.html")
}

func ShowModelEditPage(c *gin.Context) {
	modelID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", modelID))
	requestURL := fmt.Sprintf("%v/api/model?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	} else {
		delete(obj[0], ".id")
		delete(obj[0], "created")
		delete(obj[0], "id")
		delete(obj[0], "language")
		delete(obj[0], "type")
		delete(obj[0], "updated")
	}

	jsonString, _ := json.Marshal(obj[0])

	// Render the models.html page
	render(c, gin.H{
		"model_json": string(jsonString),
		"model":      obj[0]},
		"model-edit.html")
}

func ShowReleasePage(c *gin.Context) {
	releaseID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", releaseID))
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	}

	// Render the models.html page
	render(c, gin.H{
		"release": obj[0]},
		"release.html")
}

func ShowReleasesPage(c *gin.Context) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "name")
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Render the models.html page
	render(c, gin.H{
		"releases": obj},
		"releases.html")
}

func ShowReleaseEditPage(c *gin.Context) {
	releaseID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", releaseID))
	requestURL := fmt.Sprintf("%v/api/release?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	} else {
		delete(obj[0], ".id")
		delete(obj[0], "created")
		delete(obj[0], "id")
		delete(obj[0], "language")
		delete(obj[0], "type")
		delete(obj[0], "updated")
	}

	jsonString, _ := json.Marshal(obj[0])

	// Render the models.html page
	render(c, gin.H{
		"release_json": string(jsonString),
		"release":      obj[0]},
		"release-edit.html")
}

func ShowSnapshotPage(c *gin.Context) {
	snapshotID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", snapshotID))
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	}

	// Render the models.html page
	render(c, gin.H{
		"snapshot": obj[0]},
		"snapshot.html")
}

func ShowSnapshotsPage(c *gin.Context) {
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("orderdsc", "name")
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Render the models.html page
	render(c, gin.H{
		"snapshots": obj},
		"snapshots.html")
}

func ShowSnapshotEditPage(c *gin.Context) {
	snapshotID := c.Param("id")
	var obj []map[string]interface{}
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("id = \"%v\"", snapshotID))
	requestURL := fmt.Sprintf("%v/api/snapshot?%v", config.Config.APIServerURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Encountered error: Request failed with status code %v", resp.StatusCode)
		c.Status(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Printf("Encountered error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(obj) == 0 {
		obj = append(obj, map[string]interface{}{})
	} else {
		delete(obj[0], ".id")
		delete(obj[0], "created")
		delete(obj[0], "id")
		delete(obj[0], "language")
		delete(obj[0], "type")
		delete(obj[0], "updated")
	}

	jsonString, _ := json.Marshal(obj[0])

	// Render the models.html page
	render(c, gin.H{
		"snapshot_json": string(jsonString),
		"snapshot":      obj[0]},
		"snapshot-edit.html")
}

func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
