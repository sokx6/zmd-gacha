package handler

import (
	"net/http"
	"strconv"
	"zmd-gacha/internal/management/service"
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
	uid, ok := c.Get("uid").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, types.ProfileUpdateRsp{
			Message: "未授权",
		})
	}
	if err := h.Service.UpdateProfile(uid, req.Nickname, req.Profile); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.ProfileUpdateRsp{
		Code:    http.StatusOK,
		Message: "更新成功",
	})

}

func (h *UserHandler) GetUserCharacters(c echo.Context) error {
	uid := c.Get("uid").(uint)

	characters, err := h.Service.GetUserCharacters(uid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.NewCharsGetRsp(characters, http.StatusOK, "获取角色列表成功"))
}

func (h *UserHandler) GetCharFirstInfo(c echo.Context) error {
	uid := c.Get("uid").(uint)
	characterId := c.QueryParam("character_id")
	if u64, err := strconv.ParseUint(characterId, 10, 64); err != nil {
		return c.JSON(http.StatusBadRequest, types.CharFirstInfoRsp{
			Message: "不合法的角色ID",
		})
	} else {
		charId := uint(u64)
		firstAcquiredAt, firstAcquiredPool, firstAcquiredPullCount, err := h.Service.GetCharFirstInfo(uid, charId)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, types.CharFirstInfoRsp{
			Code:                   http.StatusOK,
			Message:                "获取角色首次信息成功",
			FirstAcquiredAt:        firstAcquiredAt,
			FirstAcquiredPool:      firstAcquiredPool,
			FirstAcquiredPullCount: firstAcquiredPullCount,
		})
	}
}
