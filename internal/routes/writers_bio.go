package routes

import (
	"net/http"
	"strconv"
	"trmp/internal/model"

	"github.com/gin-gonic/gin"
)

func writersBiographyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"error": "ID must be a number",
			})
			return
		}

		wb := model.WriterBiography{
			ID:          id,
			Name:        "Sergey",
			PortraitURL: "test.url",
			Lifespan:    "01.02.1900-03.09.1956",
			Country:     "Russia",
			Occuptation: "Poet",
			IsFavorite:  false,
			Content:     "content text",
		}

		ctx.JSON(http.StatusOK, wb)
		return
	}
}
