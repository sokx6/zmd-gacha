package handler

import (
	"net/http"
	"zmd-gacha/internal/service"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	var req types.ProfileUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ProfileUpdateRsp{
			Message: "不合法请求",
		})
	}

	if err := h.Service.UpdateProfile(req.User); err != nil {
		return c.JSON(http.StatusInternalServerError, types.ProfileUpdateRsp{
			Message: "更新失败",
		})
	}

	return c.JSON(http.StatusOK, types.ProfileUpdateRsp{
		Message: "更新成功",
	})

}
