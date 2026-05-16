package handler

import (
	"strconv"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/service"
	"inventory-system/pkg/errcode"
	apiResp "inventory-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	attachmentService *service.AttachmentService
}

func NewAttachmentHandler() *AttachmentHandler {
	return &AttachmentHandler{
		attachmentService: service.NewAttachmentService(),
	}
}

// Upload 上传附件
// POST /api/v1/inventories/:id/attachments
func (h *AttachmentHandler) Upload(c *gin.Context) {
	inventoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的库存ID"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("请选择要上传的文件"))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	att, errCode := h.attachmentService.Upload(uint(inventoryID), file, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "上传成功", att)
}

// List 获取库存的附件列表
// GET /api/v1/inventories/:id/attachments
func (h *AttachmentHandler) List(c *gin.Context) {
	inventoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的库存ID"))
		return
	}

	attachments, errCode := h.attachmentService.ListByInventoryID(uint(inventoryID))
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.Success(c, attachments)
}

// Download 下载附件
// GET /api/v1/attachments/:id/download
func (h *AttachmentHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的附件ID"))
		return
	}

	att, errCode := h.attachmentService.GetByID(uint(id))
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	c.FileAttachment(att.FilePath, att.DisplayName)
}

// Preview 预览图片
// GET /api/v1/attachments/:id/preview
func (h *AttachmentHandler) Preview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的附件ID"))
		return
	}

	att, errCode := h.attachmentService.GetByID(uint(id))
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	if c.Query("thumbnail") == "true" && att.ThumbnailPath != "" {
		c.File(att.ThumbnailPath)
		return
	}

	c.File(att.FilePath)
}

// Rename 重命名附件
// PUT /api/v1/attachments/:id
func (h *AttachmentHandler) Rename(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的附件ID"))
		return
	}

	var req request.RenameAttachmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("请输入新文件名"))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.attachmentService.Rename(uint(id), req.DisplayName, userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "重命名成功", nil)
}

// Delete 删除附件
// DELETE /api/v1/attachments/:id
func (h *AttachmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiResp.Error(c, errcode.ErrInvalidParam.WithMessage("无效的附件ID"))
		return
	}

	userID := c.GetUint("user_id")
	userName := c.GetString("username")

	errCode := h.attachmentService.Delete(uint(id), userID, userName)
	if errCode != nil {
		apiResp.Error(c, errCode)
		return
	}

	apiResp.SuccessWithMessage(c, "删除成功", nil)
}
