package models

import "gorm.io/gorm"

// Annotation lưu chú thích của tác giả cho từng dòng trong bài viết hoặc trang sách.
type Annotation struct {
	gorm.Model
	Content    string `gorm:"type:text;not null"`
	PostID     uint   `gorm:"index"`
	BookPageID uint   `gorm:"index"`
	LineNumber int    `gorm:"index"`
	Post       Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BookPage   BookPage `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
