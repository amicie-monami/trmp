// routes/writers.go
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

type WritersHandler struct {
	writerRepo *repository.WriterRepository
}

func NewWritersHandler(db *sql.DB) *WritersHandler {
	return &WritersHandler{
		writerRepo: repository.NewWriterRepository(db),
	}
}

func (h *WritersHandler) GetWriters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Getting all writers...")

		var writers []model.WriterCard
		var err error

		writers, err = h.writerRepo.GetAll()

		if err != nil {
			log.Printf("Error getting writers: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении списка писателей",
			})
			return
		}

		// Если нет писателей
		if len(writers) == 0 {
			ctx.JSON(http.StatusOK, []model.WriterCard{})
			return
		}

		ctx.JSON(http.StatusOK, writers)
		log.Printf("Returned %d writers", len(writers))
	}
}

func (h *WritersHandler) GetWriterBiography() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		log.Printf("Getting biography for writer ID: %s", idStr)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID must be a number",
			})
			return
		}

		writer, err := h.writerRepo.GetByID(id)
		if err != nil {
			log.Printf("Error getting writer: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении биографии",
			})
			return
		}

		if writer == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Писатель не найден",
			})
			return
		}

		ctx.JSON(http.StatusOK, writer)
		log.Printf("Returned biography for writer: %s", writer.Name)
	}
}

func (h *WritersHandler) ToggleFavorite() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID must be a number",
			})
			return
		}

		err = h.writerRepo.ToggleFavorite(id)
		if err != nil {
			log.Printf("Error toggling favorite: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при обновлении",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Статус избранного обновлен",
		})
	}
}
