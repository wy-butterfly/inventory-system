package handler

import (
	"strconv"
	"time"

	"inventory-system/internal/repository"
	"inventory-system/pkg/errcode"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logRepo *repository.LogRepo
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		logRepo: repository.NewLogRepo(),
	}
}

// ListOperationLogs 查询操作日志
// GET /api/v1/logs/operations?action=create&start_time=2026-01-01&end_time=2026-12-31&page=1&page_size=20
func (h *LogHandler) ListOperationLogs(c *gin.Context) {
	action := c.Query("action")
	targetType := c.Query("target_type")

	page := 1
	pageSize := 20

	if p := c.DefaultQuery("page", "1"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.DefaultQuery("page_size", "20"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	var startTime, endTime *time.Time
	if st := c.Query("start_time"); st != "" {
		if t, err := time.Parse("2006-01-02", st); err == nil {
			startTime = &t
		}
	}
	if et := c.Query("end_time"); et != "" {
		if t, err := time.Parse("2006-01-02", et); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			endTime = &t
		}
	}

	offset := (page - 1) * pageSize
	logs, total, err := h.logRepo.ListOperationLogs(action, targetType, nil, startTime, endTime, offset, pageSize)
	if err != nil {
		apiResp.Error(c, errcode.ErrDatabase)
		return
	}

	apiResp.PageSuccess(c, logs, total, page, pageSize)
}

// ListChangeLogs 查询变更日志
// GET /api/v1/logs/changes?inventory_id=1&page=1&page_size=20
func (h *LogHandler) ListChangeLogs(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.DefaultQuery("page", "1"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.DefaultQuery("page_size", "20"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	var inventoryID *uint
	if idStr := c.Query("inventory_id"); idStr != "" {
		if v, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			id := uint(v)
			inventoryID = &id
		}
	}

	var startTime, endTime *time.Time
	if st := c.Query("start_time"); st != "" {
		if t, err := time.Parse("2006-01-02", st); err == nil {
			startTime = &t
		}
	}
	if et := c.Query("end_time"); et != "" {
		if t, err := time.Parse("2006-01-02", et); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			endTime = &t
		}
	}

	offset := (page - 1) * pageSize
	logs, total, err := h.logRepo.ListChangeLogs(inventoryID, nil, startTime, endTime, offset, pageSize)
	if err != nil {
		apiResp.Error(c, errcode.ErrDatabase)
		return
	}

	apiResp.PageSuccess(c, logs, total, page, pageSize)
}
