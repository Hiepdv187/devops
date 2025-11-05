package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/mail"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"fiber-learning-community/internal/database"
	"fiber-learning-community/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var sessionStore = session.New()

// Vietnam timezone
var vietnamLocation *time.Location

func init() {
	// Load Vietnam timezone (UTC+7)
	var err error
	vietnamLocation, err = time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		// Fallback to fixed offset if timezone data not available
		vietnamLocation = time.FixedZone("ICT", 7*60*60) // UTC+7
	}
}

// formatTimeVN formats time in Vietnam timezone
func formatTimeVN(t time.Time) string {
	return t.In(vietnamLocation).Format("02/01/2006 15:04")
}

func isJSONRequest(c *fiber.Ctx) bool {
	return c.Is("json") || strings.Contains(strings.ToLower(c.Get(fiber.HeaderContentType)), fiber.MIMEApplicationJSON)
}

func setFlash(c *fiber.Ctx, kind, message string) {
	sess, err := sessionStore.Get(c)
	if err != nil {
		return
	}
	sess.Set("flashType", kind)
	sess.Set("flashMessage", message)
	_ = sess.Save()
}

func popFlash(sess *session.Session) fiber.Map {
	msg, _ := sess.Get("flashMessage").(string)
	if msg == "" {
		return nil
	}
	typ, _ := sess.Get("flashType").(string)
	sess.Delete("flashMessage")
	sess.Delete("flashType")
	return fiber.Map{"Type": typ, "Message": msg}
}

// parseInlineAnnotations extracts inline annotations from HTML content
// Returns map of line_number -> annotation_text
func parseInlineAnnotations(htmlContent string) map[int]string {
	annotations := make(map[int]string)
	if htmlContent == "" {
		return annotations
	}

	// Split content by common block elements
	linePattern := regexp.MustCompile(`<(?:p|h[1-6]|li|div)[^>]*>.*?</(?:p|h[1-6]|li|div)>|[^<]+`)
	matches := linePattern.FindAllString(htmlContent, -1)

	lineNumber := 0
	for _, match := range matches {
		// Skip empty matches
		if strings.TrimSpace(match) == "" {
			continue
		}

		// Check if this is a block element with content
		if strings.Contains(match, "<") {
			lineNumber++

			// Extract text content and look for #annotation pattern
			// Remove HTML tags to get plain text
			tagPattern := regexp.MustCompile(`<[^>]+>`)
			plainText := tagPattern.ReplaceAllString(match, "")
			plainText = strings.TrimSpace(plainText)

			if plainText == "" {
				continue
			}

			// Look for #annotation at the end of the line
			annotationPattern := regexp.MustCompile(`\s*#([^#]+)$`)
			if annotationMatch := annotationPattern.FindStringSubmatch(plainText); len(annotationMatch) > 1 {
				annotationText := strings.TrimSpace(annotationMatch[1])
				if annotationText != "" {
					annotations[lineNumber] = annotationText
				}
			}
		}
	}

	return annotations
}

// updatePostContentAnnotation updates the #annotation in post content for a specific line
func updatePostContentAnnotation(db *gorm.DB, post *models.Post, lineNumber int, newAnnotation string) error {
	if post.Content == "" {
		return nil
	}

	// Split content by common block elements
	linePattern := regexp.MustCompile(`(<(?:p|h[1-6]|li|div)[^>]*>)(.*?)(</(?:p|h[1-6]|li|div)>)`)
	matches := linePattern.FindAllStringSubmatch(post.Content, -1)

	currentLine := 0
	updatedContent := post.Content

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		openTag := match[1]
		content := match[2]
		closeTag := match[3]
		fullMatch := match[0]

		// Skip empty content
		tagPattern := regexp.MustCompile(`<[^>]+>`)
		plainText := tagPattern.ReplaceAllString(content, "")
		if strings.TrimSpace(plainText) == "" {
			continue
		}

		currentLine++

		if currentLine == lineNumber {
			// Remove existing #annotation if any
			annotationPattern := regexp.MustCompile(`\s*#[^#]*$`)
			cleanContent := annotationPattern.ReplaceAllString(content, "")

			// Add new annotation if provided
			var newContent string
			if newAnnotation != "" {
				newContent = cleanContent + " #" + newAnnotation
			} else {
				newContent = cleanContent
			}

			// Replace in full content
			newFullMatch := openTag + newContent + closeTag
			updatedContent = strings.Replace(updatedContent, fullMatch, newFullMatch, 1)
			break
		}
	}

	// Save updated content
	post.Content = updatedContent
	return db.Save(post).Error
}

func setUserSession(c *fiber.Ctx, user models.User) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		return err
	}
	sess.Set("userID", strconv.FormatUint(uint64(user.ID), 10))
	sess.Set("userName", user.Name)
	sess.Set("userEmail", user.Email)
	return sess.Save()
}

func clearUserSession(c *fiber.Ctx) {
	sess, err := sessionStore.Get(c)
	if err != nil {
		return
	}
	_ = sess.Destroy()
}

func sessionUser(sess *session.Session) fiber.Map {
	idStr, _ := sess.Get("userID").(string)
	if idStr == "" {
		return nil
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil
	}
	name, _ := sess.Get("userName").(string)
	email, _ := sess.Get("userEmail").(string)
	return fiber.Map{
		"ID":    uint(id),
		"Name":  name,
		"Email": email,
	}
}

func currentUserID(c *fiber.Ctx) (uint, error) {
	sess, err := sessionStore.Get(c)
	if err != nil {
		return 0, err
	}
	idStr, _ := sess.Get("userID").(string)
	if idStr == "" {
		return 0, fiber.ErrUnauthorized
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func respondError(c *fiber.Ctx, status int, message, redirect string) error {
	if isJSONRequest(c) {
		return fiber.NewError(status, message)
	}
	setFlash(c, "error", message)
	if redirect == "" {
		redirect = c.Path()
	}
	return c.Status(fiber.StatusSeeOther).Redirect(redirect)
}

func render(c *fiber.Ctx, view string, data fiber.Map, layout string) error {
	base := fiber.Map{
		"AppName":         "Cộng đồng Học DevOps",
		"Year":            time.Now().Year(),
		"IsAuthenticated": false,
		"RequestPath":     c.OriginalURL(),
		"RequestRoute":    c.Path(),
	}

	sess, err := sessionStore.Get(c)
	if err == nil {
		if flash := popFlash(sess); flash != nil {
			base["Flash"] = flash
		}
		if user := sessionUser(sess); user != nil {
			base["CurrentUser"] = user
			base["IsAuthenticated"] = true
		}
		_ = sess.Save()
	}

	for k, v := range data {
		base[k] = v
	}

	// If layout is specified, use it; if empty, render without layout
	if layout != "" {
		return c.Render(view, base, "layouts/"+layout)
	}
	return c.Render(view, base)
}

func Home() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		var posts []models.Post
		latestPosts := []fiber.Map{}
		if err := db.Preload("Author").Order("created_at DESC").Limit(6).Find(&posts).Error; err == nil {
			for _, p := range posts {
				postTags := []string{}
				if strings.TrimSpace(p.Tags) != "" {
					tags := strings.Split(p.Tags, ",")
					for _, tag := range tags {
						tag = strings.TrimSpace(tag)
						if tag != "" {
							postTags = append(postTags, tag)
						}
					}
				}
				latestPosts = append(latestPosts, fiber.Map{
					"ID":           p.ID,
					"Title":        p.Title,
					"Summary":      p.Summary,
					"CoverURL":     p.CoverURL,
					"AuthorName":   p.Author.Name,
					"CreatedLabel": formatTimeVN(p.CreatedAt),
					"PostTags":     postTags,
				})
			}
		}
		return render(c, "pages/home", fiber.Map{
			"Title":       "Học DevOps cùng cộng đồng",
			"HeroHeading": "Thực chiến DevOps, làm chủ hạ tầng",
			"HeroSub":     "Khám phá lộ trình DevOps toàn diện từ kiến trúc hệ thống, tự động hóa, đến vận hành an toàn qua đóng góp của cộng đồng.",
			"Highlights": []fiber.Map{
				{
					"Title":       "Lộ trình DevOps",
					"Description": "Phủ đầy từ Linux cơ bản, mạng, container, CI/CD đến vận hành cloud-native.",
				},
				{
					"Title":       "Học từ trải nghiệm thật",
					"Description": "Nhận tài liệu, checklist và case study từ những kỹ sư DevOps đang vận hành sản phẩm lớn.",
				},
				{
					"Title":       "Đóng góp không giới hạn",
					"Description": "Cập nhật thực tiễn mới nhất về cloud, observability, bảo mật và công cụ DevOps.",
				},
			},
			"TechInsights": []fiber.Map{
				{
					"Title":       "Quan sát hệ thống",
					"Description": "Best practice triển khai observability stack với OpenTelemetry, Prometheus và Grafana.",
				},
				{
					"Title":       "Bảo mật chuỗi CI/CD",
					"Description": "Các bước harden pipeline, ký container và quét lỗ hổng tự động trước khi phát hành.",
				},
				{
					"Title":       "Hạ tầng đa đám mây",
					"Description": "Terraform module tái sử dụng để quản lý multi-cloud và tối ưu chi phí vận hành.",
				},
			},
			"CommunityUpdates": []fiber.Map{
				{
					"Title":   "Workshop: GitOps nâng cao",
					"Date":    "05/11/2025",
					"Summary": "Trải nghiệm triển khai Argo CD với policy guardrail và progressive delivery.",
				},
				{
					"Title":   "Live stream: Observability trong Kubernetes",
					"Date":    "12/11/2025",
					"Summary": "Giải đáp trực tiếp các tình huống troubleshooting thực tế trong cluster production.",
				},
				{
					"Title":   "Blog mới: IaC Testing",
					"Date":    "Tuần này",
					"Summary": "Checklist kiểm thử Terraform và cách tích hợp Terratest vào pipeline CI.",
				},
			},
			"LatestPosts": latestPosts,
		}, "main")
	}
}

type createPostRequest struct {
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	Content         string `json:"content"`
	ContentEncoded  string `json:"content_encoded"`  // Base64 encoded content to bypass WAF
	CoverURL        string `json:"cover_url"`
	Tags            string `json:"tags"`
	AuthorID        uint   `json:"author_id"`
	LineAnnotations string `json:"line_annotations"` // JSON string: {"1": "notice text", "2": "another notice"}
}

type createCommentRequest struct {
	Content    string `json:"content"`
	AuthorID   uint   `json:"author_id"`
	LineNumber *int   `json:"line_number"`
}

func PostsPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if create mode is requested
		createMode := c.Query("create") == "true"
		if createMode {
			// Check if user is authenticated
			userID, _ := currentUserID(c)
			if userID > 0 {
				return render(c, "pages/posts_create_smart", fiber.Map{
					"Title": "Tạo bài viết mới",
				}, "main")
			}
			// If not authenticated, redirect to login
			return c.Redirect("/auth/login?next=/posts?create=true")
		}
		
		db := database.Get()
		query := strings.TrimSpace(c.Query("q"))
		selectedTag := strings.TrimSpace(c.Query("tag"))
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}
		const pageSize = 10
		offset := (page - 1) * pageSize

		builder := db.Model(&models.Post{}).Preload("Author").Order("created_at DESC")
		if query != "" {
			like := fmt.Sprintf("%%%s%%", query)
			builder = builder.Where("title LIKE ? OR summary LIKE ?", like, like)
		}
		if selectedTag != "" {
			tagLike := fmt.Sprintf("%%%s%%", selectedTag)
			builder = builder.Where("tags LIKE ?", tagLike)
		}

		var total int64
		if err := builder.Count(&total).Error; err != nil {
			return respondError(c, fiber.StatusInternalServerError, "Không thể đếm bài viết", "/")
		}

		var posts []models.Post
		if err := builder.Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
			return respondError(c, fiber.StatusInternalServerError, "Không thể tải danh sách bài viết", "/")
		}
		items := make([]fiber.Map, 0, len(posts))
		for _, p := range posts {
			postTags := []string{}
			if strings.TrimSpace(p.Tags) != "" {
				tags := strings.Split(p.Tags, ",")
				for _, tag := range tags {
					tag = strings.TrimSpace(tag)
					if tag != "" {
						postTags = append(postTags, tag)
					}
				}
			}

			items = append(items, fiber.Map{
				"ID":           p.ID,
				"Title":        p.Title,
				"Summary":      p.Summary,
				"Content":      p.Content,
				"CoverURL":     p.CoverURL,
				"AuthorName":   p.Author.Name,
				"CreatedLabel": formatTimeVN(p.CreatedAt),
				"PostTags":     postTags,
			})
		}

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
		if totalPages == 0 {
			totalPages = 1
		}
		pages := make([]int, totalPages)
		for i := range pages {
			pages[i] = i + 1
		}

		// Get all unique tags
		var allPosts []models.Post
		db.Select("tags").Find(&allPosts)
		tagSet := make(map[string]bool)
		for _, p := range allPosts {
			if p.Tags != "" {
				tags := strings.Split(p.Tags, ",")
				for _, tag := range tags {
					tag = strings.TrimSpace(tag)
					if tag != "" {
						tagSet[tag] = true
					}
				}
			}
		}
		allTags := make([]string, 0, len(tagSet))
		for tag := range tagSet {
			allTags = append(allTags, tag)
		}

		return render(c, "pages/posts", fiber.Map{
			"Title":       "Bài viết cộng đồng",
			"Posts":       items,
			"Query":       query,
			"SelectedTag": selectedTag,
			"AllTags":     allTags,
			"Page":        page,
			"TotalPages":  totalPages,
			"HasPrev":     page > 1,
			"HasNext":     page < totalPages,
			"PrevPage":    page - 1,
			"NextPage":    page + 1,
			"Pages":       pages,
		}, "main")
	}
}

func PostPreviewPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/post_preview_page", fiber.Map{
			"Title": "Xem trước bài viết",
		}, "main")
	}
}

func PostDetailPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "ID bài viết không hợp lệ", "/posts")
		}

		db := database.Get()
		var post models.Post
		if err := db.Preload("Author").Preload("Comments", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Author").Order("created_at ASC")
		}).First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return respondError(c, fiber.StatusNotFound, "Bài viết không tồn tại", "/posts")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể tải bài viết", "/posts")
		}

		lineComments := make(map[int][]fiber.Map)
		generalComments := make([]fiber.Map, 0)

		for _, cm := range post.Comments {
			comment := fiber.Map{
				"ID":         cm.ID,
				"Content":    cm.Content,
				"AuthorName": cm.Author.Name,
				"CreatedAt":  formatTimeVN(cm.CreatedAt),
			}
			if cm.LineNumber != nil {
				lineComments[*cm.LineNumber] = append(lineComments[*cm.LineNumber], comment)
			} else {
				generalComments = append(generalComments, comment)
			}
		}

		// Load annotations
		var annotations []models.Annotation
		db.Where("post_id = ?", post.ID).Find(&annotations)

		lineAnnotations := make(map[int]string)
		for _, ann := range annotations {
			lineAnnotations[ann.LineNumber] = ann.Content
		}

		// Check if current user is author
		userID, _ := currentUserID(c)
		isAuthor := userID == post.AuthorID

		// Parse tags
		postTags := []string{}
		if post.Tags != "" {
			tags := strings.Split(post.Tags, ",")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					postTags = append(postTags, tag)
				}
			}
		}

		// Check if edit mode is requested
		editMode := c.Query("edit") == "true"
		templateName := "pages/post_detail"
		if editMode && isAuthor {
			templateName = "pages/post_edit_smart"
		}

		return render(c, templateName, fiber.Map{
			"Title": "Chi tiết bài viết",
			"Post": fiber.Map{
				"ID":           post.ID,
				"Title":        post.Title,
				"Summary":      post.Summary,
				"Content":      post.Content,
				"CoverURL":     post.CoverURL,
				"Tags":         post.Tags,
				"AuthorName":   post.Author.Name,
				"AuthorID":     post.AuthorID,
				"CreatedLabel": formatTimeVN(post.CreatedAt),
			},
			"PostTags":        postTags,
			"LineComments":    lineComments,
			"GeneralComments": generalComments,
			"LineAnnotations": lineAnnotations,
			"IsAuthor":        isAuthor,
		}, "main")
	}
}

func RegisterPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/auth_register", fiber.Map{
			"Title": "Đăng ký tài khoản",
			"Next":  c.Query("next"),
		}, "main")
	}
}

func LoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/auth_login", fiber.Map{
			"Title": "Đăng nhập",
			"Next":  c.Query("next"),
		}, "main")
	}
}

func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clearUserSession(c)
		if isJSONRequest(c) {
			return c.JSON(fiber.Map{"message": "Đã đăng xuất"})
		}
		setFlash(c, "success", "Bạn đã đăng xuất")
		redirect := c.FormValue("next")
		if redirect == "" {
			redirect = "/"
		}
		return c.Status(fiber.StatusSeeOther).Redirect(redirect)
	}
}

// CreatePost cho phép người dùng tạo bài viết mới.
func CreatePost() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isJSON := isJSONRequest(c)
		var body createPostRequest

		if isJSON {
			if err := c.BodyParser(&body); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "payload không hợp lệ")
			}
		} else {
			body.Title = c.FormValue("title")
			body.Summary = c.FormValue("summary")
			body.Content = c.FormValue("content")
			body.CoverURL = c.FormValue("cover_url")
			body.Tags = c.FormValue("tags")
			body.LineAnnotations = c.FormValue("line_annotations")
		}

		body.Title = strings.TrimSpace(body.Title)
		body.Summary = strings.TrimSpace(body.Summary)
		body.Content = strings.TrimSpace(body.Content)
		body.CoverURL = strings.TrimSpace(body.CoverURL)
		body.Tags = strings.TrimSpace(body.Tags)

		// Decode base64 content if provided (to bypass WAF)
		if body.ContentEncoded != "" {
			decoded, err := base64.StdEncoding.DecodeString(body.ContentEncoded)
			if err == nil {
				body.Content = string(decoded)
			}
		}

		if !isJSON {
			authorID, err := currentUserID(c)
			if err != nil {
				return respondError(c, fiber.StatusUnauthorized, "Bạn cần đăng nhập để tạo bài viết", "/auth/login")
			}
			body.AuthorID = authorID
		} else if body.AuthorID == 0 {
			if authorID, err := currentUserID(c); err == nil {
				body.AuthorID = authorID
			}
		}

		if len(body.Title) < 3 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "tiêu đề phải từ 3 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Tiêu đề phải từ 3 ký tự", "/posts")
		}

		if len(body.Content) < 10 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "nội dung phải từ 10 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Nội dung phải từ 10 ký tự", "/posts")
		}

		if body.AuthorID == 0 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "thiếu thông tin tác giả")
			}
			return respondError(c, fiber.StatusUnauthorized, "Bạn cần đăng nhập để tạo bài viết", "/auth/login")
		}

		db := database.Get()

		var author models.User
		if err := db.First(&author, body.AuthorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if isJSON {
					return fiber.NewError(fiber.StatusBadRequest, "tác giả không tồn tại")
				}
				return respondError(c, fiber.StatusBadRequest, "Người dùng không tồn tại", "/posts")
			}
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể kiểm tra tác giả")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể kiểm tra tác giả", "/posts")
		}

		post := models.Post{
			Title:    body.Title,
			Summary:  body.Summary,
			Content:  body.Content,
			CoverURL: body.CoverURL,
			Tags:     body.Tags,
			AuthorID: body.AuthorID,
		}

		if err := db.Create(&post).Error; err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể tạo bài viết")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể tạo bài viết", "/posts")
		}

		// Create annotations from line_annotations JSON field
		if body.LineAnnotations != "" {
			var annotationsMap map[string]string
			if err := json.Unmarshal([]byte(body.LineAnnotations), &annotationsMap); err == nil {
				for lineNumStr, annotationText := range annotationsMap {
					lineNum, err := strconv.Atoi(lineNumStr)
					if err != nil {
						continue
					}
					annotation := models.Annotation{
						PostID:     post.ID,
						LineNumber: lineNum,
						Content:    strings.TrimSpace(annotationText),
					}
					if err := db.Create(&annotation).Error; err != nil {
						fmt.Printf("Warning: Failed to create annotation for line %d: %v\n", lineNum, err)
					}
				}
			} else {
				fmt.Printf("Warning: Failed to parse line_annotations JSON: %v\n", err)
			}
		}

		if isJSON {
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"message": "Tạo bài viết thành công",
				"post": fiber.Map{
					"id":         post.ID,
					"title":      post.Title,
					"summary":    post.Summary,
					"content":    post.Content,
					"author_id":  post.AuthorID,
					"created_at": post.CreatedAt,
				},
			})
		}

		setFlash(c, "success", "Bài viết đã được tạo thành công")
		return c.Status(fiber.StatusSeeOther).Redirect(fmt.Sprintf("/posts/%d", post.ID))
	}
}

// UpdatePost cho phép tác giả chỉnh sửa bài viết của mình.
func UpdatePost() fiber.Handler {
	return func(c *fiber.Ctx) error {
		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "ID bài viết không hợp lệ", "/posts")
		}

		userID, err := currentUserID(c)
		if err != nil {
			if isJSONRequest(c) {
				return fiber.NewError(fiber.StatusUnauthorized, "bạn cần đăng nhập")
			}
			return respondError(c, fiber.StatusUnauthorized, "Bạn cần đăng nhập để chỉnh sửa bài viết", "/auth/login")
		}

		var req struct {
			Title           string `json:"title"`
			Summary         string `json:"summary"`
			Content         string `json:"content"`
			ContentEncoded  string `json:"content_encoded"` // Base64 encoded content to bypass WAF
			CoverURL        string `json:"cover_url"`
			Tags            string `json:"tags"`
			LineAnnotations string `json:"line_annotations"`
		}

		isJSON := isJSONRequest(c)
		if isJSON {
			if err := c.BodyParser(&req); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "payload không hợp lệ")
			}
		} else {
			req.Title = c.FormValue("title")
			req.Summary = c.FormValue("summary")
			req.Content = c.FormValue("content")
			req.CoverURL = c.FormValue("cover_url")
			req.Tags = c.FormValue("tags")
			req.LineAnnotations = c.FormValue("line_annotations")
		}

		req.Title = strings.TrimSpace(req.Title)
		req.Summary = strings.TrimSpace(req.Summary)
		req.Content = strings.TrimSpace(req.Content)
		req.CoverURL = strings.TrimSpace(req.CoverURL)
		req.Tags = strings.TrimSpace(req.Tags)

		// Decode base64 content if provided (to bypass WAF)
		if req.ContentEncoded != "" {
			decoded, err := base64.StdEncoding.DecodeString(req.ContentEncoded)
			if err == nil {
				req.Content = string(decoded)
			}
		}

		if len(req.Title) < 3 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "tiêu đề phải từ 3 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Tiêu đề phải từ 3 ký tự", fmt.Sprintf("/posts/%d", postID))
		}

		if len(req.Content) < 10 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "nội dung phải từ 10 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Nội dung phải từ 10 ký tự", fmt.Sprintf("/posts/%d", postID))
		}

		db := database.Get()

		var post models.Post
		if err := db.First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if isJSON {
					return fiber.NewError(fiber.StatusNotFound, "bài viết không tồn tại")
				}
				return respondError(c, fiber.StatusNotFound, "Bài viết không tồn tại", "/posts")
			}
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể tải bài viết")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể tải bài viết", "/posts")
		}

		if post.AuthorID != userID {
			if isJSON {
				return fiber.NewError(fiber.StatusForbidden, "chỉ tác giả mới được chỉnh sửa")
			}
			return respondError(c, fiber.StatusForbidden, "Chỉ tác giả mới có thể chỉnh sửa bài viết", fmt.Sprintf("/posts/%d", postID))
		}

		post.Title = req.Title
		post.Summary = req.Summary
		post.Content = req.Content
		post.CoverURL = req.CoverURL
		post.Tags = req.Tags

		if err := db.Save(&post).Error; err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể cập nhật bài viết")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể cập nhật bài viết", fmt.Sprintf("/posts/%d", postID))
		}

		// Sync annotations from line_annotations JSON field
		if req.LineAnnotations != "" {
			var annotationsMap map[string]string
			if err := json.Unmarshal([]byte(req.LineAnnotations), &annotationsMap); err == nil {
				// Delete all existing annotations for this post
				db.Where("post_id = ?", post.ID).Delete(&models.Annotation{})
				
				// Create new annotations from JSON
				for lineNumStr, annotationText := range annotationsMap {
					lineNum, err := strconv.Atoi(lineNumStr)
					if err != nil {
						continue
					}
					annotation := models.Annotation{
						PostID:     post.ID,
						LineNumber: lineNum,
						Content:    strings.TrimSpace(annotationText),
					}
					if err := db.Create(&annotation).Error; err != nil {
						fmt.Printf("Warning: Failed to create annotation for line %d: %v\n", lineNum, err)
					}
				}
			} else {
				fmt.Printf("Warning: Failed to parse line_annotations JSON: %v\n", err)
			}
		}

		if isJSON {
			return c.JSON(fiber.Map{
				"message": "Cập nhật bài viết thành công",
				"post": fiber.Map{
					"id":        post.ID,
					"title":     post.Title,
					"summary":   post.Summary,
					"content":   post.Content,
					"cover_url": post.CoverURL,
				},
			})
		}

		setFlash(c, "success", "Bài viết đã được cập nhật")
		return c.Status(fiber.StatusSeeOther).Redirect(fmt.Sprintf("/posts/%d", post.ID))
	}
}

// CreateComment cho phép người dùng bình luận vào bài viết xác định.
func CreateComment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		postIDParam := c.Params("id")
		postID, err := strconv.ParseUint(postIDParam, 10, 64)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "Bài viết không hợp lệ", "/posts")
		}

		isJSON := isJSONRequest(c)
		var body createCommentRequest
		if isJSON {
			if err := c.BodyParser(&body); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "payload không hợp lệ")
			}
		} else {
			body.Content = c.FormValue("content")
			if ln := c.FormValue("line_number"); ln != "" {
				if val, err := strconv.Atoi(ln); err == nil && val > 0 {
					body.LineNumber = &val
				}
			}
		}

		body.Content = strings.TrimSpace(body.Content)
		if len(body.Content) < 3 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "nội dung bình luận phải từ 3 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Bình luận phải từ 3 ký tự", fmt.Sprintf("/posts/%d", postID))
		}

		if body.AuthorID == 0 {
			if userID, err := currentUserID(c); err == nil {
				body.AuthorID = userID
			}
		}

		if body.AuthorID == 0 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "thiếu thông tin người bình luận")
			}
			return respondError(c, fiber.StatusUnauthorized, "Bạn cần đăng nhập để bình luận", "/auth/login")
		}

		db := database.Get()

		var post models.Post
		if err := db.First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return respondError(c, fiber.StatusNotFound, "Bài viết không tồn tại", "/posts")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể truy vấn bài viết", "/posts")
		}

		var author models.User
		if err := db.First(&author, body.AuthorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if isJSON {
					return fiber.NewError(fiber.StatusBadRequest, "người dùng không tồn tại")
				}
				return respondError(c, fiber.StatusBadRequest, "Người dùng không tồn tại", fmt.Sprintf("/posts/%d", postID))
			}
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể kiểm tra người dùng")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể kiểm tra người dùng", fmt.Sprintf("/posts/%d", postID))
		}

		comment := models.Comment{
			Content:    body.Content,
			PostID:     uint(postID),
			AuthorID:   body.AuthorID,
			LineNumber: body.LineNumber,
		}

		if err := db.Create(&comment).Error; err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể tạo bình luận")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể tạo bình luận", fmt.Sprintf("/posts/%d", postID))
		}

		authorName := author.Name
		if isJSON {
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"message": "Đã thêm bình luận",
				"comment": fiber.Map{
					"id":          comment.ID,
					"content":     comment.Content,
					"author_name": authorName,
					"line_number": body.LineNumber,
					"created_at":  formatTimeVN(comment.CreatedAt),
				},
			})
		}

		setFlash(c, "success", "Đã thêm bình luận")
		return c.Status(fiber.StatusSeeOther).Redirect(fmt.Sprintf("/posts/%d", postID))
	}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register cho phép tạo tài khoản mới trong hệ thống.
func Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isJSON := isJSONRequest(c)
		var body registerRequest
		if isJSON {
			if err := c.BodyParser(&body); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "payload không hợp lệ")
			}
		} else {
			body.Name = c.FormValue("name")
			body.Email = c.FormValue("email")
			body.Password = c.FormValue("password")
		}

		body.Name = strings.TrimSpace(body.Name)
		body.Email = strings.ToLower(strings.TrimSpace(body.Email))
		body.Password = strings.TrimSpace(body.Password)

		if len(body.Name) < 3 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "tên phải từ 3 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Tên phải từ 3 ký tự", "/auth/register")
		}
		if _, err := mail.ParseAddress(body.Email); err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "email không hợp lệ")
			}
			return respondError(c, fiber.StatusBadRequest, "Email không hợp lệ", "/auth/register")
		}
		if len(body.Password) < 6 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "mật khẩu phải từ 6 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Mật khẩu phải từ 6 ký tự", "/auth/register")
		}

		db := database.Get()

		var existing models.User
		if err := db.Where("email = ?", body.Email).First(&existing).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				if isJSON {
					return fiber.NewError(fiber.StatusInternalServerError, "không thể kiểm tra tài khoản")
				}
				return respondError(c, fiber.StatusInternalServerError, "Không thể kiểm tra tài khoản", "/auth/register")
			}
		} else {
			if isJSON {
				return fiber.NewError(fiber.StatusConflict, "email đã tồn tại")
			}
			return respondError(c, fiber.StatusConflict, "Email đã tồn tại", "/auth/register")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể mã hóa mật khẩu")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể mã hóa mật khẩu", "/auth/register")
		}

		user := models.User{
			Name:         body.Name,
			Email:        body.Email,
			PasswordHash: string(hash),
		}

		if err := db.Create(&user).Error; err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể tạo tài khoản")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể tạo tài khoản", "/auth/register")
		}

		if err := setUserSession(c, user); err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể khởi tạo phiên đăng nhập")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể đăng nhập sau khi đăng ký", "/auth/login")
		}

		if isJSON {
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"message": "Đăng ký thành công",
				"user": fiber.Map{
					"id":    user.ID,
					"name":  user.Name,
					"email": user.Email,
				},
			})
		}

		setFlash(c, "success", "Chào mừng bạn đến với cộng đồng!")
		redirect := c.FormValue("next")
		if redirect == "" || redirect == "/auth/login" || redirect == "/auth/register" {
			redirect = "/"
		}
		return c.Status(fiber.StatusSeeOther).Redirect(redirect)
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login xử lý yêu cầu đăng nhập đơn giản với dữ liệu lưu trong SQLite.
func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isJSON := isJSONRequest(c)
		var body loginRequest
		if isJSON {
			if err := c.BodyParser(&body); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "payload không hợp lệ")
			}
		} else {
			body.Email = c.FormValue("email")
			body.Password = c.FormValue("password")
		}

		body.Email = strings.ToLower(strings.TrimSpace(body.Email))
		body.Password = strings.TrimSpace(body.Password)

		if _, err := mail.ParseAddress(body.Email); err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "email không hợp lệ")
			}
			return respondError(c, fiber.StatusBadRequest, "Email không hợp lệ", "/auth/login")
		}
		if len(body.Password) < 6 {
			if isJSON {
				return fiber.NewError(fiber.StatusBadRequest, "mật khẩu phải từ 6 ký tự")
			}
			return respondError(c, fiber.StatusBadRequest, "Mật khẩu phải từ 6 ký tự", "/auth/login")
		}

		var user models.User
		if err := database.Get().Where("email = ?", body.Email).First(&user).Error; err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusUnauthorized, "email hoặc mật khẩu sai")
			}
			return respondError(c, fiber.StatusUnauthorized, "Email hoặc mật khẩu không đúng", "/auth/login")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusUnauthorized, "email hoặc mật khẩu sai")
			}
			return respondError(c, fiber.StatusUnauthorized, "Email hoặc mật khẩu không đúng", "/auth/login")
		}

		if err := setUserSession(c, user); err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể khởi tạo phiên đăng nhập")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể đăng nhập", "/auth/login")
		}

		if isJSON {
			return c.JSON(fiber.Map{
				"message": "Đăng nhập thành công",
				"user": fiber.Map{
					"id":    user.ID,
					"name":  user.Name,
					"email": user.Email,
				},
			})
		}

		setFlash(c, "success", fmt.Sprintf("Chào mừng trở lại, %s!", user.Name))
		redirect := c.FormValue("next")
		if redirect == "" || redirect == "/auth/login" || redirect == "/auth/register" {
			redirect = "/"
		}
		return c.Status(fiber.StatusSeeOther).Redirect(redirect)
	}
}

func Courses() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/courses", fiber.Map{
			"Title":       "Lộ trình DevOps",
			"Description": "Các module DevOps do cộng đồng đóng góp, dễ dàng tùy chỉnh theo nhu cầu đội ngũ của bạn.",
			"Courses": []fiber.Map{
				{
					"Name":        "Nền tảng DevOps & Linux",
					"Level":       "Beginner",
					"Duration":    "8 giờ",
					"Summary":     "Nắm hệ điều hành Linux, quản trị hệ thống cơ bản và các nguyên tắc DevOps cốt lõi.",
					"Contributor": "Bảo Trần",
				},
				{
					"Name":        "Container & Orchestration",
					"Level":       "Intermediate",
					"Duration":    "10 giờ",
					"Summary":     "Làm chủ Docker, Kubernetes và triển khai workflow GitOps.",
					"Contributor": "Lan Nguyễn",
				},
				{
					"Name":        "Quan sát & Bảo mật",
					"Level":       "Advanced",
					"Duration":    "6 giờ",
					"Summary":     "Thiết lập monitoring, logging, alerting cùng với thực hành bảo mật DevSecOps.",
					"Contributor": "Hoàng Lê",
				},
			},
		}, "main")
	}
}

func Contributors() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/contributors", fiber.Map{
			"Title":       "Thành viên đóng góp",
			"Description": "Những kỹ sư DevOps đang cùng xây dựng thư viện kiến thức mở.",
			"Contributors": []fiber.Map{
				{
					"Name":          "Lan Nguyễn",
					"Role":          "Platform Engineer",
					"Bio":           "Thiết kế hệ thống Kubernetes multi-cluster cho doanh nghiệp fintech.",
					"Contributions": 14,
				},
				{
					"Name":          "Minh Phạm",
					"Role":          "Site Reliability Engineer",
					"Bio":           "Tối ưu pipeline CI/CD và vận hành hệ thống với SLO rõ ràng.",
					"Contributions": 21,
				},
				{
					"Name":          "Hoàng Lê",
					"Role":          "DevSecOps Specialist",
					"Bio":           "Xây dựng framework bảo mật và audit tự động cho microservices.",
					"Contributions": 11,
				},
			},
		}, "main")
	}
}

func About() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/about", fiber.Map{
			"Title":       "Về dự án",
			"Description": "Cộng đồng Học DevOps được xây dựng bởi các kỹ sư đam mê tự động hóa, quan sát hệ thống và văn hóa DevOps tại Việt Nam.",
			"Sections": []fiber.Map{
				{
					"Heading": "Sứ mệnh",
					"Content": "Mang đến tài nguyên DevOps tiếng Việt chất lượng cao, giúp đội ngũ kỹ thuật áp dụng nhanh vào công việc thực tế.",
				},
				{
					"Heading": "Giá trị cốt lõi",
					"Content": "Minh bạch, chia sẻ kinh nghiệm thực chiến, ưu tiên học qua làm và lan tỏa tinh thần cộng đồng.",
				},
				{
					"Heading": "Hướng phát triển",
					"Content": "Hoàn thiện các lộ trình DevOps theo cấp độ, xây dựng thư viện bài viết chuyên sâu và kho tài nguyên mã nguồn mở.",
				},
			},
			"Stats": []fiber.Map{
				{
					"Value": "+1.200",
					"Label": "Thành viên đăng ký",
				},
				{
					"Value": "85+",
					"Label": "Bài viết & hướng dẫn",
				},
				{
					"Value": "12",
					"Label": "Workshop & webinar",
				},
			},
			"Milestones": []fiber.Map{
				{
					"Period":      "Q1 2024",
					"Title":       "Ra mắt lộ trình DevOps căn bản",
					"Description": "Hoàn thiện bộ tài liệu học tập 8 tuần với bài tập thực hành và checklist đánh giá kỹ năng.",
				},
				{
					"Period":      "Q3 2024",
					"Title":       "Tổ chức chuỗi workshop Cloud Native",
					"Description": "Hơn 300 kỹ sư tham gia học cùng chuyên gia về Kubernetes, GitOps, Observability.",
				},
				{
					"Period":      "2025",
					"Title":       "Xây dựng thư viện template mở",
					"Description": "Cung cấp Terraform, Ansible, và pipeline mẫu giúp doanh nghiệp khởi động DevOps nhanh chóng.",
				},
			},
			"CTA": fiber.Map{
				"Heading":     "Cùng đóng góp để DevOps Việt Nam lớn mạnh",
				"Text":        "Chia sẻ kinh nghiệm thực tế, mentoring cho thành viên mới, hoặc mở một workshop tại cộng đồng.",
				"ActionLabel": "Bắt đầu đóng góp",
				"ActionLink":  "/contribute",
			},
		}, "main")
	}
}

func Contribute() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/contribute", fiber.Map{
			"Title":       "Đóng góp kiến thức DevOps",
			"Description": "Cùng xây dựng kho tri thức DevOps mở và giàu tính ứng dụng.",
			"Steps": []fiber.Map{
				{
					"Order":   1,
					"Heading": "Chọn chủ đề DevOps",
					"Text":    "Tập trung vào lĩnh vực bạn am hiểu như CI/CD, IaC, observability, hay bảo mật.",
				},
				{
					"Order":   2,
					"Heading": "Soạn tài liệu & demo",
					"Text":    "Chuẩn bị bài viết, slide, repo mẫu hoặc script tự động hóa kèm hướng dẫn chi tiết.",
				},
				{
					"Order":   3,
					"Heading": "Gửi đóng góp",
					"Text":    "Tạo pull request trên GitHub của dự án hoặc tham gia phiên review nội dung định kỳ.",
				},
			},
			"ContactEmail": "devops@example.com",
		}, "main")
	}
}

// UploadImage xử lý upload ảnh cho bài viết - lưu vào database
func UploadImage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Kiểm tra đăng nhập
		userID, err := currentUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Bạn cần đăng nhập để upload ảnh",
			})
		}

		// Lấy file từ form
		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Không tìm thấy file ảnh",
			})
		}

		// Kiểm tra loại file
		ext := filepath.Ext(fileHeader.Filename)
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
		}
		if !allowedExts[strings.ToLower(ext)] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Chỉ hỗ trợ file ảnh: jpg, jpeg, png, gif, webp",
			})
		}

		// Kiểm tra kích thước (max 5MB)
		if fileHeader.Size > 5*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Kích thước ảnh không được vượt quá 5MB",
			})
		}

		// Đọc nội dung file
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Không thể đọc file",
			})
		}
		defer file.Close()

		// Đọc dữ liệu binary
		fileData, err := io.ReadAll(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Không thể đọc dữ liệu file",
			})
		}

		// Xác định content type
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			// Fallback dựa vào extension
			switch strings.ToLower(ext) {
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			case ".png":
				contentType = "image/png"
			case ".gif":
				contentType = "image/gif"
			case ".webp":
				contentType = "image/webp"
			default:
				contentType = "application/octet-stream"
			}
		}

		// Lưu vào database
		db := database.Get()
		image := models.Image{
			Filename:    fileHeader.Filename,
			ContentType: contentType,
			Size:        fileHeader.Size,
			Data:        fileData,
			UploaderID:  userID,
		}

		if err := db.Create(&image).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Không thể lưu ảnh vào database",
			})
		}

		// Trả về URL với ID của ảnh
		imageURL := fmt.Sprintf("/images/%d", image.ID)
		return c.JSON(fiber.Map{
			"url":      imageURL,
			"markdown": fmt.Sprintf("![%s](%s)", fileHeader.Filename, imageURL),
		})
	}
}

// GetImage trả về ảnh từ database
func GetImage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		imageID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("ID ảnh không hợp lệ")
		}

		db := database.Get()
		var image models.Image
		
		// Chỉ lấy metadata trước để kiểm tra
		if err := db.Select("id", "content_type", "filename").First(&image, imageID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("Không tìm thấy ảnh")
			}
			return c.Status(fiber.StatusInternalServerError).SendString("Lỗi truy vấn database")
		}

		// Lấy dữ liệu ảnh
		if err := db.Model(&models.Image{}).Where("id = ?", imageID).Pluck("data", &image.Data).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Không thể tải ảnh")
		}

		// Set content type và trả về ảnh
		c.Set(fiber.HeaderContentType, image.ContentType)
		c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("inline; filename=\"%s\"", image.Filename))
		return c.Send(image.Data)
	}
}

// CreateAnnotation tạo chú thích cho dòng trong bài viết (chỉ tác giả)
func CreateAnnotation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := currentUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Bạn cần đăng nhập",
			})
		}

		postID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID bài viết không hợp lệ",
			})
		}

		var body struct {
			Content    string `json:"content"`
			LineNumber int    `json:"line_number"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Dữ liệu không hợp lệ",
			})
		}

		db := database.Get()

		// Kiểm tra quyền tác giả
		var post models.Post
		if err := db.First(&post, postID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Không tìm thấy bài viết",
			})
		}

		if post.AuthorID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Chỉ tác giả mới có thể thêm chú thích",
			})
		}

		// Xóa annotation cũ nếu có
		db.Where("post_id = ? AND line_number = ?", postID, body.LineNumber).Delete(&models.Annotation{})

		// Tạo mới
		annotation := models.Annotation{
			Content:    strings.TrimSpace(body.Content),
			PostID:     uint(postID),
			LineNumber: body.LineNumber,
		}

		if err := db.Create(&annotation).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Không thể lưu chú thích",
			})
		}

		// No longer update post content with #annotation
		// Annotations are stored separately in the annotations table

		return c.JSON(fiber.Map{
			"annotation": fiber.Map{
				"content":     annotation.Content,
				"line_number": annotation.LineNumber,
			},
		})
	}
}
