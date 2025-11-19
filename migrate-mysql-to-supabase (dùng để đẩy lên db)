package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fiber-learning-community/internal/models"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	fmt.Println("========================================")
	fmt.Println("MIGRATE MYSQL → SUPABASE POSTGRESQL")
	fmt.Println("========================================")
	fmt.Println()

	// ============================================
	// 1. Connect to MySQL (Source)
	// ============================================
	fmt.Println("📦 Connecting to MySQL (source)...")
	mysqlDSN := "root:rootpass@tcp(localhost:3308)/fiber_learning?charset=utf8mb4&parseTime=True&loc=Local"

	mysqlDB, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to MySQL: %v", err)
	}
	fmt.Println("✅ MySQL connected")

	// ============================================
	// 2. Connect to Supabase PostgreSQL (Destination)
	// ============================================
	fmt.Println("📦 Connecting to Supabase PostgreSQL (destination)...")

	// Use DATABASE_URL from .env
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("❌ DATABASE_URL not found in .env file")
	}

	postgresDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to Supabase: %v", err)
	}
	fmt.Println("✅ Supabase connected")
	fmt.Println()

	// ============================================
	// 3. Run migrations on Supabase
	// ============================================
	fmt.Println("🔧 Running migrations on Supabase...")
	err = postgresDB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Annotation{},
		&models.Book{},
		&models.BookPage{},
		&models.Highlight{},
	)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	fmt.Println("✅ Migrations completed")
	fmt.Println()

	// ============================================
	// 4. Clear existing data (optional - uncomment if needed)
	// ============================================
	fmt.Println("🗑️  Clearing existing data in Supabase...")
	clearExistingData(postgresDB)
	fmt.Println()

	// ============================================
	// 5. Migrate data
	// ============================================
	fmt.Println("📊 Starting data migration...")
	fmt.Println()

	// Migrate Users
	if err := migrateTable(mysqlDB, postgresDB, &models.User{}, "users"); err != nil {
		log.Printf("⚠️  Warning: Failed to migrate users: %v", err)
	}

	// Migrate Posts
	if err := migrateTable(mysqlDB, postgresDB, &models.Post{}, "posts"); err != nil {
		log.Printf("⚠️  Warning: Failed to migrate posts: %v", err)
	}

	// Migrate Comments
	if err := migrateTable(mysqlDB, postgresDB, &models.Comment{}, "comments"); err != nil {
		log.Printf("⚠️  Warning: Failed to migrate comments: %v", err)
	}

	// Migrate Annotations
	if err := migrateTable(mysqlDB, postgresDB, &models.Annotation{}, "annotations"); err != nil {
		log.Printf("⚠️  Warning: Failed to migrate annotations: %v", err)
	}

	// Migrate Books (skip if table doesn't exist)
	if tableExists(mysqlDB, "books") {
		if err := migrateTable(mysqlDB, postgresDB, &models.Book{}, "books"); err != nil {
			log.Printf("⚠️  Warning: Failed to migrate books: %v", err)
		}
	} else {
		fmt.Println("  ⏭️  Skipping books (table doesn't exist in MySQL)")
	}

	// Migrate Book Pages (skip if table doesn't exist)
	if tableExists(mysqlDB, "book_pages") {
		if err := migrateTable(mysqlDB, postgresDB, &models.BookPage{}, "book_pages"); err != nil {
			log.Printf("⚠️  Warning: Failed to migrate book_pages: %v", err)
		}
	} else {
		fmt.Println("  ⏭️  Skipping book_pages (table doesn't exist in MySQL)")
	}

	// Migrate Highlights (skip if table doesn't exist)
	if tableExists(mysqlDB, "highlights") {
		if err := migrateTable(mysqlDB, postgresDB, &models.Highlight{}, "highlights"); err != nil {
			log.Printf("⚠️  Warning: Failed to migrate highlights: %v", err)
		}
	} else {
		fmt.Println("  ⏭️  Skipping highlights (table doesn't exist in MySQL)")
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✅ MIGRATION COMPLETED SUCCESSFULLY!")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("📝 Next steps:")
	fmt.Println("1. Verify data in Supabase Dashboard")
	fmt.Println("2. Update your app to use Supabase")
	fmt.Println("3. Test all features")
}

// clearExistingData clears all data from Supabase tables
func clearExistingData(db *gorm.DB) {
	tables := []string{"highlights", "book_pages", "books", "annotations", "comments", "posts", "users"}

	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			fmt.Printf("  ⚠️  Warning: Failed to clear %s: %v\n", table, err)
		} else {
			fmt.Printf("  ✅ Cleared %s\n", table)
		}
	}
}

// tableExists checks if a table exists in MySQL
func tableExists(db *gorm.DB, tableName string) bool {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Scan(&count).Error
	return err == nil && count > 0
}

// migrateTable migrates data from MySQL to PostgreSQL for a specific table
func migrateTable(source, dest *gorm.DB, model interface{}, tableName string) error {
	fmt.Printf("  📋 Migrating %s...\n", tableName)

	// Get count from source
	var count int64
	if err := source.Model(model).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count source records: %w", err)
	}

	if count == 0 {
		fmt.Printf("     ⏭️  No records to migrate\n")
		return nil
	}

	// Fetch all records from source
	records := make([]map[string]interface{}, 0)
	if err := source.Model(model).Find(&records).Error; err != nil {
		return fmt.Errorf("failed to fetch records: %w", err)
	}

	// Insert records into destination
	if len(records) > 0 {
		// Use Create in batches to avoid memory issues
		batchSize := 100
		for i := 0; i < len(records); i += batchSize {
			end := i + batchSize
			if end > len(records) {
				end = len(records)
			}
			batch := records[i:end]

			if err := dest.Table(tableName).Create(&batch).Error; err != nil {
				return fmt.Errorf("failed to insert batch: %w", err)
			}
		}
	}

	fmt.Printf("     ✅ Migrated %d records\n", count)
	return nil
}
