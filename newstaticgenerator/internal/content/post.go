package content

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"

	"github.com/javierchavarri/goranite/internal/config"
)

type PostMatter struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Tags    []string  `json:"tags"`
	Summary string    `json:"summary"`
}

// Chroma theme configuration - change this to switch syntax highlighting theme
// Available styles: abap, algol, algol_nu, arduino, autumn, base16-snazzy, borland, bw, colorful, doom-one, doom-one2, dracula, emacs, friendly, fruity, github, hr_high_contrast, hrdark, igor, lovelace, manni, monokai, monokailight, murphy, native, nord, onesenterprise, paraiso-dark, paraiso-light, pastie, perldoc, pygments, rainbow_dash, rrt, solarized-dark, solarized-dark256, solarized-light, swapoff, tango, trac, vim, vs, vulcan, witchhazel, xcode, xcode-dark
// Good options for light backgrounds: github, colorful, friendly, vs, xcode
const chromaTheme = "xcode"

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

// Custom renderer for syntax highlighting
type syntaxHighlightRenderer struct{}

func newSyntaxHighlightRenderer() renderer.NodeRenderer {
	return &syntaxHighlightRenderer{}
}

func (r *syntaxHighlightRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// Register our custom handler for fenced code blocks
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
}

func (r *syntaxHighlightRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if !entering {
		return ast.WalkContinue, nil
	}

	// Get the language from the code block
	language := n.Language(source)

	// Extract the code content
	var codeBuffer bytes.Buffer
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		codeBuffer.Write(line.Value(source))
	}
	code := codeBuffer.String()

	// Try to highlight with chroma using CSS classes
	if len(language) > 0 {
		lexer := lexers.Get(string(language))
		if lexer != nil {
			style := styles.Get(chromaTheme)
			if style == nil {
				style = styles.Fallback
			}

			formatter := html.New(
				html.WithClasses(true),      // Use CSS classes instead of inline styles
				html.WithLineNumbers(false), // No line numbers for simplicity
			)

			iterator, err := lexer.Tokenise(nil, code)
			if err == nil {
				var highlightedBuffer bytes.Buffer
				err = formatter.Format(&highlightedBuffer, style, iterator)
				if err == nil {
					result := highlightedBuffer.String()
					w.WriteString(result)
					return ast.WalkContinue, nil
				}
			}
		}
	}

	// Fallback: render as plain code block
	w.WriteString("<pre><code")
	if len(language) > 0 {
		w.WriteString(" class=\"language-")
		w.Write(language)
		w.WriteString("\"")
	}
	w.WriteString(">")

	// Write the code content (HTML-escaped)
	template.HTMLEscape(w, []byte(code))

	w.WriteString("</code></pre>\n")
	return ast.WalkContinue, nil
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
		jsonBytes, err := json.Marshal(frontmatterData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
		}

		if err := json.Unmarshal(jsonBytes, &matter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal frontmatter: %w", err)
		}
	}

	// Convert markdown to HTML with our custom syntax highlighting
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
		),
	)

	// Add our custom syntax highlighting renderer to the existing renderer
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newSyntaxHighlightRenderer(), 200),
		),
	)

	var buf strings.Builder
	if err := md.Convert(body, &buf); err != nil {
		return nil, err
	}

	htmlContent := buf.String()

	// Generate slug from filename
	filename := filepath.Base(path)
	slug := strings.TrimSuffix(filename, ".md")

	// Calculate reading time
	wordCount := len(strings.Fields(string(body)))
	readingTime := (wordCount + 199) / 200
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

// GenerateChromaCSS creates a CSS file with chroma syntax highlighting styles
func GenerateChromaCSS(outputDir string) error {
	style := styles.Get(chromaTheme)
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true))

	// Create CSS file
	cssPath := filepath.Join(outputDir, "css", "chroma.css")
	if err := os.MkdirAll(filepath.Dir(cssPath), 0755); err != nil {
		return fmt.Errorf("failed to create CSS directory: %w", err)
	}

	cssFile, err := os.Create(cssPath)
	if err != nil {
		return fmt.Errorf("failed to create chroma CSS file: %w", err)
	}
	defer cssFile.Close()

	if err := formatter.WriteCSS(cssFile, style); err != nil {
		return fmt.Errorf("failed to write chroma CSS: %w", err)
	}

	fmt.Printf("✨ Generated chroma CSS at %s\n", cssPath)
	return nil
}
