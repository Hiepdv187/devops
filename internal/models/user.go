package models

import (
	"gorm.io/gorm"
)

// User lưu thông tin người dùng có thể đăng nhập hệ thống.
type User struct {
	gorm.Model
	Name         string `gorm:"size:120"`
	Email        string `gorm:"size:120;uniqueIndex"`
	PasswordHash string `gorm:"size:255"`
	Posts        []Post    `gorm:"foreignKey:AuthorID"`
	Comments     []Comment `gorm:"foreignKey:AuthorID"`
}
