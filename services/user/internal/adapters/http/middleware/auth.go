package middleware

import (
	"product-app/shared/auth"

	"github.com/labstack/echo/v4"
)

func GenerateToken(userId int64, username, email string) (string, error) {
	return auth.GenerateToken(userId, username, email)
}

func JWTMiddleware() echo.MiddlewareFunc {
	return auth.JWTMiddleware()
}
