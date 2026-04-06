package handler

import (
	"net/http"
	"zmd-gacha/internal/service"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var req types.UserRstReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.Response{
			Message: "不合法请求",
		})
	} else if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, types.Response{
			Message: "用户名或密码不能为空",
		})
	} else if len(req.Username) > 40 || len(req.Password) > 40 {
		return c.JSON(http.StatusBadRequest, types.Response{
			Message: "用户名或密码过长",
		})
	} else if req.Email == "" {
		return c.JSON(http.StatusBadRequest, types.Response{
			Message: "邮箱不能为空",
		})
	}

	if err := service.Register(req.Username, req.Password, req.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, types.Response{
			Message: "注册失败",
		})
	}

	return c.JSON(http.StatusOK, types.Response{
		Message: "注册成功",
	})
}
