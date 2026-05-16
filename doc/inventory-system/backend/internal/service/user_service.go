package service

import (
	"fmt"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/model"
	"inventory-system/internal/repository"
	"inventory-system/pkg/errcode"
	"inventory-system/pkg/hash"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepo
	logRepo  *repository.LogRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepo(),
		logRepo:  repository.NewLogRepo(),
	}
}

// Create 创建用户
func (s *UserService) Create(req *request.CreateUserReq, operatorID uint, operatorName string) (*model.User, *errcode.ErrCode) {
	existing, err := s.userRepo.FindByUsername(req.Username)
	if err == nil && existing != nil {
		return nil, errcode.ErrDuplicate.WithMessage("用户名已存在")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errcode.ErrDatabase
	}

	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, errcode.ErrServer
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
		Role:         req.Role,
		Status:       1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionCreate,
		TargetType:  model.TargetUser,
		TargetID:    &user.ID,
		Description: fmt.Sprintf("创建用户: %s (角色: %s)", user.Username, user.Role),
	})

	return user, nil
}

// Update 更新用户
func (s *UserService) Update(id uint, req *request.UpdateUserReq, operatorID uint, operatorName string) *errcode.ErrCode {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errcode.ErrNotFound
	}

	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.userRepo.Update(user); err != nil {
		return errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionUpdate,
		TargetType:  model.TargetUser,
		TargetID:    &id,
		Description: fmt.Sprintf("编辑用户: %s", user.Username),
	})

	return nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(id uint, newPassword string, operatorID uint, operatorName string) *errcode.ErrCode {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errcode.ErrNotFound
	}

	passwordHash, err := hash.HashPassword(newPassword)
	if err != nil {
		return errcode.ErrServer
	}

	user.PasswordHash = passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionUpdate,
		TargetType:  model.TargetUser,
		TargetID:    &id,
		Description: fmt.Sprintf("重置用户密码: %s", user.Username),
	})

	return nil
}

// List 查询用户列表
func (s *UserService) List(query *request.UserQuery) ([]model.User, int64, *errcode.ErrCode) {
	users, total, err := s.userRepo.List(
		query.Keyword, query.Role, query.Status,
		query.GetOffset(), query.GetPageSize(),
	)
	if err != nil {
		return nil, 0, errcode.ErrDatabase
	}
	return users, total, nil
}
