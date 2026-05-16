package service

import (
	"inventory-system/internal/model"
	"inventory-system/internal/repository"
	"inventory-system/pkg/errcode"
	"inventory-system/pkg/hash"
	"inventory-system/pkg/jwt"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepo
	logRepo  *repository.LogRepo
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepo(),
		logRepo:  repository.NewLogRepo(),
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password, ip string) (string, *model.User, *errcode.ErrCode) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, errcode.ErrLoginFailed
		}
		return "", nil, errcode.ErrDatabase
	}

	if user.Status == 0 {
		return "", nil, errcode.ErrAccountDisabled
	}

	if !hash.CheckPassword(password, user.PasswordHash) {
		return "", nil, errcode.ErrLoginFailed
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, errcode.ErrServer
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      user.ID,
		UserName:    user.Username,
		Action:      model.ActionLogin,
		TargetType:  model.TargetUser,
		TargetID:    &user.ID,
		Description: "用户登录",
		IPAddress:   ip,
	})

	return token, user, nil
}

// GetCurrentUser 获取当前登录用户信息
func (s *AuthService) GetCurrentUser(userID uint) (*model.User, *errcode.ErrCode) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return user, nil
}
