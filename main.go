package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	xhtml "golang.org/x/net/html"

	"fiber-learning-community/internal/database"
	"fiber-learning-community/internal/handlers"
)

var (
	markdownPolicy *bluemonday.Policy
)

func Markdown(input string) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.HardLineBreak
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(input))

	renderer := mdhtml.NewRenderer(mdhtml.RendererOptions{Flags: mdhtml.CommonFlags})
	rendered := markdown.Render(doc, renderer)

	flattened := flattenLists(rendered)
	safeHTML := markdownPolicy.SanitizeBytes(flattened)

	return template.HTML(safeHTML)
}

func flattenLists(input []byte) []byte {
	root, err := xhtml.Parse(bytes.NewReader(input))
	if err != nil {
		return input
	}

	body := findNode(root, "body")
	target := root
	if body != nil {
		target = body
	}

	flattenListNodes(target)

	var buf bytes.Buffer
	for child := target.FirstChild; child != nil; child = child.NextSibling {
		if err := xhtml.Render(&buf, child); err != nil {
			return input
		}
	}
	return buf.Bytes()
}

func findNode(node *xhtml.Node, name string) *xhtml.Node {
	if node.Type == xhtml.ElementNode && node.Data == name {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findNode(child, name); found != nil {
			return found
		}
	}
	return nil
}

func flattenListNodes(node *xhtml.Node) {
	for child := node.FirstChild; child != nil; {
		next := child.NextSibling
		if child.Type == xhtml.ElementNode && (child.Data == "ol" || child.Data == "ul") {
			fragments := transformListNode(child)
			for _, fragment := range fragments {
				node.InsertBefore(fragment, child)
			}
			node.RemoveChild(child)
		} else {
			flattenListNodes(child)
		}
		child = next
	}
}

func transformListNode(list *xhtml.Node) []*xhtml.Node {
	ordered := list.Data == "ol"
	counter := 1
	if ordered {
		for _, attr := range list.Attr {
			if attr.Key == "start" {
				if value, err := strconv.Atoi(attr.Val); err == nil {
					counter = value
				}
				break
			}
		}
	}

	var fragments []*xhtml.Node
	for item := list.FirstChild; item != nil; item = item.NextSibling {
		if item.Type != xhtml.ElementNode || item.Data != "li" {
			continue
		}

		// Giữ danh sách con, đừng flatten sớm
		hasNestedList := false
		for c := item.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == xhtml.ElementNode && (c.Data == "ul" || c.Data == "ol") {
				hasNestedList = true
				break
			}
		}

		prefix := "- "
		if ordered {
			prefix = strconv.Itoa(counter) + ". "
		}

		paragraph := &xhtml.Node{Type: xhtml.ElementNode, Data: "p"}
		paragraph.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: prefix})
		appendListItemContent(paragraph, item)
		fragments = append(fragments, paragraph)

		if hasNestedList {
			for c := item.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == xhtml.ElementNode && (c.Data == "ul" || c.Data == "ol") {
					nested := transformListNode(c)
					fragments = append(fragments, nested...)
				}
			}
		}

		if ordered {
			counter++
		}
	}
	return fragments
}

func appendListItemContent(dest, src *xhtml.Node) {
	for child := src.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == xhtml.ElementNode && child.Data == "p" {
			appendListItemContent(dest, child)
			if child.NextSibling != nil {
				dest.AppendChild(&xhtml.Node{Type: xhtml.ElementNode, Data: "br"})
			}
			continue
		}
		dest.AppendChild(cloneNode(child))
	}
}

func cloneNode(n *xhtml.Node) *xhtml.Node {
	clone := &xhtml.Node{
		Type:      n.Type,
		Data:      n.Data,
		Namespace: n.Namespace,
		Attr:      append([]xhtml.Attribute(nil), n.Attr...),
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		clone.AppendChild(cloneNode(child))
	}
	return clone
}

func init() {
	markdownPolicy = bluemonday.UGCPolicy()
	markdownPolicy.AllowRelativeURLs(true)
	markdownPolicy.AllowURLSchemes("http", "https", "data")
	markdownPolicy.AllowImages()
	markdownPolicy.AllowElements("figure", "figcaption")
	markdownPolicy.AllowAttrs("class").OnElements("figure", "figcaption")
	markdownPolicy.AllowAttrs("src", "alt", "title", "loading", "width", "height", "class").OnElements("img")
}

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
	engine.AddFunc("markdown", Markdown)
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
