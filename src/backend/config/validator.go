package config

import (
	"fmt"
	"strings"
)

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	// 验证 Nacos 配置
	if c.Nacos.ServerAddr == "" {
		return fmt.Errorf("nacos server address is required")
	}
	if c.Nacos.Port == 0 {
		return fmt.Errorf("nacos port is required")
	}
	if c.Nacos.Namespace == "" {
		return fmt.Errorf("nacos namespace is required")
	}
	if c.Nacos.Group == "" {
		return fmt.Errorf("nacos group is required")
	}
	if c.Nacos.DataId == "" {
		return fmt.Errorf("nacos data id is required")
	}

	// 验证数据库配置
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port == 0 {
		return fmt.Errorf("database port is required")
	}
	if c.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// 验证 Redis 配置
	if c.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}
	if c.Redis.Port == 0 {
		return fmt.Errorf("redis port is required")
	}

	return nil
}

// ValidateUser 验证用户输入
func ValidateUser(username, email, password string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

// ValidateTrade 验证交易参数
func ValidateTrade(amount, price float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	return nil
}

// ValidatePosition 验证仓位参数
func ValidatePosition(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}
