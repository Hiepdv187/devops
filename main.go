package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	"fiber-learning-community/internal/database"
	"fiber-learning-community/internal/handlers"
)

func main() {
	database.Init()

	engine := html.New("./views", ".html")
	engine.AddFunc("now", func() time.Time {
		return time.Now()
	})
	engine.AddFunc("date", func(layout string, t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format(layout)
	})
	engine.AddFunc("markdown", func(text string) template.HTML {
		if strings.TrimSpace(text) == "" {
			return template.HTML("")
		}
		rendered := blackfriday.Run([]byte(text), blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs|blackfriday.HardLineBreak))
		sanitized := bluemonday.UGCPolicy().SanitizeBytes(rendered)
		return template.HTML(sanitized)
	})
	engine.AddFunc("json", func(v interface{}) (template.JS, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return template.JS(b), nil
	})

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	app.Static("/static", "./public")

	app.Get("/", handlers.Home())
	app.Get("/courses", handlers.Courses())
	app.Get("/contributors", handlers.Contributors())
	app.Get("/about", handlers.About())
	app.Get("/contribute", handlers.Contribute())
	app.Get("/posts", handlers.PostsPage())
	app.Get("/posts/:id", handlers.PostDetailPage())
	app.Get("/auth/register", handlers.RegisterPage())
	app.Get("/auth/login", handlers.LoginPage())
	app.Post("/auth/register", handlers.Register())
	app.Post("/auth/login", handlers.Login())
	app.Post("/auth/logout", handlers.Logout())
	app.Post("/posts", handlers.CreatePost())
	app.Post("/posts/:id/comments", handlers.CreateComment())
	app.Post("/posts/:id/annotations", handlers.CreateAnnotation())
	app.Post("/posts/:id/edit", handlers.UpdatePost())
	app.Post("/upload/image", handlers.UploadImage())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
