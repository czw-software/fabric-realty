package service

import (
	"application/model"
	"application/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

// Register 用户注册
func (s *UserService) Register(user *model.User) error {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := model.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("username already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = string(hashedPassword)

	// 创建用户
	if err := model.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, string, error) {
	var user model.User
	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", fmt.Errorf("user not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid password")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return &user, token, nil
}
