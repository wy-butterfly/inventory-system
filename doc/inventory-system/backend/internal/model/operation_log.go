package model

import "time"

// OperationLog 操作日志表
type OperationLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	UserName    string    `gorm:"type:varchar(50);not null" json:"user_name"`
	Action      string    `gorm:"type:varchar(50);not null;index" json:"action"`
	TargetType  string    `gorm:"type:varchar(50);not null" json:"target_type"`
	TargetID    *uint     `gorm:"default:null" json:"target_id"`
	Description string    `gorm:"type:varchar(500);default:''" json:"description"`
	IPAddress   string    `gorm:"type:varchar(50);default:''" json:"ip_address"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

const (
	ActionCreate   = "create"
	ActionUpdate   = "update"
	ActionDelete   = "delete"
	ActionRestore  = "restore"
	ActionUpload   = "upload"
	ActionDownload = "download"
	ActionExport   = "export"
	ActionLogin    = "login"
	ActionLogout   = "logout"
)

const (
	TargetInventory  = "inventory"
	TargetAttachment = "attachment"
	TargetUser       = "user"
)
