package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"fiber-learning-community/internal/database"
	"fiber-learning-community/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Helper function to get current user from session
func getUserForBooks(c *fiber.Ctx) *models.User {
	userID, err := currentUserID(c)
	if err != nil || userID == 0 {
		return nil
	}

	db := database.Get()
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil
	}
	return &user
}

// BooksPage hiển thị danh sách sách
func BooksPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		user := getUserForBooks(c)

		var books []models.Book
		query := db.Order("created_at DESC")
		
		// Chỉ hiển thị sách published hoặc sách của mình
		if user != nil {
			query = query.Where("published = ? OR author_id = ?", true, user.ID)
		} else {
			query = query.Where("published = ?", true)
		}
		
		if err := query.Find(&books).Error; err != nil {
			return c.Status(500).SendString("Lỗi tải danh sách sách")
		}

		// Load author names
		for i := range books {
			var author models.User
			if err := db.First(&author, books[i].AuthorID).Error; err == nil {
				books[i].AuthorName = author.Name
			}
		}

		return render(c, "pages/books", fiber.Map{
			"Title": "Sách",
			"Books": books,
		}, "main")
	}
}

// BookDetailPage hiển thị chi tiết sách và cho phép chỉnh sửa
func BookDetailPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		user := getUserForBooks(c)
		
		bookID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("ID không hợp lệ")
		}

		var book models.Book
		if err := db.Preload("Pages", func(db *gorm.DB) *gorm.DB {
			return db.Order("page_number ASC")
		}).First(&book, bookID).Error; err != nil {
			return c.Status(404).SendString("Không tìm thấy sách")
		}

		// Load author
		var author models.User
		if err := db.First(&author, book.AuthorID).Error; err == nil {
			book.AuthorName = author.Name
		}

		isAuthor := user != nil && user.ID == book.AuthorID

		return render(c, "pages/book_detail", fiber.Map{
			"Title":    book.Title,
			"Book":     book,
			"User":     user,
			"IsAuthor": isAuthor,
		}, "main")
	}
}

// BookReadPage hiển thị sách ở chế độ đọc với page flip
func BookReadPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		user := getUserForBooks(c)

		bookID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("ID không hợp lệ")
		}

		var book models.Book
		if err := db.Preload("Pages", func(db *gorm.DB) *gorm.DB {
			return db.Order("page_number ASC")
		}).First(&book, bookID).Error; err != nil {
			return c.Status(404).SendString("Không tìm thấy sách")
		}

		// Load annotations for each page
		for i := range book.Pages {
			var annotations []models.Annotation
			db.Where("book_page_id = ?", book.Pages[i].ID).Find(&annotations)

			// Convert to map[lineNumber]content
			annotationsMap := make(map[int]string)
			for _, ann := range annotations {
				annotationsMap[ann.LineNumber] = ann.Content
			}
			book.Pages[i].Annotations = annotations
		}

		// Load author
		var author models.User
		if err := db.First(&author, book.AuthorID).Error; err == nil {
			book.AuthorName = author.Name
		}

		isAuthor := user != nil && user.ID == book.AuthorID

		// Check if request accepts JSON (for AJAX/fetch calls)
		if c.Get("Accept") == "application/json" {
			return c.JSON(fiber.Map{
				"id":          book.ID,
				"title":       book.Title,
				"description": book.Description,
				"cover_url":   book.CoverURL,
				"cover_color": book.CoverColor,
				"author_id":   book.AuthorID,
				"author_name": book.AuthorName,
				"published":   book.Published,
				"pages":       book.Pages,
			})
		}

		return render(c, "pages/book_read", fiber.Map{
			"Title":    book.Title,
			"Book":     book,
			"IsAuthor": isAuthor,
		}, "main")
	}
}

// CreateBook tạo sách mới
func CreateBook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Chưa đăng nhập"})
		}

		db := database.Get()

		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CoverURL    string `json:"cover_url"`
			CoverColor  string `json:"cover_color"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
		}

		if strings.TrimSpace(req.Title) == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Tiêu đề không được để trống"})
		}

		coverColor := strings.TrimSpace(req.CoverColor)
		if coverColor == "" {
			coverColor = "#1e293b"
		}

		book := models.Book{
			Title:       strings.TrimSpace(req.Title),
			Description: strings.TrimSpace(req.Description),
			CoverURL:    strings.TrimSpace(req.CoverURL),
			CoverColor:  coverColor,
			AuthorID:    user.ID,
			Published:   false,
		}

		if err := db.Create(&book).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi tạo sách"})
		}

		// Tạo trang đầu tiên tự động
		firstPage := models.BookPage{
			BookID:     book.ID,
			PageNumber: 1,
			Title:      "Trang 1",
			Content:    "<p>Bắt đầu viết nội dung của bạn ở đây...</p>",
		}
		
		if err := db.Create(&firstPage).Error; err != nil {
			// Log error but don't fail book creation
			println("Warning: Failed to create first page:", err.Error())
		}

		return c.JSON(fiber.Map{"success": true, "book_id": book.ID})
	}
}

// UpdateBook cập nhật thông tin sách
func UpdateBook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Redirect("/auth/login")
		}

		db := database.Get()
		bookID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID không hợp lệ"})
		}

		var book models.Book
		if err := db.First(&book, bookID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy sách"})
		}

		if book.AuthorID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Không có quyền"})
		}

		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CoverURL    string `json:"cover_url"`
			CoverColor  string `json:"cover_color"`
			Published   bool   `json:"published"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
		}

		book.Title = strings.TrimSpace(req.Title)
		book.Description = strings.TrimSpace(req.Description)
		book.CoverURL = strings.TrimSpace(req.CoverURL)
		if req.CoverColor != "" {
			book.CoverColor = strings.TrimSpace(req.CoverColor)
		}
		book.Published = req.Published

		if err := db.Save(&book).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi cập nhật sách"})
		}

		return c.JSON(fiber.Map{"success": true})
	}
}

// CreateBookPage tạo trang mới cho sách
func CreateBookPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Chưa đăng nhập"})
		}

		db := database.Get()
		bookID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID không hợp lệ"})
		}

		var book models.Book
		if err := db.First(&book, bookID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy sách"})
		}

		if book.AuthorID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Không có quyền"})
		}

		var req struct {
			Title           string `json:"title"`
			Content         string `json:"content"`
			PageNumber      int    `json:"page_number"`
			LineAnnotations string `json:"line_annotations"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
		}

		// Determine page number
		pageNumber := req.PageNumber
		if pageNumber <= 0 {
			// Get next page number
			var maxPage int
			db.Model(&models.BookPage{}).Where("book_id = ?", bookID).Select("COALESCE(MAX(page_number), 0)").Scan(&maxPage)
			pageNumber = maxPage + 1
		}

		page := models.BookPage{
			BookID:     uint(bookID),
			PageNumber: pageNumber,
			Title:      strings.TrimSpace(req.Title),
			Content:    req.Content,
		}

		if err := db.Create(&page).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi tạo trang"})
		}

		// Save annotations
		if req.LineAnnotations != "" {
			var annotationsMap map[string]string
			if err := json.Unmarshal([]byte(req.LineAnnotations), &annotationsMap); err == nil {
				for lineNumStr, annotationText := range annotationsMap {
					lineNum, err := strconv.Atoi(lineNumStr)
					if err != nil {
						continue
					}
					annotation := models.Annotation{
						BookPageID: page.ID,
						LineNumber: lineNum,
						Content:    strings.TrimSpace(annotationText),
					}
					db.Create(&annotation)
				}
			}
		}

		return c.JSON(fiber.Map{"success": true, "page_id": page.ID})
	}
}

// UpdateBookPage cập nhật trang sách
func UpdateBookPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Chưa đăng nhập"})
		}

		db := database.Get()
		bookID, _ := strconv.Atoi(c.Params("bookId"))
		pageID, _ := strconv.Atoi(c.Params("pageId"))

		var book models.Book
		if err := db.First(&book, bookID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy sách"})
		}

		if book.AuthorID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Không có quyền"})
		}

		var page models.BookPage
		if err := db.First(&page, pageID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy trang"})
		}

		var req struct {
			Title           string `json:"title"`
			Content         string `json:"content"`
			LineAnnotations string `json:"line_annotations"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
		}

		page.Title = strings.TrimSpace(req.Title)
		page.Content = req.Content

		if err := db.Save(&page).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi cập nhật trang"})
		}

		// Update annotations
		db.Where("book_page_id = ?", page.ID).Delete(&models.Annotation{})
		
		if req.LineAnnotations != "" {
			var annotationsMap map[string]string
			if err := json.Unmarshal([]byte(req.LineAnnotations), &annotationsMap); err == nil {
				for lineNumStr, annotationText := range annotationsMap {
					lineNum, err := strconv.Atoi(lineNumStr)
					if err != nil {
						continue
					}
					annotation := models.Annotation{
						BookPageID: page.ID,
						LineNumber: lineNum,
						Content:    strings.TrimSpace(annotationText),
					}
					db.Create(&annotation)
				}
			}
		}

		return c.JSON(fiber.Map{"success": true})
	}
}

// DeleteBookPage xóa trang sách
func DeleteBookPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Chưa đăng nhập"})
		}

		db := database.Get()
		bookID, _ := strconv.Atoi(c.Params("bookId"))
		pageID, _ := strconv.Atoi(c.Params("pageId"))

		var book models.Book
		if err := db.First(&book, bookID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy sách"})
		}

		if book.AuthorID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Không có quyền"})
		}

		if err := db.Delete(&models.BookPage{}, pageID).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi xóa trang"})
		}

		return c.JSON(fiber.Map{"success": true})
	}
}
