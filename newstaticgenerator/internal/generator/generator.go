package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/javierchavarri/goranite/internal/config"
	"github.com/javierchavarri/goranite/internal/content"
)

type Generator struct {
	config    *config.Config
	templates *template.Template
}

func New(configPath, templatesDir string) (*Generator, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	tmpl, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Generator{
		config:    cfg,
		templates: tmpl,
	}, nil
}

func (g *Generator) Build(contentDir, staticDir, outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy static files
	if err := g.copyStaticFiles(staticDir, outputDir); err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}

	// Load posts
	posts, err := content.LoadPosts(filepath.Join(contentDir, "posts"))
	if err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}

	// Generate index page
	if err := g.generateIndex(posts, outputDir); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}

	// Generate individual post pages
	for _, post := range posts {
		if err := g.generatePost(&post, outputDir); err != nil {
			return fmt.Errorf("failed to generate post %s: %w", post.Title, err)
		}
	}

	return nil
}

func (g *Generator) generateIndex(posts []config.Post, outputDir string) error {
	data := config.NewTemplateData(g.config)
	data.Posts = posts

	file, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return g.templates.ExecuteTemplate(file, "index.html", data)
}

func (g *Generator) generatePost(post *config.Post, outputDir string) error {
	data := config.NewTemplateData(g.config)
	data.Post = post
	data.Title = post.Title
	data.Description = post.Summary

	// Create post directory
	postDir := filepath.Join(outputDir, post.Slug)
	if err := os.MkdirAll(postDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(postDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return g.templates.ExecuteTemplate(file, "post.html", data)
}

func (g *Generator) copyStaticFiles(staticDir, outputDir string) error {
	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, src, info.Mode())
	})
}
