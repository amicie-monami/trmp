package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"trmp/internal/database/repository"
	"trmp/internal/model"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchRepo *repository.SearchRepository
}

func NewSearchHandler(db *sql.DB) *SearchHandler {
	return &SearchHandler{
		searchRepo: repository.NewSearchRepository(db),
	}
}

// SearchAllHandler общий поиск по статьям и писателям (раздельные результаты)
func (h *SearchHandler) SearchAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)

		// параметры поиска
		query := strings.TrimSpace(ctx.Query("q"))
		tagsParam := strings.TrimSpace(ctx.Query("tags"))

		// теги
		var tags []string
		if tagsParam != "" {
			tags = parseTagsParam(tagsParam)
		}

		if query == "" && len(tags) == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"articles": []model.ArticleCard{},
				"writers":  []model.WriterCard{},
				"counts": gin.H{
					"articles": 0,
					"writers":  0,
				},
				"query": query,
				"tags":  tags,
			})
			return
		}

		log.Printf("Search all: query='%s', tags=%v, userID=%d",
			query, tags, userID)

		articles, writers, err := h.searchRepo.SearchAll(query, tags, userID)
		if err != nil {
			log.Printf("Error searching: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при поиске",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"articles": articles,
			"writers":  writers,
			"counts": gin.H{
				"articles": len(articles),
				"writers":  len(writers),
			},
			"query": query,
			"tags":  tags,
		})
	}
}

func (h *SearchHandler) SearchArticles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)

		// параметры поиска
		query := strings.TrimSpace(ctx.Query("q"))
		tagsParam := strings.TrimSpace(ctx.Query("tags"))

		// теги (могут быть через запятую: "литература,история")
		var tags []string
		if tagsParam != "" {
			tags = parseTagsParam(tagsParam)
		}

		// нет ни запроса, ни тегов - возвращаем пустой результат
		if query == "" && len(tags) == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"articles": []model.ArticleCard{},
				"writers":  []model.WriterCard{},
				"counts": gin.H{
					"articles": 0,
					"writers":  0,
					"total":    0,
				},
				"query": query,
				"tags":  tags,
			})
			return
		}

		log.Printf("Search articles: query='%s', tags=%v, userID=%d",
			query, tags, userID)

		articles, err := h.searchRepo.SearchArticles(query, tags, userID)
		if err != nil {
			log.Printf("Error searching articles: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при поиске статей",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"articles": articles,
			"count":    len(articles),
			"query":    query,
			"tags":     tags,
		})
	}
}

func (h *SearchHandler) SearchWriters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)

		query := strings.TrimSpace(ctx.Query("q"))
		tagsParam := strings.TrimSpace(ctx.Query("tags"))

		var tags []string
		if tagsParam != "" {
			tags = parseTagsParam(tagsParam)
		}

		log.Printf("Search writers: query='%s', tags=%v, userID=%d",
			query, tags, userID)

		writers, err := h.searchRepo.SearchWriters(query, tags, userID)
		if err != nil {
			log.Printf("Error searching writers: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при поиске писателей",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"writers": writers,
			"count":   len(writers),
			"query":   query,
			"tags":    tags,
		})
	}
}

// GetTagsHandler возвращает все уникальные теги
func (h *SearchHandler) GetTags() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tags, err := h.searchRepo.GetAllTags()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении тегов",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"tags":  tags,
			"count": len(tags),
		})
	}
}

// Вспомогательная функция для парсинга тегов из параметра
func parseTagsParam(tagsParam string) []string {
	if tagsParam == "" {
		return []string{}
	}

	var tags []string
	// Разделяем по запятой
	rawTags := strings.Split(tagsParam, ",")
	for _, tag := range rawTags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}
