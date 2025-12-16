package routes

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	log.Println("Setting up routes...")

	authHandler := NewAuthHandler(db)
	writersHandler := NewWritersHandler(db)
	articlesHandler := NewArticlesHandler(db)
	favoritesHandler := NewFavoritesHandler(db)
	searchHandler := NewSearchHandler(db)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register())
			auth.POST("/login", authHandler.Login())
		}

		writers := api.Group("/writers")
		writers.Use(AuthMiddleware())
		{
			writers.GET("", writersHandler.GetWriters())
			writers.GET("/:id/bio", writersHandler.GetWriterBiography())
		}

		articles := api.Group("/articles")
		articles.Use(AuthMiddleware())
		{
			articles.GET("", articlesHandler.GetArticles())
			articles.GET("/:id", articlesHandler.GetArticle())
		}

		favorites := api.Group("/favorites")
		favorites.Use(AuthMiddleware())
		{
			//writers
			favorites.GET("/writers", favoritesHandler.GetFavoriteWriters())
			favorites.POST("/writers/:id/toggle", favoritesHandler.ToggleWriterFavorite())

			//articles
			favorites.GET("/articles", favoritesHandler.GetFavoriteArticles())
			favorites.POST("/articles/:id/toggle", favoritesHandler.ToggleArticleFavorite())
		}

		// Search routes
		search := api.Group("/search")
		search.Use(AuthMiddleware())
		{
			search.GET("", searchHandler.SearchAll())
			search.GET("/articles", searchHandler.SearchArticles()) // ?q=текст&tags=тег1,тег2
			search.GET("/writers", searchHandler.SearchWriters())   // ?q=текст&tags=тег1,тег2
			search.GET("/tags", searchHandler.GetTags())            // все уникальные теги
		}
	}

	log.Println("Routes setup complete")
}

// 		// search := api.Group("/search")
// 		// {
// 		// 	search.GET("/search", searchHandler())
// 		// }

// 		// favorites := api.Group("/favorites")
// 		// {
// 		// 	favorites.GET("favorites")
// 		// }

// 		// me := api.Group("me")
// 		// {
// 		// 	me.GET("/me")
// 		// }
// 	}
// }
