package routes

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ВРЕМЕННО для тестирования - всегда user_id = 1
		ctx.Set("user_id", uint(1))
		ctx.Set("user_email", "test@mail.ru")
		log.Println("DEBUG: Using hardcoded user ID 1")
		ctx.Next()
	}
}

// // AuthMiddleware извлекает и проверяет JWT токен
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		if authHeader == "" {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
// 			ctx.Abort()
// 			return
// 		}

// 		// Проверяем формат "Bearer <token>"
// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
// 			ctx.Abort()
// 			return
// 		}

// 		tokenString := parts[1]
// 		claims, err := utils.ValidateJWT(tokenString)
// 		if err != nil {
// 			log.Printf("Token validation error: %v", err)
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
// 			ctx.Abort()
// 			return
// 		}

// 		// Сохраняем данные пользователя в контекст
// 		ctx.Set("user_id", claims.UserID)
// 		ctx.Set("user_email", claims.Email)

// 		log.Printf("Authenticated user ID: %d, Email: %s", claims.UserID, claims.Email)
// 		ctx.Next()
// 	}
// }

// GetUserIDFromContext вспомогательная функция для получения user_id из контекста
func GetUserIDFromContext(ctx *gin.Context) (int, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, false
	}

	// Преобразуем uint в int
	userIDUint, ok := userID.(uint)
	if !ok {
		return 0, false
	}

	return int(userIDUint), true
}
