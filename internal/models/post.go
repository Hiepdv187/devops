package models

import "gorm.io/gorm"

// Post đại diện cho bài viết do người dùng tạo.
type Post struct {
	gorm.Model
	Title    string    `gorm:"size:200;not null"`
	Summary  string    `gorm:"size:255"`
	Content  string    `gorm:"type:text"`
	AuthorID uint      `gorm:"index"`
	Author   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Comments []Comment `gorm:"constraint:OnDelete:CASCADE;"`
}
