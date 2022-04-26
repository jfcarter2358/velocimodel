// routes.go

package main

import (
	"asset-manager/api"
)

func initializeRoutes() {
	router.GET("/health", api.GetHealth)
	apiRoutes := router.Group("/api")
	{
		// Asset
		apiRoutes.DELETE("/asset", api.DeleteAsset)
		apiRoutes.GET("/asset", api.GetAssets)
		apiRoutes.POST("/asset", api.PostAsset)
		apiRoutes.PUT("/asset", api.PutAsset)
		apiRoutes.POST("/asset/file", api.CreateFileAsset)
		apiRoutes.GET("/asset/file/:id", api.DownloadFileAsset)
		apiRoutes.POST("/asset/git", api.CreateGitAsset)
		apiRoutes.POST("/asset/git/sync", api.SyncGitAsset)
	}
}
