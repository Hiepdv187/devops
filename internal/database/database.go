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
		var dsn string

		// Priority 1: Check for DATABASE_URL (Supabase standard)
		databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
		if databaseURL != "" {
			// Add prefer_simple_protocol if not already present (required for Supabase pooler)
			if !strings.Contains(databaseURL, "prefer_simple_protocol") {
				if strings.Contains(databaseURL, "?") {
					dsn = databaseURL + "&prefer_simple_protocol=true"
				} else {
					dsn = databaseURL + "?prefer_simple_protocol=true"
				}
			} else {
				dsn = databaseURL
			}
			log.Println("✓ Using DATABASE_URL for connection")
		} else {
			// Priority 2: Check for DATABASE_DSN
			databaseDSN := strings.TrimSpace(os.Getenv("DATABASE_DSN"))
			if databaseDSN != "" {
				dsn = databaseDSN
				log.Println("✓ Using DATABASE_DSN for connection")
			} else {
				// Priority 3: Build from individual environment variables
				host := getEnv("DB_HOST", "localhost")
				port := getEnv("DB_PORT", "5432")
				user := getEnv("DB_USER", "postgres")
				password := getEnv("DB_PASSWORD", "")
				database := getEnv("DB_NAME", "postgres")
				sslmode := getEnv("DB_SSLMODE", "disable")

				dsn = "host=" + host + " user=" + user + " password=" + password + " dbname=" + database + " port=" + port + " sslmode=" + sslmode + " prefer_simple_protocol=true"
				log.Println("✓ Using individual DB_* variables for connection")
			}
		}

		// Custom logger to only show errors, not slow SQL warnings
		customLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             10 * time.Second, // Very high threshold to avoid SLOW SQL logs
				LogLevel:                  logger.Error,     // Only log errors
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)

		dbConfig := &gorm.Config{
			Logger:                                   customLogger,
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

		// Configure connection pool for Supabase pooler
		// Use smaller pool size to avoid prepared statement conflicts
		sqlDB.SetMaxOpenConns(5)                   // Reduced for pooler compatibility
		sqlDB.SetMaxIdleConns(2)                   // Keep minimal idle connections
		sqlDB.SetConnMaxLifetime(1 * time.Minute)  // Shorter lifetime for pooler
		sqlDB.SetConnMaxIdleTime(30 * time.Second) // Close idle connections quickly

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
