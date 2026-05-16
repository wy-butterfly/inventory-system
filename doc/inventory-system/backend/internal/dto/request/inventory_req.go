package request

// CreateInventoryReq 新增库存请求参数
type CreateInventoryReq struct {
	MaterialCode  string  `json:"material_code" binding:"required,min=1,max=50"`
	MaterialName  string  `json:"material_name" binding:"required,min=1,max=200"`
	Category      string  `json:"category" binding:"required"`
	Specification string  `json:"specification" binding:"omitempty,max=200"`
	Unit          string  `json:"unit" binding:"required,min=1,max=20"`
	Quantity      float64 `json:"quantity" binding:"required,min=0"`
	SafetyStock   float64 `json:"safety_stock" binding:"omitempty,min=0"`
	Location      string  `json:"location" binding:"omitempty,max=200"`
	Status        *int8   `json:"status" binding:"omitempty,oneof=0 1"`
	Remark        string  `json:"remark" binding:"omitempty,max=1000"`
}

// UpdateInventoryReq 编辑库存请求参数
type UpdateInventoryReq struct {
	MaterialName  *string  `json:"material_name" binding:"omitempty,min=1,max=200"`
	Category      *string  `json:"category" binding:"omitempty"`
	Specification *string  `json:"specification" binding:"omitempty,max=200"`
	Unit          *string  `json:"unit" binding:"omitempty,min=1,max=20"`
	Quantity      *float64 `json:"quantity" binding:"omitempty,min=0"`
	SafetyStock   *float64 `json:"safety_stock" binding:"omitempty,min=0"`
	Location      *string  `json:"location" binding:"omitempty,max=200"`
	Status        *int8    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark        *string  `json:"remark" binding:"omitempty,max=1000"`
	Version       int      `json:"version" binding:"required,min=1"`
}

// InventoryQuery 库存列表查询参数
type InventoryQuery struct {
	PageQuery
	Keyword     string   `form:"keyword" binding:"omitempty,max=100"`
	Category    string   `form:"category" binding:"omitempty"`
	Status      *int8    `form:"status" binding:"omitempty,oneof=0 1"`
	MinQuantity *float64 `form:"min_quantity" binding:"omitempty,min=0"`
	MaxQuantity *float64 `form:"max_quantity" binding:"omitempty,min=0"`
	Location    string   `form:"location" binding:"omitempty,max=200"`
	ShowDeleted bool     `form:"show_deleted"`
}

// BatchDeleteReq 批量删除请求参数
type BatchDeleteReq struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
