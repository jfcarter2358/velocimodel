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

	apiRoutes := router.Group("/script/api")
	{
		// apiRoutes.GET("/asset", api.GetAssets)
		// apiRoutes.GET("/asset/:id", api.GetAsset)
		apiRoutes.GET("/model", api.GetModels)
		apiRoutes.GET("/model/:id", api.GetModel)
		// apiRoutes.GET("/release", api.GetReleases)
		// apiRoutes.GET("/release/:id", api.GetRelease)
		// apiRoutes.GET("/snapshot", api.GetModels)
		// apiRoutes.GET("/snapshot/:id", api.GetModel)
	}

	uiRoutes := router.Group("/ui")
	{
		uiRoutes.GET("/assets", page.ShowAssetsPage)
		uiRoutes.GET("/asset/:id", page.ShowAssetPage)
		uiRoutes.GET("/asset/:id/edit", page.ShowAssetEditPage)
		uiRoutes.GET("/dashboard", page.ShowDashboardPage)
		uiRoutes.GET("/models", page.ShowModelsPage)
		uiRoutes.GET("/model/:id", page.ShowModelPage)
		uiRoutes.GET("/model/:id/edit", page.ShowModelEditPage)
		uiRoutes.GET("/releases", page.ShowReleasesPage)
		uiRoutes.GET("/release/:id", page.ShowReleasePage)
		uiRoutes.GET("/release/:id/edit", page.ShowReleaseEditPage)
		uiRoutes.GET("/snapshots", page.ShowSnapshotsPage)
		uiRoutes.GET("/snapshot/:id", page.ShowSnapshotPage)
		uiRoutes.GET("/snapshot/:id/edit", page.ShowSnapshotEditPage)
	}
}
