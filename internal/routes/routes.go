package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter, db *sql.DB) {
	r.Use()

	authHandler := NewAuthHandler(db)
	writersHandler := NewWritersHandler(db)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login())
			auth.POST("/register", authHandler.Register())
		}

		writers := api.Group("/writers")
		{
			writers.GET("", writersHandler.GetWriters())
			writers.GET("/:id/bio", writersHandler.GetWriterBiography())
		}

		// articles := api.Group("/articles")
		// {
		// 	articles.GET("", articlesHandler())
		// 	articles.GET("/:id", articleHandler())
		// }

		// search := api.Group("/search")
		// {
		// 	search.GET("/search", searchHandler())
		// }

		// favorites := api.Group("/favorites")
		// {
		// 	favorites.GET("favorites")
		// }

		// me := api.Group("me")
		// {
		// 	me.GET("/me")
		// }
	}
}
