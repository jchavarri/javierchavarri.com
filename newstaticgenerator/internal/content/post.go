package content

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"html/template"

	"github.com/javierchavarri/goranite/internal/config"
)

type PostMatter struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Tags    []string  `json:"tags"`
	Summary string    `json:"summary"`
}

func parseFrontmatter(content []byte) (map[string]interface{}, []byte, error) {
	if !bytes.HasPrefix(content, []byte("---\n")) {
		return nil, content, nil // No frontmatter
	}

	// Find end of frontmatter
	end := bytes.Index(content[4:], []byte("\n---\n"))
	if end == -1 {
		return nil, nil, errors.New("unclosed frontmatter")
	}

	jsonContent := content[4 : end+4]
	bodyContent := content[end+8:]

	var frontmatter map[string]interface{}
	if err := json.Unmarshal(jsonContent, &frontmatter); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON frontmatter: %w", err)
	}

	return frontmatter, bodyContent, nil
}

func LoadPosts(contentDir string) ([]config.Post, error) {
	var posts []config.Post

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		post, err := LoadPost(path)
		if err != nil {
			return fmt.Errorf("failed to load post %s: %w", path, err)
		}

		posts = append(posts, *post)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort posts by date (newest first)
	for i := 0; i < len(posts)-1; i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].Date.Before(posts[j].Date) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}

	return posts, nil
}

func LoadPost(path string) (*config.Post, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse frontmatter
	frontmatterData, body, err := parseFrontmatter(content)
	if err != nil {
		return nil, err
	}

	// Convert map to our struct
	var matter PostMatter
	if frontmatterData != nil {
		// Convert map[string]interface{} to JSON then to struct
		jsonBytes, err := json.Marshal(frontmatterData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
		}

		if err := json.Unmarshal(jsonBytes, &matter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal frontmatter: %w", err)
		}
	}

	// Convert markdown to HTML
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf strings.Builder
	if err := md.Convert(body, &buf); err != nil {
		return nil, err
	}

	htmlContent := buf.String()

	// Add syntax highlighting
	htmlContent = addSyntaxHighlighting(htmlContent)

	// Generate slug from filename
	filename := filepath.Base(path)
	slug := strings.TrimSuffix(filename, ".md")

	// Calculate reading time (rough estimate: 200 words per minute)
	wordCount := len(strings.Fields(string(body)))
	readingTime := (wordCount + 199) / 200 // Round up
	if readingTime < 1 {
		readingTime = 1
	}

	return &config.Post{
		Title:       matter.Title,
		Date:        matter.Date,
		Tags:        matter.Tags,
		Summary:     matter.Summary,
		Content:     template.HTML(htmlContent),
		URL:         "/" + slug + "/",
		ReadingTime: readingTime,
		Slug:        slug,
	}, nil
}

func addSyntaxHighlighting(html string) string {
	// Simple regex to find code blocks
	codeBlockRegex := regexp.MustCompile(`<pre><code class="language-(\w+)">(.*?)</code></pre>`)

	return codeBlockRegex.ReplaceAllStringFunc(html, func(match string) string {
		submatch := codeBlockRegex.FindStringSubmatch(match)
		if len(submatch) < 3 {
			return match
		}

		language := submatch[1]
		code := submatch[2]

		// Use chroma for syntax highlighting
		var buf strings.Builder
		err := quick.Highlight(&buf, code, language, "html", "terminal")
		if err != nil {
			return match // Return original if highlighting fails
		}

		return buf.String()
	})
}
