package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"

	"fiber-learning-community/internal/database"
	"fiber-learning-community/internal/handlers"
)

func main() {
	database.Init()

	engine := html.New("./views", ".html")

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
