package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"trmp/internal/database/repository"

	"github.com/gin-gonic/gin"
)

type ArticlesHandler struct {
	articleRepo *repository.ArticleRepository
}

func NewArticlesHandler(db *sql.DB) *ArticlesHandler {
	return &ArticlesHandler{
		articleRepo: repository.NewArticleRepository(db),
	}
}

func (h *ArticlesHandler) GetArticles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)

		articles, err := h.articleRepo.GetAllWithFavorites(userID)
		if err != nil {
			log.Printf("Error getting articles: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении списка статей",
			})
			return
		}

		ctx.JSON(http.StatusOK, articles)
	}
}

func (h *ArticlesHandler) GetArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)
		id, _ := strconv.Atoi(ctx.Param("id"))

		article, err := h.articleRepo.GetByIDWithFavorite(id, userID)
		if err != nil || article == nil {
			log.Printf("Error getting article: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении статьи",
			})
			return
		}

		ctx.JSON(http.StatusOK, article)
	}
}
