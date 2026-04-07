package handler

import (
	"strconv"
	"zmd-gacha/internal/service"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

type GachaHandler struct {
	Service *service.GachaService
}

func NewGachaHandler(service *service.GachaService) *GachaHandler {
	return &GachaHandler{Service: service}
}

func (h *GachaHandler) PullOnce(c echo.Context) error {
	poolId := c.QueryParam("pool_id")
	if poolId == "" {
		return c.JSON(400, types.ErrorRsp{Message: "缺少pool_id参数"})
	}

	poolIdUint, err := strconv.ParseUint(poolId, 10, 64)
	if err != nil {
		return c.JSON(400, types.ErrorRsp{Message: "pool_id参数无效"})
	}

	uid := c.Get("uid").(uint)

	result, err := h.Service.PullOnce(uint(poolIdUint), uid)
	if err != nil {
		return c.JSON(500, types.ErrorRsp{Message: "抽卡失败"})
	}

	return c.JSON(200, types.PullOnceRsp{
		Character: result,
		Message:   "抽卡成功",
	})
}
