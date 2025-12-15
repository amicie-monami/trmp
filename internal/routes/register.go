package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"trmp/internal/database/repository"
	"trmp/internal/model"
	"trmp/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(db),
	}
}

func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RegisterRequest

		// Валидация входных данных
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Printf("Validation error: %v", err)
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
		// Нормализация данных
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.Name = strings.TrimSpace(req.Name)
		req.Password = strings.TrimSpace(req.Password)

		// Проверка существования пользователя
		existingUser, err := h.userRepo.FindByEmail(req.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при проверке пользователя",
				"msg":   err.Error(),
			})
			return
		}

		if existingUser != nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Пользователь с таким email уже существует",
			})
			return
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при создании пользователя",
			})
			return
		}

		// Создание пользователя
		user := &model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		if err := h.userRepo.CreateUser(user); err != nil {
			// Проверяем, если это ошибка уникальности email
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "Пользователь с таким email уже существует",
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка при сохранении пользователя",
				})
			}
			return
		}

		// Получаем полные данные пользователя (опционально, можно использовать тот же объект)
		fullUser, err := h.userRepo.GetUserByID(user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при получении данных пользователя",
			})
			return
		}

		// Генерация JWT токена
		token, err := utils.GenerateJWT(fullUser.ID, fullUser.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при генерации токена",
			})
			return
		}

		// Скрываем пароль в ответе
		fullUser.Password = ""

		ctx.JSON(http.StatusCreated, AuthResponse{
			Token: token,
			User:  *fullUser,
		})
	}
}
