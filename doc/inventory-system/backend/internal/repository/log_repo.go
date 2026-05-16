package repository

import (
	"inventory-system/internal/model"
	"inventory-system/pkg/database"
	"time"
)

type LogRepo struct{}

func NewLogRepo() *LogRepo {
	return &LogRepo{}
}

// CreateOperationLog 记录操作日志
func (r *LogRepo) CreateOperationLog(log *model.OperationLog) error {
	return database.DB.Create(log).Error
}

// CreateChangeLog 记录变更日志
func (r *LogRepo) CreateChangeLog(log *model.ChangeLog) error {
	return database.DB.Create(log).Error
}

// CreateChangeLogs 批量记录变更日志
func (r *LogRepo) CreateChangeLogs(logs []model.ChangeLog) error {
	if len(logs) == 0 {
		return nil
	}
	return database.DB.Create(&logs).Error
}

// ListOperationLogs 查询操作日志
func (r *LogRepo) ListOperationLogs(action, targetType string, userID *uint,
	startTime, endTime *time.Time, offset, limit int) ([]model.OperationLog, int64, error) {

	var logs []model.OperationLog
	var total int64

	query := database.DB.Model(&model.OperationLog{})

	if action != "" {
		query = query.Where("action = ?", action)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, total, err
}

// ListChangeLogs 查询变更日志
func (r *LogRepo) ListChangeLogs(inventoryID *uint, userID *uint,
	startTime, endTime *time.Time, offset, limit int) ([]model.ChangeLog, int64, error) {

	var logs []model.ChangeLog
	var total int64

	query := database.DB.Model(&model.ChangeLog{})

	if inventoryID != nil {
		query = query.Where("inventory_id = ?", *inventoryID)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, total, err
}

// GetRecentOperations 获取最近操作日志（仪表盘用）
func (r *LogRepo) GetRecentOperations(limit int) ([]model.OperationLog, error) {
	var logs []model.OperationLog
	err := database.DB.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
