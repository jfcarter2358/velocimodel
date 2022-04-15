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
		apiRoutes.POST("/asset/upload", api.UploadAsset)
	}
}
