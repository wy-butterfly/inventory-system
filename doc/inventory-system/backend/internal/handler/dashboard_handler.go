package handler

import (
	"inventory-system/internal/repository"
	"inventory-system/pkg/errcode"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	inventoryRepo *repository.InventoryRepo
	logRepo       *repository.LogRepo
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		inventoryRepo: repository.NewInventoryRepo(),
		logRepo:       repository.NewLogRepo(),
	}
}

// GetDashboard 获取仪表盘数据
// GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	total, categoryStats, warnings, err := h.inventoryRepo.GetDashboardStats()
	if err != nil {
		apiResp.Error(c, errcode.ErrDatabase)
		return
	}

	recentOps, _ := h.logRepo.GetRecentOperations(10)

	apiResp.Success(c, gin.H{
		"total_count":    total,
		"category_stats": categoryStats,
		"warnings":       warnings,
		"recent_logs":    recentOps,
	})
}

// HealthCheck 健康检查
// GET /api/health
func (h *DashboardHandler) HealthCheck(c *gin.Context) {
	apiResp.Success(c, gin.H{
		"status": "ok",
	})
}
