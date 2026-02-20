package middleware

import (
	"product-app/shared/auth"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return auth.JWTMiddleware()
}
