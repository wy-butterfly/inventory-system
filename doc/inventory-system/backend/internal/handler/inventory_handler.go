package handler

import (
	"strconv"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/service"
	"inventory-system/pkg/errcode"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
	exportService    *service.ExportService
}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{
		inventoryService: service.NewInventoryService(),
		exportService:    service.NewExportService(),
	}
}

// Create 新增库存
// POST /api/v1/inventories
func (h *InventoryHandler) Create(c *gin.Context) {
	var req request.CreateInventoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage(err.Error()))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	inv, errCode := h.inventoryService.Create(&req, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "新增成功", inv)
}

// GetByID 获取库存详情
// GET /api/v1/inventories/:id
func (h *InventoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的ID"))
		return
	}

	inv, errCode := h.inventoryService.GetByID(uint(id))
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.Success(c, inv)
}

// List 分页查询库存列表
// GET /api/v1/inventories?page=1&page_size=20&keyword=xxx
func (h *InventoryHandler) List(c *gin.Context) {
	var query request.InventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage(err.Error()))
		return
	}

	inventories, total, errCode := h.inventoryService.List(&query)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.PageSuccess(c, inventories, total, query.GetPage(), query.GetPageSize())
}

// Update 编辑库存
// PUT /api/v1/inventories/:id
func (h *InventoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的ID"))
		return
	}

	var req request.UpdateInventoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage(err.Error()))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.inventoryService.Update(uint(id), &req, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "编辑成功", nil)
}

// Delete 删除库存（软删除）
// DELETE /api/v1/inventories/:id
func (h *InventoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的ID"))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.inventoryService.Delete(uint(id), userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "删除成功", nil)
}

// BatchDelete 批量删除
// POST /api/v1/inventories/batch-delete
func (h *InventoryHandler) BatchDelete(c *gin.Context) {
	var req request.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("请选择要删除的记录"))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.inventoryService.BatchDelete(req.IDs, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "批量删除成功", nil)
}

// Restore 恢复已删除的记录
// PUT /api/v1/inventories/:id/restore
func (h *InventoryHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的ID"))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.inventoryService.Restore(uint(id), userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "恢复成功", nil)
}

// Export 导出库存数据
// GET /api/v1/inventories/export
func (h *InventoryHandler) Export(c *gin.Context) {
	var query request.InventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage(err.Error()))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	excelFile, errCode := h.exportService.ExportInventories(&query, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=inventory_export.xlsx")

	if err := excelFile.Write(c.Writer); err != nil {
		apiResp.Error(c, errcode.ErrServer)
		return
	}
}
