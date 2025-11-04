package models

import (
	"time"
	"gorm.io/gorm"
)

type Highlight struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	BookPageID   uint   `gorm:"not null;index" json:"book_page_id"`
	UserID       uint   `gorm:"not null;index" json:"user_id"`
	
	// Highlight data
	Color        string `gorm:"type:varchar(20);not null" json:"color"`          // e.g., "#ffeb3b", "#4caf50"
	HighlightedText string `gorm:"type:text;not null" json:"highlighted_text"`  // The actual text that was highlighted
	Note         string `gorm:"type:text" json:"note"`                          // Optional note for the highlight
	
	// Position data for re-rendering
	StartOffset  int    `gorm:"not null" json:"start_offset"`  // Character offset from start of page content
	EndOffset    int    `gorm:"not null" json:"end_offset"`    // Character offset from start of page content
	
	// Relations
	BookPage     BookPage `gorm:"foreignKey:BookPageID" json:"-"`
	User         User     `gorm:"foreignKey:UserID" json:"-"`
}
