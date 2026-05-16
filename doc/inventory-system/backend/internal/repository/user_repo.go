package repository

import (
	"inventory-system/internal/model"
	"inventory-system/pkg/database"
)

// UserRepo 用户数据访问
type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// FindByUsername 根据用户名查找用户
func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepo) Update(user *model.User) error {
	return database.DB.Save(user).Error
}

// List 分页查询用户列表
func (r *UserRepo) List(keyword, role string, status *int8, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := database.DB.Model(&model.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR display_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// GetNamesByIDs 批量获取用户名称
func (r *UserRepo) GetNamesByIDs(ids []uint) (map[uint]string, error) {
	var users []model.User
	err := database.DB.Select("id, display_name, username").Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}

	nameMap := make(map[uint]string)
	for _, u := range users {
		name := u.DisplayName
		if name == "" {
			name = u.Username
		}
		nameMap[u.ID] = name
	}
	return nameMap, nil
}
