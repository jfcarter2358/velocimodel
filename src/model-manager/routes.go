// routes.go

package main

import (
	"model-manager/api"
)

func initializeRoutes() {
	router.GET("/health", api.GetHealth)
	apiRoutes := router.Group("/api")
	{
		// Model
		apiRoutes.DELETE("/model", api.DeleteModel)
		apiRoutes.GET("/model", api.GetModels)
		apiRoutes.POST("/model", api.PostModel)
		apiRoutes.PUT("/model", api.PutModel)
		apiRoutes.POST("/model/asset", api.AddAsset)
		apiRoutes.DELETE("/model/asset", api.DeleteAsset)
		apiRoutes.GET("/model/archive/:id", api.DownloadModel)
		// Release
		// apiRoutes.DELETE("/release", api.DeleteRelease)
		apiRoutes.GET("/release", api.GetReleases)
		apiRoutes.POST("/release", api.PostRelease)
		apiRoutes.POST("/release/snapshot", api.CreateRelease)
		apiRoutes.GET("/release/archive/:id", api.DownloadRelease)
		// Snapshot
		// apiRoutes.DELETE("/snapshot", api.DeleteSnapshot)
		apiRoutes.GET("/snapshot", api.GetSnapshots)
		apiRoutes.POST("/snapshot", api.PostSnapshot)
		// apiRoutes.PUT("/snapshot", api.PutSnapshot)
		apiRoutes.POST("/snapshot/model", api.CreateSnapshot)
		apiRoutes.GET("/snapshot/archive/:id", api.DownloadSnapshot)
	}
}
