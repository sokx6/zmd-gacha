package handler

import (
	"net/http"
	"zmd-gacha/internal/management/service"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

type GachaHandler struct {
	Service *service.GachaService
}

func NewGachaHandler(service *service.GachaService) *GachaHandler {
	return &GachaHandler{Service: service}
}

func (h *GachaHandler) CreateCharacter(c echo.Context) error {
	var req types.CharCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, types.ErrorRsp{Message: "请求参数无效"})
	}

	character, err := h.Service.CreateCharacter(req.Name, req.Rank, req.IsLimited, req.IsUp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.CharCreateRsp{
		Character: character,
		Message:   "创建角色成功",
	})
}

func (h *GachaHandler) CreatePool(c echo.Context) error {
	var req types.PoolCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, types.ErrorRsp{Message: "请求参数无效"})
	}
	poolId, err := h.Service.CreatePool(req.Pool, req.Config)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.PoolCreateRsp{
		PoolID:  poolId,
		Message: "创建卡池成功",
	})
}

func (h *GachaHandler) InsertCharacterToPool(c echo.Context) error {
	var req types.InsertCharReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, types.ErrorRsp{Message: "请求参数无效"})
	}
	err := h.Service.InsertCharacterToPool(req.PoolId, req.CharacterId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.InsertCharRsp{
		Message: "插入角色成功",
	})
}
