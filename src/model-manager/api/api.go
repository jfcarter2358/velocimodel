package api

import (
	"fmt"
	"log"
	"model-manager/model"
	"model-manager/release"
	"model-manager/snapshot"
	"model-manager/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const LIMIT_DEFAULT = "0"
const FILTER_DEFAULT = ""
const COUNT_DEFAULT = "false"
const ORDERASC_DEFAULT = "NA"
const ORDERDSC_DEFAULT = "NA"

var Healthy = false

// Health API

func GetHealth(c *gin.Context) {
	if Healthy == false {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Status(http.StatusOK)
	return
}

// Model API

func AddAsset(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	log.Println("Adding asset")
	err := model.AddAsset(input["model"], input["asset"])
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func DeleteAsset(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := model.DeleteAsset(input["model"], input["asset"])
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func DeleteModel(c *gin.Context) {
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := model.DeleteModel(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetModels(c *gin.Context) {
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
	data, err := model.GetModels(limit, filter, count, orderasc, orderdsc)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostModel(c *gin.Context) {
	var input model.Model
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	fileID := uuid.New().String()
	input.ID = fileID
	err := model.RegisterModel(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

func PutModel(c *gin.Context) {
	var input model.Model
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := model.UpdateModel(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func DownloadModel(c *gin.Context) {
	modelID := c.Param("id")

	models, err := model.GetModels(LIMIT_DEFAULT, fmt.Sprintf("id = \"%v\"", modelID), COUNT_DEFAULT, ORDERASC_DEFAULT, ORDERDSC_DEFAULT)

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	log.Printf("MODEL: %v", models[0])

	localPath, filename, err := utils.CollectObjects(models[0])

	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	c.FileAttachment(localPath, filename)

	c.Status(http.StatusOK)
}

// Release API

func CreateRelease(c *gin.Context) {
	var input snapshot.Snapshot
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	fileID, err := release.CreateReleaseFromSnapshot(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

func DeleteRelease(c *gin.Context) {
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := release.DeleteRelease(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetReleases(c *gin.Context) {
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
	data, err := release.GetReleases(limit, filter, count)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostRelease(c *gin.Context) {
	var input release.Release
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	fileID := uuid.New().String()
	input.ID = fileID
	err := release.RegisterRelease(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

// Snapshot API

func CreateSnapshot(c *gin.Context) {
	var input model.Model
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	fileID, err := snapshot.CreateSnapshotFromModel(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

func DeleteSnapshot(c *gin.Context) {
	var input []string
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := snapshot.DeleteSnapshot(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func GetSnapshots(c *gin.Context) {
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
	data, err := snapshot.GetSnapshots(limit, filter, count)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func PostSnapshot(c *gin.Context) {
	var input snapshot.Snapshot
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	fileID := uuid.New().String()
	input.ID = fileID
	err := snapshot.RegisterSnapshot(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

func PutSnapshot(c *gin.Context) {
	var input snapshot.Snapshot
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	err := snapshot.UpdateSnapshot(input)
	if err != nil {
		utils.Error(err, c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
