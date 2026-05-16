package repository

import (
	"inventory-system/internal/model"
	"inventory-system/pkg/database"
)

type AttachmentRepo struct{}

func NewAttachmentRepo() *AttachmentRepo {
	return &AttachmentRepo{}
}

// Create 保存附件记录
func (r *AttachmentRepo) Create(att *model.Attachment) error {
	return database.DB.Create(att).Error
}

// FindByID 根据ID查找附件
func (r *AttachmentRepo) FindByID(id uint) (*model.Attachment, error) {
	var att model.Attachment
	err := database.DB.Where("id = ? AND is_deleted = 0", id).First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

// ListByInventoryID 查询某个库存记录的所有附件
func (r *AttachmentRepo) ListByInventoryID(inventoryID uint) ([]model.Attachment, error) {
	var attachments []model.Attachment
	err := database.DB.Where("inventory_id = ? AND is_deleted = 0", inventoryID).
		Order("created_at DESC").Find(&attachments).Error
	return attachments, err
}

// UpdateDisplayName 更新显示名称（重命名）
func (r *AttachmentRepo) UpdateDisplayName(id uint, displayName string) error {
	return database.DB.Model(&model.Attachment{}).
		Where("id = ?", id).
		Update("display_name", displayName).Error
}

// SoftDelete 软删除附件
func (r *AttachmentRepo) SoftDelete(id uint) error {
	return database.DB.Model(&model.Attachment{}).
		Where("id = ?", id).
		Update("is_deleted", 1).Error
}

// SoftDeleteByInventoryID 根据库存ID软删除所有附件
func (r *AttachmentRepo) SoftDeleteByInventoryID(inventoryID uint) error {
	return database.DB.Model(&model.Attachment{}).
		Where("inventory_id = ? AND is_deleted = 0", inventoryID).
		Update("is_deleted", 1).Error
}

// ListAll 分页查询所有附件
func (r *AttachmentRepo) ListAll(fileType string, inventoryID *uint,
	offset, limit int) ([]model.Attachment, int64, error) {

	var attachments []model.Attachment
	var total int64

	query := database.DB.Model(&model.Attachment{}).Where("is_deleted = 0")

	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}
	if inventoryID != nil {
		query = query.Where("inventory_id = ?", *inventoryID)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&attachments).Error
	return attachments, total, err
}
