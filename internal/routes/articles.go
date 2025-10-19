package routes

import (
	"net/http"
	"strconv"
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

func articleHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"error": "ID must be a number",
			})
			return
		}

		articles := model.Article{
			ID: id, Title: "Заголовок", Description: "Описание",
			CoverURL: "test.url", IsFavorite: false, Content: "Текст",
		}

		ctx.JSON(http.StatusOK, articles)
		return
	}
}
