package handlers

import (
	"errors"
	"fmt"
	"net/mail"
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

func render(c *fiber.Ctx, view string, data fiber.Map) error {
	base := fiber.Map{
		"AppName":         "Cộng đồng Học DevOps",
		"Year":            time.Now().Year(),
		"IsAuthenticated": false,
		"RequestPath":     c.OriginalURL(),
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

	return c.Render(view, base)
}

func Home() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		var posts []models.Post
		latestPosts := []fiber.Map{}
		if err := db.Preload("Author").Order("created_at DESC").Limit(3).Find(&posts).Error; err == nil {
			for _, p := range posts {
				latestPosts = append(latestPosts, fiber.Map{
					"ID":           p.ID,
					"Title":        p.Title,
					"Summary":      p.Summary,
					"AuthorName":   p.Author.Name,
					"CreatedLabel": p.CreatedAt.Format("02/01/2006 15:04"),
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
			"LatestPosts": latestPosts,
		})
	}
}

type createPostRequest struct {
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Content  string `json:"content"`
	AuthorID uint   `json:"author_id"`
}

type createCommentRequest struct {
	Content  string `json:"content"`
	AuthorID uint   `json:"author_id"`
}

func PostsPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.Get()
		var posts []models.Post
		if err := db.Preload("Author").Order("created_at DESC").Find(&posts).Error; err != nil {
			return respondError(c, fiber.StatusInternalServerError, "Không thể tải danh sách bài viết", "/")
		}
		items := make([]fiber.Map, 0, len(posts))
		for _, p := range posts {
			items = append(items, fiber.Map{
				"ID":           p.ID,
				"Title":        p.Title,
				"Summary":      p.Summary,
				"Content":      p.Content,
				"AuthorName":   p.Author.Name,
				"CreatedLabel": p.CreatedAt.Format("02/01/2006 15:04"),
			})
		}
		return render(c, "pages/posts", fiber.Map{
			"Title": "Bài viết cộng đồng",
			"Posts": items,
		})
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

		comments := make([]fiber.Map, 0, len(post.Comments))
		for _, cm := range post.Comments {
			comments = append(comments, fiber.Map{
				"ID":         cm.ID,
				"Content":    cm.Content,
				"AuthorName": cm.Author.Name,
				"CreatedAt":  cm.CreatedAt.Format("02/01/2006 15:04"),
			})
		}

		return render(c, "pages/post_detail", fiber.Map{
			"Title": "Chi tiết bài viết",
			"Post": fiber.Map{
				"ID":           post.ID,
				"Title":        post.Title,
				"Summary":      post.Summary,
				"Content":      post.Content,
				"AuthorName":   post.Author.Name,
				"CreatedLabel": post.CreatedAt.Format("02/01/2006 15:04"),
			},
			"Comments": comments,
		})
	}
}

func RegisterPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/auth_register", fiber.Map{
			"Title": "Đăng ký tài khoản",
		})
	}
}

func LoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/auth_login", fiber.Map{
			"Title": "Đăng nhập",
			"Next":  c.Query("next"),
		})
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
		}

		body.Title = strings.TrimSpace(body.Title)
		body.Summary = strings.TrimSpace(body.Summary)
		body.Content = strings.TrimSpace(body.Content)

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
			AuthorID: body.AuthorID,
		}

		if err := db.Create(&post).Error; err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể tạo bài viết")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể tạo bài viết", "/posts")
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
			Content:  body.Content,
			PostID:   uint(postID),
			AuthorID: body.AuthorID,
		}

		if err := db.Create(&comment).Error; err != nil {
			if isJSON {
				return fiber.NewError(fiber.StatusInternalServerError, "không thể tạo bình luận")
			}
			return respondError(c, fiber.StatusInternalServerError, "Không thể tạo bình luận", fmt.Sprintf("/posts/%d", postID))
		}

		if isJSON {
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"message": "Thêm bình luận thành công",
				"comment": fiber.Map{
					"id":         comment.ID,
					"content":    comment.Content,
					"post_id":    comment.PostID,
					"author_id":  comment.AuthorID,
					"created_at": comment.CreatedAt,
				},
			})
		}

		setFlash(c, "success", "Bình luận đã được gửi")
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
		if redirect == "" {
			redirect = "/posts"
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
		if redirect == "" {
			redirect = "/posts"
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
		})
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
		})
	}
}

func About() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, "pages/about", fiber.Map{
			"Title":       "Về dự án",
			"Description": "Sứ mệnh của cộng đồng học DevOps và cách chúng tôi vận hành.",
			"Sections": []fiber.Map{
				{
					"Heading": "Tầm nhìn",
					"Content": "Tạo ra nền tảng học DevOps thực tế, cập nhật liên tục, giúp đội ngũ kỹ thuật vận hành hệ thống ổn định.",
				},
				{
					"Heading": "Giá trị cốt lõi",
					"Content": "Chia sẻ kinh nghiệm thực chiến, học hỏi không ngừng và vận hành minh bạch.",
				},
				{
					"Heading": "Bạn có thể làm gì?",
					"Content": "Viết nội dung, chia sẻ scripts, template infrastructure, review tài liệu và tổ chức workshop.",
				},
			},
		})
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
			"ContactEmail": "hello@hocdevops.community",
		})
	}
}
