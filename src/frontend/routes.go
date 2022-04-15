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
		apiRoutes.DELETE("/:path", api.DoDelete)
		apiRoutes.GET("/:path", api.DoGet)
		apiRoutes.POST("/:path", api.DoPost)
		apiRoutes.PUT("/:path", api.DoPut)
		apiRoutes.POST("/:path/upload", api.DoUpload)
	}

	uiRoutes := router.Group("/ui")
	{
		uiRoutes.GET("/dashboard", page.ShowHomePage)
	}
}
