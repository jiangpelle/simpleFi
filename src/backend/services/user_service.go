package services

import (
	"errors"
	"time"

	"defi-backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Register 用户注册
func (s *UserService) Register(username, email, password string) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		LastLogin: time.Now(),
		IsActive:  true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 更新最后登录时间
	user.LastLogin = time.Now()
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(userID uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint, profile *models.UserProfile) error {
	profile.UserID = userID
	return s.db.Save(profile).Error
}

// UpdateWalletAddress 更新钱包地址
func (s *UserService) UpdateWalletAddress(userID uint, walletAddress string) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("wallet_address", walletAddress).Error
}
