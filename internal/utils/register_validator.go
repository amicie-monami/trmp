package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterValidator(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()

			field = strings.ToLower(field)

			switch tag {
			case "required":
				errors[field] = "Это поле обязательно для заполнения"
			case "email":
				errors[field] = "Неверный формат email"
			case "min":
				if field == "password" {
					errors[field] = "Пароль должен содержать минимум 6 символов"
				} else if field == "name" {
					errors[field] = "Имя должно содержать минимум 2 символа"
				} else {
					errors[field] = "Значение слишком короткое"
				}
			default:
				errors[field] = "Некорректное значение"
			}
		}
	}

	return errors
}
