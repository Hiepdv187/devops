package models

import "gorm.io/gorm"

// Comment lưu bình luận của người dùng trên bài viết.
type Comment struct {
	gorm.Model
	Content    string `gorm:"type:text;not null"`
	PostID     uint   `gorm:"index"`
	AuthorID   uint   `gorm:"index"`
	LineNumber *int   `gorm:"index"`
	Post       Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Author     User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
