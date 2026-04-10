package handler

import (
	"net/http"
	"strconv"
	"zmd-gacha/internal/game/service"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

type GachaHandler struct {
	Service *service.GachaService
}

func NewGachaHandler(service *service.GachaService) *GachaHandler {
	return &GachaHandler{Service: service}
}

func (h *GachaHandler) Pull(c echo.Context) error {
	poolId := c.QueryParam("pool_id")
	times := c.QueryParam("times")
	if poolId == "" {
		return c.JSON(400, types.ErrorRsp{Message: "缺少pool_id参数"})
	}

	poolIdUint, err := strconv.ParseUint(poolId, 10, 64)
	if err != nil {
		return c.JSON(400, types.ErrorRsp{Message: "pool_id参数无效"})
	}

	uid := c.Get("uid").(uint)
	var result models.Character
	var results []models.Character
	switch times {
	case "1":
		result, err = h.Service.PullOnce(uint(poolIdUint), uid)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, types.PullOnceRsp{
			Character: result,
			Message:   "抽卡成功",
		})
	case "10":
		results, err = h.Service.PullTen(uint(poolIdUint), uid)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, types.PullTenRsp{
			Characters: results,
			Message:    "抽卡成功",
		})
	default:
		return c.JSON(400, types.ErrorRsp{Message: "times参数无效"})
	}
}

func (h *GachaHandler) GetPoolInfo(c echo.Context) error {
	poolId := c.QueryParam("pool_id")
	if poolId == "" {
		return c.JSON(400, types.ErrorRsp{Message: "缺少pool_id参数"})
	}
	poolIdUint, err := strconv.ParseUint(poolId, 10, 64)
	if err != nil {
		return c.JSON(400, types.ErrorRsp{Message: "pool_id参数无效"})
	}
	poolInfo, err := h.Service.GetPoolInfo(uint(poolIdUint))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.PoolInfoRsp{
		Pool:    poolInfo,
		Message: "获取卡池信息成功",
	})
}
