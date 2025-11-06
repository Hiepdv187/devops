package handlers

import (
	"log"
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

		// One-time migration: Update all existing books to published = true
		// This can be removed after first run
		db.Model(&models.Book{}).Where("published = ?", false).Update("published", true)

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

		isAuthenticated := user != nil

		return render(c, "pages/books", fiber.Map{
			"Title":           "Sách",
			"Books":           books,
			"IsAuthenticated": isAuthenticated,
			"CurrentUser":     user,
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

		// Load author
		var author models.User
		if err := db.First(&author, book.AuthorID).Error; err == nil {
			book.AuthorName = author.Name
		}

		isAuthor := user != nil && user.ID == book.AuthorID

		isAuthenticated := user != nil

		// Check if request accepts JSON (for AJAX/fetch calls)
		if c.Get("Accept") == "application/json" {
			return c.JSON(fiber.Map{
				"id":              book.ID,
				"title":           book.Title,
				"description":     book.Description,
				"cover_url":       book.CoverURL,
				"cover_color":     book.CoverColor,
				"author_id":       book.AuthorID,
				"author_name":     book.AuthorName,
				"published":       book.Published,
				"pages":           book.Pages,
				"is_author":       isAuthor,
				"is_authenticated": isAuthenticated,
			})
		}

		return render(c, "pages/book_read", fiber.Map{
			"Title":    book.Title,
			"Book":     book,
			"IsAuthor": isAuthor,
		}, "empty")
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
			Title        string `json:"title"`
			Description  string `json:"description"`
			CoverURL     string `json:"cover_url"`
			CoverColor   string `json:"cover_color"`
			BookTag      string `json:"book_tag"`
			BookCategory string `json:"book_category"`
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
			Title:        strings.TrimSpace(req.Title),
			Description:  strings.TrimSpace(req.Description),
			CoverURL:     strings.TrimSpace(req.CoverURL),
			CoverColor:   coverColor,
			AuthorID:     user.ID,
			Published:    true, // Mặc định publish sách mới để mọi người có thể xem
			BookTag:      strings.TrimSpace(req.BookTag),
			BookCategory: strings.TrimSpace(req.BookCategory),
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
			Title        string `json:"title"`
			Description  string `json:"description"`
			CoverURL     string `json:"cover_url"`
			CoverColor   string `json:"cover_color"`
			Published    bool   `json:"published"`
			BookTag      string `json:"book_tag"`
			BookCategory string `json:"book_category"`
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
		book.BookTag = strings.TrimSpace(req.BookTag)
		book.BookCategory = strings.TrimSpace(req.BookCategory)

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
			Title      string `json:"title"`
			Content    string `json:"content"`
			PageNumber int    `json:"page_number"`
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

		return c.JSON(page)
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
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
		}

		page.Title = strings.TrimSpace(req.Title)
		page.Content = req.Content

		if err := db.Save(&page).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi cập nhật trang"})
		}

		return c.JSON(fiber.Map{"success": true})
	}
}

// DeleteBook xóa cứng sách và tất cả nội dung liên quan
func DeleteBook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Chưa đăng nhập"})
		}

		db := database.Get()
		bookID, _ := strconv.Atoi(c.Params("id"))

		// Tìm sách
		var book models.Book
		if err := db.First(&book, bookID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy sách"})
		}

		// Kiểm tra quyền - chỉ author mới được xóa
		if book.AuthorID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Bạn không có quyền xóa sách này"})
		}

		// Xóa trong transaction để đảm bảo data integrity
		err := db.Transaction(func(tx *gorm.DB) error {
			// 1. Lấy tất cả pages của sách
			var pages []models.BookPage
			if err := tx.Where("book_id = ?", bookID).Find(&pages).Error; err != nil {
				return err
			}

			// 2. Xóa tất cả highlights và annotations của từng page
			for _, page := range pages {
				// Xóa highlights
				if err := tx.Unscoped().Where("book_page_id = ?", page.ID).Delete(&models.Highlight{}).Error; err != nil {
					return err
				}
			}

			// 3. Xóa tất cả pages
			if err := tx.Unscoped().Where("book_id = ?", bookID).Delete(&models.BookPage{}).Error; err != nil {
				return err
			}

			// 4. Xóa sách
			if err := tx.Unscoped().Delete(&book).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("Error deleting book %d: %v", bookID, err)
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi khi xóa sách"})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Đã xóa sách thành công",
		})
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

// SaveHighlight lưu highlight mới
func SaveHighlight() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Chưa đăng nhập"})
		}

		db := database.Get()
		bookID, _ := strconv.Atoi(c.Params("bookId"))
		pageID, _ := strconv.Atoi(c.Params("pageId"))

			// Verify book ownership/access
		var book models.Book
		if err := db.First(&book, bookID).Error; err != nil {
			log.Printf("Book not found: %v", err)
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy sách"})
		}

		// Verify page belongs to book
		var page models.BookPage
		if err := db.Where("id = ? AND book_id = ?", pageID, bookID).First(&page).Error; err != nil {
			log.Printf("Page not found: pageID=%d, bookID=%d, error=%v", pageID, bookID, err)
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy trang"})
		}

			var payload struct {
			Color           string `json:"color"`
			HighlightedText string `json:"highlighted_text"`
			Note            string `json:"note"`
			StartOffset     int    `json:"start_offset"`
			EndOffset       int    `json:"end_offset"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
		}

		highlight := models.Highlight{
			BookPageID:      page.ID,
			UserID:          user.ID,
			Color:           payload.Color,
			HighlightedText: payload.HighlightedText,
			Note:            payload.Note,
			StartOffset:     payload.StartOffset,
			EndOffset:       payload.EndOffset,
		}

		if err := db.Create(&highlight).Error; err != nil {
			log.Printf("Error creating highlight: %v, PageID: %d, UserID: %d", err, page.ID, user.ID)
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi lưu highlight", "details": err.Error()})
		}

		return c.JSON(highlight)
	}
}

// GetHighlights lấy highlights của user hiện tại cho một page
func GetHighlights() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		bookID, _ := strconv.Atoi(c.Params("bookId"))
		pageID, _ := strconv.Atoi(c.Params("pageId"))

		// Get book to check author
		var book models.Book
		if err := db.First(&book, bookID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy sách"})
		}

		// Get current user (may be nil if not logged in)
		user := getUserForBooks(c)

		// Load highlights based on user role:
		// - Not logged in: see only author's highlights
		// - Author: see only their own highlights
		// - Other users: see author's highlights + their own highlights
		var highlights []models.Highlight
		
		if user == nil {
			// Not logged in: load only author's highlights
			if err := db.Where("book_page_id = ? AND user_id = ?", pageID, book.AuthorID).Find(&highlights).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Lỗi tải highlights"})
			}
		} else if book.AuthorID == user.ID {
			// Author: load only their own highlights
			if err := db.Where("book_page_id = ? AND user_id = ?", pageID, user.ID).Find(&highlights).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Lỗi tải highlights"})
			}
		} else {
			// Non-author logged in: load author's highlights + own highlights
			if err := db.Where("book_page_id = ? AND (user_id = ? OR user_id = ?)", pageID, book.AuthorID, user.ID).Find(&highlights).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Lỗi tải highlights"})
			}
		}

		// Return highlights
		type HighlightResponse struct {
			ID              uint   `json:"id"`
			Color           string `json:"color"`
			HighlightedText string `json:"highlighted_text"`
			Note            string `json:"note"`
			StartOffset     int    `json:"start_offset"`
			EndOffset       int    `json:"end_offset"`
		}

		// Initialize as empty slice instead of nil to ensure JSON returns [] not null
		response := make([]HighlightResponse, 0)
		for _, h := range highlights {
			response = append(response, HighlightResponse{
				ID:              h.ID,
				Color:           h.Color,
				HighlightedText: h.HighlightedText,
				Note:            h.Note,
				StartOffset:     h.StartOffset,
				EndOffset:       h.EndOffset,
			})
		}

		return c.JSON(response)
	}
}

// DeleteHighlight xóa highlight
func DeleteHighlight() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getUserForBooks(c)
		if user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Chưa đăng nhập"})
		}

		db := database.Get()
		highlightID, _ := strconv.Atoi(c.Params("highlightId"))

		var highlight models.Highlight
		if err := db.First(&highlight, highlightID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy highlight"})
		}

		if highlight.UserID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Không có quyền"})
		}

		if err := db.Delete(&highlight).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi xóa highlight"})
		}

		return c.JSON(fiber.Map{"success": true})
	}
}

// SearchBooks API endpoint for searching books
func SearchBooks() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		user := getUserForBooks(c)

		query := c.Query("q")
		if query == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Query parameter 'q' is required"})
		}

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		searchTerm := "%" + strings.ToLower(query) + "%"

		// Build query với full-text search
		var books []models.Book
		var total int64

		baseQuery := db.Model(&models.Book{})

		// Chỉ tìm sách published hoặc sách của mình
		if user != nil {
			baseQuery = baseQuery.Where("(published = ? OR author_id = ?)", true, user.ID)
		} else {
			baseQuery = baseQuery.Where("published = ?", true)
		}

		// Search trong title, description, và author name
		searchQuery := baseQuery.
			Joins("LEFT JOIN users ON users.id = books.author_id").
			Where("LOWER(books.title) LIKE ? OR LOWER(books.description) LIKE ? OR LOWER(users.name) LIKE ?",
				searchTerm, searchTerm, searchTerm)

		// Count total results
		searchQuery.Count(&total)

		// Get paginated results
		if err := searchQuery.
			Order("books.created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&books).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Lỗi tìm kiếm"})
		}

		// Load author names
		for i := range books {
			var author models.User
			if err := db.First(&author, books[i].AuthorID).Error; err == nil {
				books[i].AuthorName = author.Name
			}
		}

		return c.JSON(fiber.Map{
			"success": true,
			"results": books,
			"total":   total,
			"page":    page,
			"limit":   limit,
			"query":   query,
		})
	}
}
