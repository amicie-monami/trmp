package routes

import (
	"net/http"
	"trmp/internal/model"

	"github.com/gin-gonic/gin"
)

func articlesHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articles := []model.Article{
			{
				ID: 1, Title: "Заголовок", Description: "Описание",
				CoverURL: "test.url", IsFavorite: false, Content: "Текст",
			},
			{
				ID: 2, Title: "Заголовок", Description: "Описание",
				CoverURL: "test.url", IsFavorite: false, Content: "Текст",
			},
		}

		ctx.JSON(http.StatusOK, articles)
		return
	}
}
