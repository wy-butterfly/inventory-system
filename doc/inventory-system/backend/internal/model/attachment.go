package model

// Attachment 附件表
type Attachment struct {
	BaseModel
	InventoryID   uint   `gorm:"not null;index" json:"inventory_id"`
	FileName      string `gorm:"type:varchar(200);not null" json:"file_name"`
	DisplayName   string `gorm:"type:varchar(200);not null" json:"display_name"`
	FilePath      string `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize      int64  `gorm:"not null" json:"file_size"`
	FileType      string `gorm:"type:varchar(20);not null" json:"file_type"`
	MimeType      string `gorm:"type:varchar(100);not null" json:"mime_type"`
	Extension     string `gorm:"type:varchar(10);not null" json:"extension"`
	ThumbnailPath string `gorm:"type:varchar(500);default:''" json:"thumbnail_path"`
	UploadedBy    uint   `gorm:"not null" json:"uploaded_by"`
	IsDeleted     int8   `gorm:"type:tinyint;not null;default:0" json:"is_deleted"`

	UploaderName string `gorm:"-" json:"uploader_name,omitempty"`
}

func (Attachment) TableName() string {
	return "attachments"
}

const (
	FileTypeImage    = "image"
	FileTypeDocument = "document"
)
