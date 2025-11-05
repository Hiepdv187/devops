package database

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
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
		// PostgreSQL connection parameters
		host := getEnv("DB_HOST", "ep-odd-morning-a7ox4rz0-pooler.ap-southeast-2.aws.neon.tech")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "neondb_owner")
		password := getEnv("DB_PASSWORD", "npg_ZCva2xmGOt3g")
		database := getEnv("DB_NAME", "wedevops")
		sslmode := getEnv("DB_SSLMODE", "require")

		// Tạo PostgreSQL DSN with prefer_simple_protocol to avoid prepared statement issues
		defaultDSN := "host=" + host + " user=" + user + " password=" + password + " dbname=" + database + " port=" + port + " sslmode=" + sslmode + " prefer_simple_protocol=true"
		dsn := strings.TrimSpace(getEnv("DATABASE_DSN", defaultDSN))
		if dsn == "" {
			dsn = defaultDSN
		}

		dbConfig := &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Warn),
			NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
			PrepareStmt:                              false, // Disable prepared statement cache to avoid schema mismatch
			DisableForeignKeyConstraintWhenMigrating: true,
		}

		var err error
		db, err = gorm.Open(postgres.Open(dsn), dbConfig)
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to retrieve sql DB instance: %v", err)
		}
		
		// Close all existing connections to clear prepared statements
		sqlDB.SetMaxOpenConns(0)
		sqlDB.SetMaxIdleConns(0)
		time.Sleep(100 * time.Millisecond) // Wait for connections to close
		
		// Set new connection pool settings
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute) // Shorter lifetime to avoid stale connections

		if err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Annotation{}, &models.Book{}, &models.BookPage{}, &models.Highlight{}); err != nil {
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
