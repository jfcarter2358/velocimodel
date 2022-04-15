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
		// Release
		apiRoutes.DELETE("/release", api.DeleteRelease)
		apiRoutes.GET("/release", api.GetReleases)
		apiRoutes.POST("/release", api.PostRelease)
		// Snapshot
		apiRoutes.DELETE("/snapshot", api.DeleteSnapshot)
		apiRoutes.GET("/snapshot", api.GetSnapshots)
		apiRoutes.POST("/snapshot", api.PostSnapshot)
		apiRoutes.PUT("/snapshot", api.PutSnapshot)
	}
}
