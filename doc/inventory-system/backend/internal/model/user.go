package model

// User 用户表
type User struct {
	BaseModel
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"type:varchar(200);not null" json:"-"`
	DisplayName  string `gorm:"type:varchar(100);default:''" json:"display_name"`
	Role         string `gorm:"type:varchar(20);not null;default:'viewer'" json:"role"`
	Status       int8   `gorm:"type:tinyint;not null;default:1" json:"status"`
}

func (User) TableName() string {
	return "users"
}

const (
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsOperator() bool {
	return u.Role == RoleOperator
}

func (u *User) HasWritePermission() bool {
	return u.Role == RoleAdmin || u.Role == RoleOperator
}
