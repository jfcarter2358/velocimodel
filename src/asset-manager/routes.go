// routes.go

package main

import (
	"asset-manager/api"
)

func initializeRoutes() {
	apiRoutes := router.Group("/api")
	{
		// Asset
		apiRoutes.POST("/asset/reload", api.ReloadAsset)
		apiRoutes.POST("/asset/sync", api.SyncAsset)
		apiRoutes.DELETE("/asset", api.DeleteAsset)
		apiRoutes.GET("/asset", api.GetAsset)
		apiRoutes.POST("/asset", api.PostAsset)
		apiRoutes.PUT("/asset", api.PutAsset)
		apiRoutes.POST("/asset/upload", api.UploadAsset)
	}
}
