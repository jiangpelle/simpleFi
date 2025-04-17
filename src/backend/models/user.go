package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string    `gorm:"uniqueIndex;not null" json:"username"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	Password      string    `gorm:"not null" json:"-"`
	WalletAddress string    `gorm:"uniqueIndex" json:"wallet_address"`
	LastLogin     time.Time `json:"last_login"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
}

type UserProfile struct {
	gorm.Model
	UserID      uint   `gorm:"uniqueIndex" json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	KYCVerified bool   `gorm:"default:false" json:"kyc_verified"`
}
