package handler

import (
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
	"github.com/gin-gonic/gin"
)

// StatsHandler 统计处理器。
type StatsHandler struct {
	svc *service.StatsService
}

func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

func (h *StatsHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats()
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, stats)
}
