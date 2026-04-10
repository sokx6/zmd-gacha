package handler

import (
	"net/http"
	"strconv"
	"zmd-gacha/internal/game/service"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) GetUserCharacters(c echo.Context) error {
	uid := c.Get("uid").(uint)

	characters, err := h.Service.GetUserCharacters(uid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.CharsGetRsp{
		Message:    "获取角色列表成功",
		Characters: characters,
	})
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
			Message:                "获取角色首次信息成功",
			FirstAcquiredAt:        firstAcquiredAt,
			FirstAcquiredPool:      firstAcquiredPool,
			FirstAcquiredPullCount: firstAcquiredPullCount,
		})
	}
}
