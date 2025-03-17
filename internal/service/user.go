package service

import (
	"errors"

	"gorm.io/gorm"

	"cinexus/internal/database"
	"cinexus/internal/model"
	"cinexus/pkg/jwt"
)

// UserService 用户服务
type UserService struct{}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user model.User

	// 查询用户
	err := database.DB.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 检查密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) error {
	// 检查用户名是否已存在
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		database.DB.Model(&model.User{}).Where("email = ?", req.Email).Count(&count)
		if count > 0 {
			return errors.New("邮箱已存在")
		}
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     "user",
		Status:   1,
	}

	return database.DB.Create(&user).Error
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, req *UpdateUserRequest) error {
	// 检查邮箱是否已存在
	if req.Email != "" {
		var count int64
		database.DB.Model(&model.User{}).Where("email = ? AND id != ?", req.Email, id).Count(&count)
		if count > 0 {
			return errors.New("邮箱已存在")
		}
	}

	// 更新用户
	return database.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"avatar":   req.Avatar,
	}).Error
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(id uint, req *UpdatePasswordRequest) error {
	var user model.User

	// 查询用户
	err := database.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 检查旧密码
	if !user.CheckPassword(req.OldPassword) {
		return errors.New("旧密码错误")
	}

	// 更新密码
	user.Password = req.NewPassword
	return database.DB.Save(&user).Error
}
