// pages.go

package page

import (
	"encoding/json"
	"frontend/action"
	"frontend/config"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func RedirectIndexPage(c *gin.Context) {
	c.Redirect(301, config.Config.HTTPBasePath+"/ui/dashboard")
}

func RedirectProxyIndexPage(c *gin.Context) {
	c.Redirect(301, config.Config.HTTPBasePath+"/ui/dashboard")
}

func ShowDashboardPage(c *gin.Context) {
	assets, err := action.GetAssetsLimitLatest(c, "10")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	models, err := action.GetModelsLimitLatest(c, "10")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	releases, err := action.GetReleasesLimitLatest(c, "10")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	snapshots, err := action.GetSnapshotsLimitLatest(c, "10")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"assets":         assets,
		"models":         models,
		"releases":       releases,
		"snapshots":      snapshots,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"dashboard.html")
}

func ShowAssetPage(c *gin.Context) {
	assetID := c.Param("id")

	asset, err := action.GetAssetByID(c, assetID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if asset == nil {
		c.Status(http.StatusNotFound)
		return
	}

	modelIDList := make([]string, len(asset["models"].([]interface{})))
	for idx, val := range asset["models"].([]interface{}) {
		modelIDList[idx] = val.(string)
	}
	models, err := action.GetModelsByIDList(c, modelIDList)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	tagObj := make([]map[string]string, len(asset["tags"].([]interface{})))
	for idx, val := range asset["tags"].([]interface{}) {
		tagObj[idx] = map[string]string{
			"value": val.(string),
			// "class": "w3-round velocimodel-green",
		}
	}

	tagJSON, _ := json.Marshal(tagObj)
	metadataJSON, _ := json.MarshalIndent(asset["metadata"], "", "    ")

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"asset":          asset,
		"models":         models,
		"tag_json":       string(tagJSON),
		"metadata_json":  string(metadataJSON),
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"asset.html")
}

func ShowAssetsPage(c *gin.Context) {
	assets, err := action.GetAssetsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	secrets, err := action.GetSecretsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	credentials := make([]string, 0)
	for key := range secrets {
		if strings.HasPrefix(key, "git_") {
			parts := strings.Split(key, "_")
			if !action.Contains(credentials, parts[1]) {
				credentials = append(credentials, parts[1])
			}
		}
	}
	sort.Strings(credentials)

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"assets":         assets,
		"credentials":    credentials,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"assets.html")
}

func ShowAssetCodePage(c *gin.Context) {
	assetID := c.Param("id")

	asset, err := action.GetAssetByID(c, assetID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if asset == nil {
		c.Status(http.StatusNotFound)
		return
	}

	delete(asset, ".id")
	delete(asset, "created")
	delete(asset, "id")
	delete(asset, "url")
	delete(asset, "type")
	delete(asset, "updated")

	jsonString, _ := json.Marshal(asset)

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"asset_id":       assetID,
		"asset_json":     string(jsonString),
		"asset":          asset,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"asset-code.html")
}

func ShowModelPage(c *gin.Context) {
	modelID := c.Param("id")

	model, err := action.GetModelByID(c, modelID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if model == nil {
		c.Status(http.StatusNotFound)
		return
	}

	assetIDList := make([]string, len(model["assets"].([]interface{})))
	for idx, val := range model["assets"].([]interface{}) {
		assetIDList[idx] = val.(string)
	}
	assets, err := action.GetAssetsByIDList(c, assetIDList)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for idx, val := range assets {
		assetMetadataJSON, _ := json.MarshalIndent(val["metadata"].(map[string]interface{}), "", "    ")
		val["metadata_json"] = string(assetMetadataJSON)
		assets[idx] = val
	}

	snapshotIDList := make([]string, len(model["snapshots"].([]interface{})))
	for idx, val := range model["snapshots"].([]interface{}) {
		snapshotIDList[idx] = val.(string)
	}
	snapshots, err := action.GetSnapshotsByIDList(c, snapshotIDList)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	releaseIDList := make([]string, len(model["releases"].([]interface{})))
	for idx, val := range model["releases"].([]interface{}) {
		releaseIDList[idx] = val.(string)
	}
	releases, err := action.GetReleasesByIDList(c, releaseIDList)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	secrets, err := action.GetSecretsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	credentials := make([]string, 0)
	for key := range secrets {
		if strings.HasPrefix(key, "git_") {
			parts := strings.Split(key, "_")
			if !action.Contains(credentials, parts[1]) {
				credentials = append(credentials, parts[1])
			}
		}
	}
	sort.Strings(credentials)

	allAssets, err := action.GetAssetsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for idx, val := range allAssets {
		delete(val, "created")
		delete(val, "url")
		delete(val, "tags")
		delete(val, "metadata")
		delete(val, "models")
		allAssets[idx] = val
	}

	tagObj := make([]map[string]string, len(model["tags"].([]interface{})))
	for idx, val := range model["tags"].([]interface{}) {
		tagObj[idx] = map[string]string{
			"value": val.(string),
			// "class": "w3-round velocimodel-green",
		}
	}

	tagJSON, _ := json.Marshal(tagObj)
	metadataJSON, _ := json.MarshalIndent(model["metadata"], "", "    ")

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"model":          model,
		"assets":         assets,
		"snapshots":      snapshots,
		"releases":       releases,
		"credentials":    credentials,
		"all_assets":     allAssets,
		"tag_json":       string(tagJSON),
		"metadata_json":  string(metadataJSON),
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"model.html")
}

func ShowModelsPage(c *gin.Context) {
	models, err := action.GetModelsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"models":         models,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"models.html")
}

func ShowModelCodePage(c *gin.Context) {
	modelID := c.Param("id")

	model, err := action.GetModelByID(c, modelID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if model == nil {
		c.Status(http.StatusNotFound)
		return
	}

	delete(model, ".id")
	delete(model, "created")
	delete(model, "id")
	delete(model, "language")
	delete(model, "type")
	delete(model, "updated")

	jsonString, _ := json.Marshal(model)

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"model_id":       modelID,
		"model_json":     string(jsonString),
		"model":          model,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"model-code.html")
}

func ShowReleasePage(c *gin.Context) {
	releaseID := c.Param("id")

	release, err := action.GetReleaseByID(c, releaseID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if release == nil {
		c.Status(http.StatusNotFound)
		return
	}

	assetIDList := make([]string, len(release["assets"].([]interface{})))
	for idx, val := range release["assets"].([]interface{}) {
		assetIDList[idx] = val.(string)
	}
	assets, err := action.GetAssetsByIDList(c, assetIDList)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	tagObj := make([]map[string]string, len(release["tags"].([]interface{})))
	for idx, val := range release["tags"].([]interface{}) {
		tagObj[idx] = map[string]string{
			"value": val.(string),
			// "class": "w3-round velocimodel-green",
		}
	}

	tagJSON, _ := json.Marshal(tagObj)
	metadataJSON, _ := json.MarshalIndent(release["metadata"], "", "    ")

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"release":        release,
		"assets":         assets,
		"tag_json":       string(tagJSON),
		"metadata_json":  string(metadataJSON),
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"release.html")
}

func ShowReleasesPage(c *gin.Context) {
	releases, err := action.GetReleasesAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	allSnapshots, err := action.GetSnapshotsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for idx, val := range allSnapshots {
		delete(val, "created")
		delete(val, "type")
		delete(val, "tags")
		delete(val, "metadata")
		delete(val, "assets")
		delete(val, "language")
		delete(val, "releases")
		allSnapshots[idx] = val
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"releases":       releases,
		"all_snapshots":  allSnapshots,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"releases.html")
}

func ShowReleaseCodePage(c *gin.Context) {
	releaseID := c.Param("id")

	release, err := action.GetReleaseByID(c, releaseID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if release == nil {
		c.Status(http.StatusNotFound)
		return
	}

	delete(release, ".id")
	delete(release, "created")
	delete(release, "id")
	delete(release, "language")
	delete(release, "type")
	delete(release, "updated")

	jsonString, _ := json.Marshal(release)

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"release_id":     releaseID,
		"release_json":   string(jsonString),
		"release":        release,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"release-code.html")
}

func ShowSnapshotPage(c *gin.Context) {
	snapshotID := c.Param("id")

	snapshot, err := action.GetSnapshotByID(c, snapshotID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if snapshot == nil {
		c.Status(http.StatusNotFound)
		return
	}

	assetIDList := make([]string, len(snapshot["assets"].([]interface{})))
	for idx, val := range snapshot["assets"].([]interface{}) {
		assetIDList[idx] = val.(string)
	}
	assets, err := action.GetAssetsByIDList(c, assetIDList)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	releaseIDList := make([]string, len(snapshot["releases"].([]interface{})))
	for idx, val := range snapshot["releases"].([]interface{}) {
		releaseIDList[idx] = val.(string)
	}
	releases, err := action.GetReleasesByIDList(c, releaseIDList)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	tagObj := make([]map[string]string, len(snapshot["tags"].([]interface{})))
	for idx, val := range snapshot["tags"].([]interface{}) {
		tagObj[idx] = map[string]string{
			"value": val.(string),
			// "class": "w3-round velocimodel-green",
		}
	}

	tagJSON, _ := json.Marshal(tagObj)
	metadataJSON, _ := json.MarshalIndent(snapshot["metadata"], "", "    ")

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"snapshot":       snapshot,
		"assets":         assets,
		"releases":       releases,
		"tag_json":       string(tagJSON),
		"metadata_json":  string(metadataJSON),
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"snapshot.html")
}

func ShowSnapshotsPage(c *gin.Context) {
	snapshots, err := action.GetSnapshotsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	allModels, err := action.GetModelsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for idx, val := range allModels {
		delete(val, "created")
		delete(val, "type")
		delete(val, "tags")
		delete(val, "metadata")
		delete(val, "assets")
		delete(val, "snapshots")
		delete(val, "releases")
		allModels[idx] = val
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"snapshots":      snapshots,
		"all_models":     allModels,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"snapshots.html")
}

func ShowSnapshotCodePage(c *gin.Context) {
	snapshotID := c.Param("id")

	snapshot, err := action.GetSnapshotByID(c, snapshotID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if snapshot == nil {
		c.Status(http.StatusNotFound)
		return
	}

	delete(snapshot, ".id")
	delete(snapshot, "created")
	delete(snapshot, "id")
	delete(snapshot, "language")
	delete(snapshot, "type")
	delete(snapshot, "updated")

	jsonString, _ := json.Marshal(snapshot)

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	// Render the models.html page
	render(c, gin.H{
		"snapshot_id":    snapshotID,
		"snapshot_json":  string(jsonString),
		"snapshot":       snapshot,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"snapshot-code.html")
}

func ShowUsersPage(c *gin.Context) {
	users, err := action.GetUsersAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	params, err := action.GetParamsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	available_groups := strings.Split(params["available_groups"].(string), ",")

	render(c, gin.H{
		"users":            users,
		"user_data":        userData,
		"is_admin_role":    isAdminRole,
		"is_admin_group":   isAdminGroup,
		"available_groups": available_groups,
		"base_path":        config.Config.HTTPBasePath},
		"users.html")
}

func ShowUserPage(c *gin.Context) {
	userID := c.Param("id")

	user, err := action.GetUserByID(userID, c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	params, err := action.GetParamsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	available_groups := strings.Split(params["available_groups"].(string), ",")
	available_roles := []string{"read", "write", "admin"}

	groups_assigned := map[string]bool{}
	roles_assigned := map[string]bool{}

	for _, group := range available_groups {
		if StringSliceContains(groups, group) {
			groups_assigned[group] = true
		} else {
			groups_assigned[group] = false
		}
	}

	for _, role := range available_roles {
		if StringSliceContains(roles, role) {
			roles_assigned[role] = true
		} else {
			roles_assigned[role] = false
		}
	}

	render(c, gin.H{
		"user":            user,
		"user_data":       userData,
		"is_admin_role":   isAdminRole,
		"is_admin_group":  isAdminGroup,
		"groups_assigned": groups_assigned,
		"roles_assigned":  roles_assigned,
		"base_path":       config.Config.HTTPBasePath},
		"user.html")
}

func ShowParamsPage(c *gin.Context) {
	params, err := action.GetParamsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"params":         params,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"params.html")
}

func ShowSecretsPage(c *gin.Context) {
	secrets, err := action.GetSecretsAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	userData, err := action.GetUserData(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	isAdminRole := false
	roles := make([]string, len(userData["roles"].([]interface{})))
	for idx, val := range userData["roles"].([]interface{}) {
		roles[idx] = val.(string)
	}
	if StringSliceContains(roles, "admin") {
		isAdminRole = true
	}
	isAdminGroup := false
	groups := make([]string, len(userData["groups"].([]interface{})))
	for idx, val := range userData["groups"].([]interface{}) {
		groups[idx] = val.(string)
	}
	if StringSliceContains(groups, "admin") {
		isAdminGroup = true
	}

	render(c, gin.H{
		"secrets":        secrets,
		"user_data":      userData,
		"is_admin_role":  isAdminRole,
		"is_admin_group": isAdminGroup,
		"base_path":      config.Config.HTTPBasePath},
		"secrets.html")
}

func ShowLoginPage(c *gin.Context) {
	render(c, gin.H{}, "login.html")
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

func StringSliceContains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
