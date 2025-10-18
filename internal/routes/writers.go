package routes

import (
	"net/http"
	"trmp/internal/model"

	"github.com/gin-gonic/gin"
)

func writersHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		writer := model.WriterCard{
			ID:          1,
			Name:        "Sergey",
			PortraitURL: "test.url",
			IsFavorite:  false,
		}

		ctx.JSON(http.StatusOK, writer)
	}
}
