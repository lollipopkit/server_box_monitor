package web

import "github.com/labstack/echo/v4"

func ok(c echo.Context, data any) error {
	return c.JSON(200, map[string]any{
		"code": respCodeOK,
		"data": data,
	})
}

func fail(c echo.Context, code int, data any) error {
	return c.JSON(500, map[string]any{
		"code": code,
		"msg":  data,
	})
}
