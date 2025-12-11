package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter, db *sql.DB) {
	r.Use()

	api := r.Group("/api")
	{
		api.GET("/writers", writersHandler())
		api.GET("/writers/:id", writersBiographyHandler())
		api.GET("/articles", articlesHandler())
		api.GET("/articles/:id", articleHandler())
		api.GET("/search", searchHandler())
		// api.GET("favorites")
		// api.GET("/me")
	}
}
