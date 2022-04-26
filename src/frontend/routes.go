// routes.go

package main

import (
	"frontend/api"
	"frontend/page"
)

func initializeRoutes() {
	router.Static("/static/css", "./static/css")
	router.Static("/static/img", "./static/img")
	router.Static("/static/js", "./static/js")

	router.GET("/", page.RedirectIndexPage)

	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/asset", api.GetAssets)
		apiRoutes.GET("/asset/:id", api.GetAsset)
		apiRoutes.PUT("/asset/:id", api.UpdateAsset)
		apiRoutes.POST("/asset/file", api.CreateFileAsset)
		apiRoutes.POST("/asset/git", api.CreateGitAsset)
		apiRoutes.POST("/asset/git/sync/:id", api.SyncGitAsset)
		apiRoutes.GET("/asset/file/:id", api.DownloadAsset)
		apiRoutes.GET("/model", api.GetModels)
		apiRoutes.GET("/model/:id", api.GetModel)
		apiRoutes.PUT("/model/:id", api.UpdateModel)
		apiRoutes.POST("/model/:id/snapshot", api.CreateSnapshot)
		apiRoutes.POST("/model/asset", api.ModelAddAsset)
		apiRoutes.DELETE("/model/asset", api.ModelDeleteAsset)
		apiRoutes.POST("/model", api.CreateNewModel)
		apiRoutes.GET("/model/archive/:id", api.DownloadModel)
		apiRoutes.GET("/release", api.GetReleases)
		apiRoutes.GET("/release/:id", api.GetRelease)
		apiRoutes.GET("/release/archive/:id", api.DownloadRelease)
		apiRoutes.GET("/snapshot", api.GetSnapshots)
		apiRoutes.GET("/snapshot/:id", api.GetSnapshot)
		apiRoutes.PUT("/snapshot/:id", api.UpdateSnapshot)
		apiRoutes.POST("/snapshot/:id/release", api.CreateRelease)
		apiRoutes.GET("/snapshot/archive/:id", api.DownloadSnapshot)
	}

	uiRoutes := router.Group("/ui")
	{
		uiRoutes.GET("/assets", page.ShowAssetsPage)
		uiRoutes.GET("/asset/:id", page.ShowAssetPage)
		uiRoutes.GET("/asset/:id/code", page.ShowAssetCodePage)
		uiRoutes.GET("/dashboard", page.ShowDashboardPage)
		uiRoutes.GET("/models", page.ShowModelsPage)
		uiRoutes.GET("/model/:id", page.ShowModelPage)
		uiRoutes.GET("/model/:id/code", page.ShowModelCodePage)
		uiRoutes.GET("/releases", page.ShowReleasesPage)
		uiRoutes.GET("/release/:id", page.ShowReleasePage)
		uiRoutes.GET("/release/:id/code", page.ShowReleaseCodePage)
		uiRoutes.GET("/snapshots", page.ShowSnapshotsPage)
		uiRoutes.GET("/snapshot/:id", page.ShowSnapshotPage)
		uiRoutes.GET("/snapshot/:id/code", page.ShowSnapshotCodePage)
	}
}
