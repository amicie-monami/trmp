package routes

import (
	"log"
	"net/http"
	"strings"
	"trmp/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("=== Login handler called ===")

		var req LoginRequest

		// Валидация входных данных
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Printf("Login validation error: %v", err)
			validationErrors := utils.RegisterValidator(err)

			if len(validationErrors) > 0 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error":   "Неверные данные",
					"details": validationErrors,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Неверный формат данных",
				})
			}
			return
		}

		// Нормализация email
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.Password = strings.TrimSpace(req.Password)

		log.Printf("Login attempt for email: %s", req.Email)

		// Поиск пользователя по email
		user, err := h.userRepo.FindByEmail(req.Email)
		if err != nil {
			log.Printf("Error finding user: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка сервера",
			})
			return
		}

		// Если пользователь не найден
		if user == nil {
			log.Printf("User not found: %s", req.Email)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный email или пароль",
			})
			return
		}

		// Проверка пароля
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			log.Printf("Invalid password for user: %s", req.Email)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный email или пароль",
			})
			return
		}

		log.Printf("Login successful for: %s (ID: %d)", user.Email, user.ID)

		// Генерация JWT токена
		token, err := utils.GenerateJWT(user.ID, user.Email)
		if err != nil {
			log.Printf("Error generating JWT: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при генерации токена",
			})
			return
		}

		// Скрываем пароль в ответе
		user.Password = ""

		ctx.JSON(http.StatusOK, AuthResponse{
			Token: token,
			User:  *user,
		})

		log.Printf("Login completed successfully for: %s", user.Email)
	}
}
