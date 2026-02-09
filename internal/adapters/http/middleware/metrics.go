package middleware

import (
	"strconv"
	"time"

	"product-app/internal/ports"

	"github.com/labstack/echo/v4"
)

func MetricsMiddleware(metrics ports.MetricsPort) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			status := c.Response().Status
			path := c.Path()
			method := c.Request().Method

			metrics.IncHTTPRequests(
				method,
				path,
				strconv.Itoa(status),
			)

			metrics.ObserveHTTPLatency(
				method,
				path,
				time.Since(start).Seconds(),
			)

			return err
		}
	}
}
