// routes.go

package main

import (
	"service-manager/api"
)

func initializeRoutes() {
	router.GET("/health", api.GetHealth)
	router.GET("/status", api.GetStatuses)
	apiRoutes := router.Group("/api")
	{
		// Service
		apiRoutes.DELETE("/service", api.DeleteService)
		apiRoutes.GET("/service", api.GetServices)
		apiRoutes.POST("/service", api.PostService)
		apiRoutes.PUT("/service", api.PutService)
		// Params
		apiRoutes.DELETE("/param", api.DeleteParam)
		apiRoutes.GET("/param", api.GetParams)
		apiRoutes.POST("/param", api.PostParam)
		apiRoutes.PUT("/param", api.PutParam)
		// Secrets
		apiRoutes.DELETE("/secret", api.DeleteSecret)
		apiRoutes.GET("/secret", api.GetSecrets)
		apiRoutes.POST("/secret", api.PostSecret)
		apiRoutes.PUT("/secret", api.PutSecret)
	}
}
