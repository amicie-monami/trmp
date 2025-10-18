package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
	r.Use()

	api := r.Group("/api")
	{
		api.GET("/writers", writersHandler())
		api.GET("/writers/:id", writersBiographyHandler())
		api.GET("/articles", articlesHandler())
		// api.GET("/search")
		// api.GET("favorites")
		// api.GET("/me")
	}
}
