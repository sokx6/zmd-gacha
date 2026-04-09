package middleware

import (
	"errors"
	"net/http"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

func AppHTTPErrorHandler(err error, c echo.Context) {
	var he *echo.HTTPError
	if errors.As(err, &he) {
		c.JSON(he.Code, map[string]interface{}{
			"error": he.Message,
		})
	}

	var appErr *types.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, map[string]interface{}{
			"code":  appErr.Code,
			"error": appErr.Message,
		})
	}
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code":  http.StatusInternalServerError,
		"error": "服务器内部错误",
	})
}
