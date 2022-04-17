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
		apiRoutes.GET("/model", api.GetModels)
		apiRoutes.GET("/model/:id", api.GetModel)
	}

	uiRoutes := router.Group("/ui")
	{
		uiRoutes.GET("/assets", page.ShowAssetsPage)
		uiRoutes.GET("/dashboard", page.ShowDashboardPage)
		uiRoutes.GET("/models", page.ShowModelsPage)
		uiRoutes.GET("/model/:id", page.ShowModelPage)
		uiRoutes.GET("/model/:id/edit", page.ShowModelEditPage)
		uiRoutes.GET("/releases", page.ShowReleasesPage)
		uiRoutes.GET("/snapshots", page.ShowSnapshotsPage)
	}
}
