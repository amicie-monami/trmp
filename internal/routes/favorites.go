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

type FavoritesHandler struct {
	favoritesRepo *repository.FavoritesRepository
	writerRepo    *repository.WriterRepository
	articleRepo   *repository.ArticleRepository
}

func NewFavoritesHandler(db *sql.DB) *FavoritesHandler {
	return &FavoritesHandler{
		favoritesRepo: repository.NewFavoritesRepository(db),
		writerRepo:    repository.NewWriterRepository(db),
		articleRepo:   repository.NewArticleRepository(db),
	}
}

// ToggleWriterFavorite переключает избранное для писателя
func (h *FavoritesHandler) ToggleWriterFavorite() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Получаем user_id из контекста (из токена)
		userID, exists := GetUserIDFromContext(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Не авторизован",
			})
			return
		}

		// Получаем writer_id из параметра URL
		writerIDStr := ctx.Param("id")
		writerID, err := strconv.Atoi(writerIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID писателя должно быть числом",
			})
			return
		}

		// Проверяем существование писателя
		writer, err := h.writerRepo.GetByID(writerID)
		if err != nil || writer == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Писатель не найден",
			})
			return
		}

		// Переключаем статус избранного
		isFavorite, err := h.favoritesRepo.ToggleWriterFavorite(userID, writerID)
		if err != nil {
			log.Printf("Error toggling writer favorite: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при обновлении избранного",
			})
			return
		}

		status := "добавлен в избранное"
		if !isFavorite {
			status = "удален из избранного"
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":     "Писатель " + status,
			"writer_id":   writerID,
			"is_favorite": isFavorite,
		})

		log.Printf("User %d toggled favorite for writer %d: %v", userID, writerID, isFavorite)
	}
}

// GetFavoriteWriters возвращает избранных писателей пользователя
func (h *FavoritesHandler) GetFavoriteWriters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := GetUserIDFromContext(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Не авторизован",
			})
			return
		}

		// Получаем список избранных писателей
		writerIDs, err := h.favoritesRepo.GetFavoriteWriters(userID)
		if err != nil {
			log.Printf("Error getting favorite writers: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении избранных писателей",
			})
			return
		}

		// Получаем полную информацию о каждом писателе
		var writers []model.WriterCard
		for _, writerID := range writerIDs {
			writer, err := h.writerRepo.GetByID(writerID)
			if err != nil || writer == nil {
				continue
			}

			writers = append(writers, model.WriterCard{
				ID:          writer.ID,
				Name:        writer.Name,
				PortraitURL: writer.PortraitURL,
				Tags:        writer.Tags,
				IsFavorite:  true, // Все они избранные
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"writers": writers,
			"count":   len(writers),
		})

		log.Printf("User %d has %d favorite writers", userID, len(writers))
	}
}

// GetFavoriteArticles возвращает избранные статьи пользователя
func (h *FavoritesHandler) GetFavoriteArticles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := GetUserIDFromContext(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Не авторизован",
			})
			return
		}

		// Получаем список избранных статей
		articleIDs, err := h.favoritesRepo.GetFavoriteArticles(userID)
		if err != nil {
			log.Printf("Error getting favorite articles: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении избранных статей",
			})
			return
		}

		// Получаем полную информацию о каждой статье
		var articles []model.ArticleCard
		for _, articleID := range articleIDs {
			article, err := h.articleRepo.GetByID(articleID)
			if err != nil || article == nil {
				continue
			}

			articles = append(articles, model.ArticleCard{
				ID:          article.ID,
				CoverURL:    article.CoverURL,
				Title:       article.Title,
				Tags:        article.Tags,
				Description: article.Description,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"articles": articles,
			"count":    len(articles),
		})

		log.Printf("User %d has %d favorite articles", userID, len(articles))
	}
}

// ToggleArticleFavorite переключает избранное для статьи
func (h *FavoritesHandler) ToggleArticleFavorite() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)
		articleID, _ := strconv.Atoi(ctx.Param("id"))

		// Проверяем существование статьи
		article, err := h.articleRepo.GetByID(articleID)
		if err != nil || article == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Статья не найдена"})
			return
		}

		isFavorite, err := h.favoritesRepo.ToggleArticleFavorite(userID, articleID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении избранного"})
			return
		}

		status := "добавлена в избранное"
		if !isFavorite {
			status = "удалена из избранного"
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":     "Статья " + status,
			"article_id":  articleID,
			"is_favorite": isFavorite,
		})
	}
}
