package controller

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

func parsePositiveIDParam(c echo.Context, name string) (int64, error) {
	param := c.Param(name)
	id, err := strconv.Atoi(param)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return int64(id), nil
}
