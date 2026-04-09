package middleware

import (
	"log/slog"
	"time"
	"zmd-gacha/internal/utils"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traceID, err := utils.GenerateRefreshToken(32)
		c.Set("trace_id", traceID)
		if err != nil {
			return err
		}
		start := time.Now()
		slog.Info("请求开始",
			slog.String("trace_id", traceID),
			slog.String("请求方法", c.Request().Method),
			slog.String("请求路径", c.Request().URL.Path),
			slog.String("客户端IP", c.RealIP()),
		)
		err = next(c)
		duration := time.Since(start).Milliseconds()
		status := c.Response().Status
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}
		slog.Log(c.Request().Context(), level, "请求结束",
			slog.String("trace_id", traceID),
			slog.Int("状态码", status),
			slog.Int64("耗时(ms)", duration),
		)

		return err
	}
}
