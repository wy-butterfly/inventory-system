package repository

import (
	"inventory-system/internal/model"
	"inventory-system/pkg/database"

	"gorm.io/gorm"
)

// InventoryRepo 库存数据访问
type InventoryRepo struct{}

func NewInventoryRepo() *InventoryRepo {
	return &InventoryRepo{}
}

// FindByID 根据ID查找库存（排除已删除的）
func (r *InventoryRepo) FindByID(id uint) (*model.Inventory, error) {
	var inv model.Inventory
	err := database.DB.Where("id = ? AND is_deleted = 0", id).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// FindByIDWithDeleted 根据ID查找库存（包含已删除的）
func (r *InventoryRepo) FindByIDWithDeleted(id uint) (*model.Inventory, error) {
	var inv model.Inventory
	err := database.DB.First(&inv, id).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// FindByMaterialCode 根据物料编码查找库存
func (r *InventoryRepo) FindByMaterialCode(code string) (*model.Inventory, error) {
	var inv model.Inventory
	err := database.DB.Where("material_code = ? AND is_deleted = 0", code).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// Create 新增库存记录
func (r *InventoryRepo) Create(inv *model.Inventory) error {
	return database.DB.Create(inv).Error
}

// Update 更新库存记录（带乐观锁）
func (r *InventoryRepo) Update(inv *model.Inventory, oldVersion int) error {
	result := database.DB.Model(inv).
		Where("id = ? AND version = ? AND is_deleted = 0", inv.ID, oldVersion).
		Updates(map[string]interface{}{
			"material_name": inv.MaterialName,
			"category":      inv.Category,
			"specification": inv.Specification,
			"unit":          inv.Unit,
			"quantity":      inv.Quantity,
			"safety_stock":  inv.SafetyStock,
			"location":      inv.Location,
			"status":        inv.Status,
			"remark":        inv.Remark,
			"version":       oldVersion + 1,
			"updated_by":    inv.UpdatedBy,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// SoftDelete 软删除
func (r *InventoryRepo) SoftDelete(id uint, operatorID uint) error {
	return database.DB.Model(&model.Inventory{}).
		Where("id = ? AND is_deleted = 0", id).
		Updates(map[string]interface{}{
			"is_deleted": 1,
			"updated_by": operatorID,
		}).Error
}

// BatchSoftDelete 批量软删除
func (r *InventoryRepo) BatchSoftDelete(ids []uint, operatorID uint) error {
	return database.DB.Model(&model.Inventory{}).
		Where("id IN ? AND is_deleted = 0", ids).
		Updates(map[string]interface{}{
			"is_deleted": 1,
			"updated_by": operatorID,
		}).Error
}

// Restore 恢复已删除的记录
func (r *InventoryRepo) Restore(id uint) error {
	return database.DB.Model(&model.Inventory{}).
		Where("id = ? AND is_deleted = 1", id).
		Update("is_deleted", 0).Error
}

// List 分页查询库存列表
func (r *InventoryRepo) List(keyword, category string, status *int8,
	minQty, maxQty *float64, location string, showDeleted bool,
	offset, limit int) ([]model.Inventory, int64, error) {

	var inventories []model.Inventory
	var total int64

	query := database.DB.Model(&model.Inventory{})

	if showDeleted {
		query = query.Where("is_deleted = 1")
	} else {
		query = query.Where("is_deleted = 0")
	}

	if keyword != "" {
		query = query.Where("(material_code LIKE ? OR material_name LIKE ?)",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if minQty != nil {
		query = query.Where("quantity >= ?", *minQty)
	}
	if maxQty != nil {
		query = query.Where("quantity <= ?", *maxQty)
	}
	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&inventories).Error

	return inventories, total, err
}

// ListAll 查询全部数据（用于导出，不分页）
func (r *InventoryRepo) ListAll(keyword, category string, status *int8,
	minQty, maxQty *float64) ([]model.Inventory, error) {

	var inventories []model.Inventory
	query := database.DB.Where("is_deleted = 0")

	if keyword != "" {
		query = query.Where("(material_code LIKE ? OR material_name LIKE ?)",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if minQty != nil {
		query = query.Where("quantity >= ?", *minQty)
	}
	if maxQty != nil {
		query = query.Where("quantity <= ?", *maxQty)
	}

	err := query.Order("updated_at DESC").Limit(10000).Find(&inventories).Error
	return inventories, err
}

// GetDashboardStats 获取仪表盘统计数据
func (r *InventoryRepo) GetDashboardStats() (int64, []map[string]interface{}, []model.Inventory, error) {
	var total int64
	database.DB.Model(&model.Inventory{}).Where("is_deleted = 0 AND status = 1").Count(&total)

	var categoryStats []map[string]interface{}
	database.DB.Model(&model.Inventory{}).
		Select("category, COUNT(*) as count, SUM(quantity) as total_quantity").
		Where("is_deleted = 0 AND status = 1").
		Group("category").
		Find(&categoryStats)

	var warnings []model.Inventory
	database.DB.Where("is_deleted = 0 AND status = 1 AND safety_stock > 0 AND quantity < safety_stock").
		Order("quantity ASC").Limit(10).Find(&warnings)

	return total, categoryStats, warnings, nil
}
