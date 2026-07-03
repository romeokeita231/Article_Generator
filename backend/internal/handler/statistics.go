package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/romeokeita231/Article_Generator/internal/service"
	"github.com/romeokeita231/Article_Generator/internal/common"
	_ "github.com/romeokeita231/Article_Generator/internal/model"
	
)
// StatisticsHandler 统计分析处理器
type StatisticsHandler struct {
    statisticsService *service.StatisticsService
}

// NewStatisticsHandler 创建统计分析处理器
func NewStatisticsHandler(statisticsService *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

// GetStatistics 
// @Summary 获取系统统计数据（仅管理员）
// @Tags statisticsHandler
// @Accept json
// @Produce json
// @Success 200 {object} common.BaseResponse{data=model.StatisticsVO}
// @Router /statistics/overview [get]
func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
    statistics, err := h.statisticsService.GetStatistics()
    if err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrSystem.WithMessage("获取统计数据失败: "+err.Error())))
        return
    }
    c.JSON(http.StatusOK, common.Success(statistics))
}
