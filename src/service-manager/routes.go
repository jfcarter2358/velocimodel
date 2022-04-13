// routes.go

package main

import (
	"service-manager/api"
)

func initializeRoutes() {
	apiRoutes := router.Group("/api")
	{
		// Service
		apiRoutes.POST("/service/reload", api.ReloadService)
		apiRoutes.POST("/service/sync", api.SyncService)
		apiRoutes.DELETE("/service", api.DeleteService)
		apiRoutes.GET("/service", api.GetService)
		apiRoutes.POST("/service", api.PostService)
		apiRoutes.PUT("/service", api.PutService)
		// Params
		apiRoutes.DELETE("/param", api.DeleteParam)
		apiRoutes.GET("/param", api.GetParam)
		apiRoutes.POST("/param", api.PostParam)
		apiRoutes.PUT("/param", api.PutParam)
		// Secrets
		apiRoutes.DELETE("/secret", api.DeleteSecret)
		apiRoutes.GET("/secret", api.GetSecret)
		apiRoutes.POST("/secret", api.PostSecret)
		apiRoutes.PUT("/secret", api.PutSecret)
	}
}
