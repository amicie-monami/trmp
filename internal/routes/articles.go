package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"trmp/internal/database/repository"
	"trmp/internal/model"

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

// GetArticles возвращает список всех статей (карточки)
func (h *ArticlesHandler) GetArticles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Getting all articles...")

		articles, err := h.articleRepo.GetAll()
		if err != nil {
			log.Printf("Error getting articles: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении списка статей",
			})
			return
		}

		// Если массив пустой, возвращаем пустой массив
		if articles == nil {
			ctx.JSON(http.StatusOK, []model.ArticleCard{})
			return
		}

		ctx.JSON(http.StatusOK, articles)
		log.Printf("Successfully returned %d articles", len(articles))
	}
}

// GetArticle возвращает полную статью по ID
func (h *ArticlesHandler) GetArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		log.Printf("Getting article with ID: %s", idStr)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID must be a number",
			})
			return
		}

		article, err := h.articleRepo.GetByID(id)
		if err != nil {
			log.Printf("Error getting article: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении статьи",
			})
			return
		}

		if article == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Статья не найдена",
			})
			return
		}

		ctx.JSON(http.StatusOK, article)
		log.Printf("Successfully returned article: %s (ID: %d)", article.Title, article.ID)
	}
}
