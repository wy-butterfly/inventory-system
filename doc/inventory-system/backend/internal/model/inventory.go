package model

// Inventory 库存表
type Inventory struct {
	BaseModel
	MaterialCode  string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"material_code"`
	MaterialName  string  `gorm:"type:varchar(200);not null" json:"material_name"`
	Category      string  `gorm:"type:varchar(50);not null" json:"category"`
	Specification string  `gorm:"type:varchar(200);default:''" json:"specification"`
	Unit          string  `gorm:"type:varchar(20);not null" json:"unit"`
	Quantity      float64 `gorm:"type:decimal(12,2);not null;default:0" json:"quantity"`
	SafetyStock   float64 `gorm:"type:decimal(12,2);default:0" json:"safety_stock"`
	Location      string  `gorm:"type:varchar(200);default:''" json:"location"`
	Status        int8    `gorm:"type:tinyint;not null;default:1" json:"status"`
	Remark        string  `gorm:"type:varchar(1000);default:''" json:"remark"`
	Version       int     `gorm:"type:int;not null;default:1" json:"version"`
	IsDeleted     int8    `gorm:"type:tinyint;not null;default:0" json:"is_deleted"`
	CreatedBy     uint    `gorm:"not null" json:"created_by"`
	UpdatedBy     *uint   `gorm:"default:null" json:"updated_by"`

	Attachments []Attachment `gorm:"foreignKey:InventoryID" json:"attachments,omitempty"`
	CreatorName string       `gorm:"-" json:"creator_name,omitempty"`
	Warning     bool         `gorm:"-" json:"warning"`
}

func (Inventory) TableName() string {
	return "inventories"
}

const (
	CategoryRawMaterial     = "raw_material"
	CategoryFinishedProduct = "finished_product"
	CategorySparePart       = "spare_part"
	CategoryOther           = "other"
)

var ValidCategories = []string{
	CategoryRawMaterial,
	CategoryFinishedProduct,
	CategorySparePart,
	CategoryOther,
}

func IsValidCategory(category string) bool {
	for _, c := range ValidCategories {
		if c == category {
			return true
		}
	}
	return false
}
