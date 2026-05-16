package model

import "time"

// ChangeLog 变更日志表
type ChangeLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InventoryID uint      `gorm:"not null;index" json:"inventory_id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	UserName    string    `gorm:"type:varchar(50);not null" json:"user_name"`
	FieldName   string    `gorm:"type:varchar(50);not null" json:"field_name"`
	FieldLabel  string    `gorm:"type:varchar(50);not null" json:"field_label"`
	OldValue    string    `gorm:"type:varchar(500);default:''" json:"old_value"`
	NewValue    string    `gorm:"type:varchar(500);default:''" json:"new_value"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (ChangeLog) TableName() string {
	return "change_logs"
}

// FieldLabelMap 字段名到中文标签的映射
var FieldLabelMap = map[string]string{
	"material_code": "物料编码",
	"material_name": "物料名称",
	"category":      "分类",
	"specification": "规格型号",
	"unit":          "单位",
	"quantity":      "库存数量",
	"safety_stock":  "安全库存",
	"location":      "仓库位置",
	"status":        "状态",
	"remark":        "备注",
}
