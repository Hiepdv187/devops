package database

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"fiber-learning-community/internal/models"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init khởi tạo kết nối GORM.
func Init() *gorm.DB {
	once.Do(func() {
		defaultDSN := "root:rootpass@tcp(127.0.0.1:3308)/fiber_learning?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := strings.TrimSpace(getEnv("DATABASE_DSN", defaultDSN))
		if dsn == "" {
			dsn = defaultDSN
		}

		dbConfig := &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Warn),
			NamingStrategy: schema.NamingStrategy{SingularTable: false},
		}

		var err error
		db, err = gorm.Open(mysql.Open(dsn), dbConfig)
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to retrieve sql DB instance: %v", err)
		}
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(15 * time.Minute)

		if err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Annotation{}, &models.Image{}); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}

		seedDemoUser()
	})

	return db
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

// Get trả về instance GORM đã được khởi tạo.
func Get() *gorm.DB {
	if db == nil {
		return Init()
	}
	return db
}

func seedDemoUser() {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		log.Printf("failed to count users: %v", err)
		return
	}
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("devops123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to seed user: %v", err)
		return
	}

	demo := models.User{
		Name:         "DevOps Maintainer",
		Email:        "admin@hocdevops.community",
		PasswordHash: string(hash),
	}

	if err := db.Create(&demo).Error; err != nil {
		log.Printf("failed to create demo user: %v", err)
	}
}
