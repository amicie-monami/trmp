package routes

import (
	"database/sql"
	"log"
	"net/http"
	"trmp/internal/database/repository"
	"trmp/internal/model"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	progressRepo *repository.ProgressRepository
}

func NewProgressHandler(db *sql.DB) *ProgressHandler {
	return &ProgressHandler{
		progressRepo: repository.NewProgressRepository(db),
	}
}

func (h *ProgressHandler) GetReadingProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)

		progress, err := h.progressRepo.GetUserProgress(userID)
		if err != nil {
			log.Printf("Error getting reading progress: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении прогресса",
			})
			return
		}

		ctx.JSON(http.StatusOK, progress)
	}
}

func (h *ProgressHandler) SaveReadingProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)

		var req model.ProgressUpdateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Неверный формат данных",
			})
			return
		}

		if req.Progress < 0 || req.Progress > 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Прогресс должен быть от 0.0 до 1.0",
			})
			return
		}

		err := h.progressRepo.UpdateProgress(userID, req.Type, req.ID, req.Progress)
		if err != nil {
			log.Printf("Error saving progress: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при сохранении прогресса",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":  "Прогресс сохранен",
			"type":     req.Type,
			"id":       req.ID,
			"progress": req.Progress,
		})
	}
}

func (h *ProgressHandler) BulkSaveReadingProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, _ := GetUserIDFromContext(ctx)

		var req model.BulkProgressRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Неверный формат данных",
			})
			return
		}

		for id, progress := range req.Writers {
			if progress < 0 || progress > 1 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Прогресс для писателя " + id + " должен быть от 0.0 до 1.0",
				})
				return
			}
		}

		for id, progress := range req.Articles {
			if progress < 0 || progress > 1 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Прогресс для статьи " + id + " должен быть от 0.0 до 1.0",
				})
				return
			}
		}

		err := h.progressRepo.BulkUpdateProgress(userID, req.Writers, req.Articles)
		if err != nil {
			log.Printf("Error bulk saving progress: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при массовом сохранении прогресса",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Прогресс массово сохранен",
			"counts": gin.H{
				"writers":  len(req.Writers),
				"articles": len(req.Articles),
			},
		})
	}
}
