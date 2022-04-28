// routes.go

package main

import (
	"frontend/api"
	"frontend/auth"
	"frontend/middleware"
	"frontend/page"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initializeRoutes() {
	router.Static("/static/css", "./static/css")
	router.Static("/static/img", "./static/img")
	router.Static("/static/js", "./static/js")

	router.GET("/", page.RedirectIndexPage)

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

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
		uiRoutes.GET("/assets", middleware.EnsureLoggedIn(), page.ShowAssetsPage)
		uiRoutes.GET("/asset/:id", middleware.EnsureLoggedIn(), page.ShowAssetPage)
		uiRoutes.GET("/asset/:id/code", middleware.EnsureLoggedIn(), page.ShowAssetCodePage)
		uiRoutes.GET("/dashboard", middleware.EnsureLoggedIn(), page.ShowDashboardPage)
		uiRoutes.GET("/models", middleware.EnsureLoggedIn(), page.ShowModelsPage)
		uiRoutes.GET("/model/:id", middleware.EnsureLoggedIn(), page.ShowModelPage)
		uiRoutes.GET("/model/:id/code", middleware.EnsureLoggedIn(), page.ShowModelCodePage)
		uiRoutes.GET("/releases", middleware.EnsureLoggedIn(), page.ShowReleasesPage)
		uiRoutes.GET("/release/:id", middleware.EnsureLoggedIn(), page.ShowReleasePage)
		uiRoutes.GET("/release/:id/code", middleware.EnsureLoggedIn(), page.ShowReleaseCodePage)
		uiRoutes.GET("/snapshots", middleware.EnsureLoggedIn(), page.ShowSnapshotsPage)
		uiRoutes.GET("/snapshot/:id", middleware.EnsureLoggedIn(), page.ShowSnapshotPage)
		uiRoutes.GET("/snapshot/:id/code", middleware.EnsureLoggedIn(), page.ShowSnapshotCodePage)
	}

	authRoutes := router.Group("/auth")
	{
		authRoutes.GET("/redirect", auth.HandleRedirect)
		authRoutes.GET("/login", auth.HandleLogin)
		authRoutes.GET("/logout", auth.HandleLogout)
	}
}
