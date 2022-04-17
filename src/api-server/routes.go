// routes.go

package main

import (
	"api-server/api"
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

		// Model
		apiRoutes.DELETE("/model", api.DeleteModel)
		apiRoutes.GET("/model", api.GetModels)
		apiRoutes.POST("/model", api.PostModel)
		apiRoutes.PUT("/model", api.PutModel)

		// Param
		apiRoutes.DELETE("/param", api.DeleteParam)
		apiRoutes.GET("/param", api.GetParams)
		apiRoutes.POST("/param", api.PostParam)
		apiRoutes.PUT("/param", api.PutParam)

		// Release
		apiRoutes.DELETE("/release", api.DeleteRelease)
		apiRoutes.GET("/release", api.GetReleases)
		apiRoutes.POST("/release", api.PostRelease)
		apiRoutes.PUT("/release", api.PutRelease)

		// Secret
		apiRoutes.DELETE("/secret", api.DeleteSecret)
		apiRoutes.GET("/secret", api.GetSecrets)
		apiRoutes.POST("/secret", api.PostSecret)
		apiRoutes.PUT("/secret", api.PutSecret)

		// Secret
		apiRoutes.DELETE("/service", api.DeleteService)
		apiRoutes.GET("/service", api.GetServices)
		apiRoutes.POST("/service", api.PostService)
		apiRoutes.PUT("/service", api.PutService)

		// Snapshot
		apiRoutes.DELETE("/snapshot", api.DeleteSnapshot)
		apiRoutes.GET("/snapshot", api.GetSnapshots)
		apiRoutes.POST("/snapshot", api.PostSnapshot)
		apiRoutes.PUT("/snapshot", api.PutSnapshot)
	}
}