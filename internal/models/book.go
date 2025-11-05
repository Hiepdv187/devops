package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Title        string         `gorm:"not null" json:"title"`
	Description  string         `json:"description"`
	CoverURL     string         `json:"cover_url"`
	CoverColor   string         `gorm:"default:#1e293b" json:"cover_color"`
	AuthorID     uint           `gorm:"not null" json:"author_id"`
	AuthorName   string         `gorm:"-" json:"author_name"`
	Published    bool           `gorm:"default:false" json:"published"`
	BookTag      string         `gorm:"index" json:"book_tag"`      // Tag for book (e.g., "linux", "golang")
	BookCategory string         `gorm:"index" json:"book_category"` // Category for filtering
	Pages        []BookPage     `gorm:"foreignKey:BookID" json:"pages,omitempty"`
}

type BookPage struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	BookID     uint           `gorm:"not null;index" json:"book_id"`
	PageNumber int            `gorm:"not null" json:"page_number"`
	Title      string         `json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
}
