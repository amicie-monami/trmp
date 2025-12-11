package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func searchHandler() gin.HandlerFunc {
	result := map[string][]any{
		"writers": {
			// wr1, wr2,
		},
		"articles": {},
	}

	return func(ctx *gin.Context) {
		like := ctx.Query("like")
		if like == "" {
			ctx.JSON(http.StatusOK, map[string][]any{"writers": {"nul"}, "articles": {"nul"}})
			return
		}

		// wr1 := model.WriterCard{ID: 1}
		// wr2 := model.WriterCard{ID: 2}
		// ar1 := model.Article{ID: 3}
		log.Println(like)

		ctx.JSON(http.StatusOK, result)
		return
	}
}
