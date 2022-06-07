// routes.go

package main

import (
	"api-server/api"
	"api-server/middleware"
)

func initializeRoutes() {
	router.GET("/health", api.GetHealth)
	router.GET("/status", api.GetStatuses)
	apiRoutes := router.Group("/api")
	{
		// middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"),
		// Asset
		apiRoutes.DELETE("/asset", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.DeleteAsset)
		apiRoutes.GET("/asset", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.GetAssets)
		apiRoutes.POST("/asset", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostAsset)
		apiRoutes.PUT("/asset", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PutAsset)
		apiRoutes.POST("/asset/file", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.CreateFileAsset)
		apiRoutes.GET("/asset/file/:id", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.DownloadFileAsset)
		apiRoutes.POST("/asset/git", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.CreateGitAsset)
		apiRoutes.POST("/asset/git/sync", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.SyncGitAsset)

		// Model
		apiRoutes.DELETE("/model", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.DeleteModel)
		apiRoutes.GET("/model", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.GetModels)
		apiRoutes.POST("/model", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostModel)
		apiRoutes.PUT("/model", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PutModel)
		apiRoutes.DELETE("/model/asset", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.DeleteModelAsset)
		apiRoutes.POST("/model/asset", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostModelAsset)
		apiRoutes.GET("/model/archive/:id", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.DownloadModel)

		// Param
		apiRoutes.DELETE("/param", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.DeleteParam)
		apiRoutes.GET("/param", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.GetParams)
		apiRoutes.POST("/param", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.PostParam)
		apiRoutes.PUT("/param", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.PutParam)

		// Release
		// apiRoutes.DELETE("/release", api.DeleteRelease)
		apiRoutes.GET("/release", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.GetReleases)
		apiRoutes.POST("/release", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostRelease)
		apiRoutes.POST("/release/snapshot", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostReleaseSnapshot)
		apiRoutes.GET("/release/archive/:id", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.DownloadRelease)
		// apiRoutes.PUT("/release", api.PutRelease)

		// Secret
		apiRoutes.DELETE("/secret", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.DeleteSecret)
		apiRoutes.GET("/secret", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.GetSecrets)
		apiRoutes.POST("/secret", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.PostSecret)
		apiRoutes.PUT("/secret", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.PutSecret)

		// Service
		apiRoutes.DELETE("/service", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.DeleteService)
		apiRoutes.GET("/service", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.GetServices)
		apiRoutes.POST("/service", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostService)
		apiRoutes.PUT("/service", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PutService)

		// Snapshot
		// apiRoutes.DELETE("/snapshot", api.DeleteSnapshot)
		apiRoutes.GET("/snapshot", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.GetSnapshots)
		apiRoutes.POST("/snapshot", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostSnapshot)
		apiRoutes.POST("/snapshot/model", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("write"), api.PostSnapshotModel)
		apiRoutes.GET("/snapshot/archive/:id", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("read"), api.DownloadSnapshot)
		// apiRoutes.PUT("/snapshot", api.PutSnapshot)

		// Secret
		apiRoutes.DELETE("/user", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.DeleteUser)
		apiRoutes.GET("/user", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.GetUsers)
		apiRoutes.POST("/user", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.PostUser)
		apiRoutes.PUT("/user", middleware.EnsureLoggedIn(), middleware.EnsureRoleAllowed("admin"), middleware.EnsureGroupAllowed("admin"), api.PutUser)
	}
}
