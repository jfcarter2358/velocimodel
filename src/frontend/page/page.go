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

func ShowHomePage(c *gin.Context) {
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

	log.Println(obj)

	// Render the dashboard.html page
	render(c, gin.H{
		"title":    "Dashboard",
		"location": "Dashboard",
		"models":   obj},
		"dashboard.html")
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
