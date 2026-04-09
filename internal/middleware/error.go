package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

func AppHTTPErrorHandler(err error, c echo.Context) {
	traceID := "-"
	if tID, ok := c.Get("trace_id").(string); ok {
		traceID = tID
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		slog.Warn("HTTP错误",
			slog.String("trace_id", traceID),
			slog.Int("状态码", he.Code),
			slog.String("错误信息", he.Message.(string)),
		)
		c.JSON(he.Code, map[string]interface{}{
			"code":    he.Code,
			"message": he.Message,
		})
		return
	}

	var appErr *types.AppError
	if errors.As(err, &appErr) {
		slog.Warn("应用错误",
			slog.String("trace_id", traceID),
			slog.Int("状态码", appErr.Code),
			slog.String("错误信息", appErr.Message),
			slog.Any("原始错误", appErr.Err),
		)
		c.JSON(appErr.Code, map[string]interface{}{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	slog.Error("未知错误",
		slog.String("trace_id", traceID),
		slog.Any("错误信息", err),
	)

	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code":    http.StatusInternalServerError,
		"message": "服务器内部错误",
	})
}
