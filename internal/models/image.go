package models

import "gorm.io/gorm"

// Image đại diện cho ảnh được upload vào hệ thống.
type Image struct {
	gorm.Model
	Filename    string `gorm:"size:255;not null"`           // Tên file gốc
	ContentType string `gorm:"size:100;not null"`           // MIME type (image/jpeg, image/png, etc.)
	Size        int64  `gorm:"not null"`                    // Kích thước file (bytes)
	Data        []byte `gorm:"type:longblob;not null"`      // Dữ liệu ảnh (binary)
	UploaderID  uint   `gorm:"index"`                       // ID người upload
	Uploader    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
