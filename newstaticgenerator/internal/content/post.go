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

	"github.com/alecthomas/chroma"
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

// Nord-Light custom style definition
var nordLightEntries = chroma.StyleEntries{
	// Background - Snow Storm (light colors)
	chroma.Background: "bg:#eceff4", // nord6 - lightest
	chroma.PreWrapper: "bg:#eceff4", // nord6

	// Text - Polar Night (dark colors for contrast)
	chroma.Text: "#2e3440", // nord0 - darkest

	// Comments - muted Polar Night
	chroma.Comment:            "#4c566a italic", // nord3
	chroma.CommentSingle:      "#4c566a italic", // nord3
	chroma.CommentMultiline:   "#4c566a italic", // nord3
	chroma.CommentSpecial:     "#434c5e italic", // nord2
	chroma.CommentHashbang:    "#4c566a italic", // nord3
	chroma.CommentPreproc:     "#5e81ac",        // nord10
	chroma.CommentPreprocFile: "#5e81ac",        // nord10

	// Keywords - Frost (darker blues for light background)
	chroma.Keyword:            "#5e81ac bold", // nord10 - darker blue
	chroma.KeywordConstant:    "#5e81ac bold", // nord10
	chroma.KeywordDeclaration: "#5e81ac bold", // nord10
	chroma.KeywordNamespace:   "#5e81ac bold", // nord10
	chroma.KeywordPseudo:      "#5e81ac bold", // nord10
	chroma.KeywordReserved:    "#5e81ac bold", // nord10
	chroma.KeywordType:        "#5e81ac bold", // nord10

	// Operators - same as keywords for consistency
	chroma.Operator:     "#5e81ac",      // nord10
	chroma.OperatorWord: "#5e81ac bold", // nord10

	// Strings - Aurora green (darkened for contrast on light bg)
	chroma.LiteralString:          "#2d5016",        // Darker version of nord14
	chroma.LiteralStringSingle:    "#2d5016",        // Darker green
	chroma.LiteralStringDouble:    "#2d5016",        // Darker green
	chroma.LiteralStringBacktick:  "#2d5016",        // Darker green
	chroma.LiteralStringChar:      "#2d5016",        // Darker green
	chroma.LiteralStringDoc:       "#4c566a italic", // nord3 for docstrings
	chroma.LiteralStringEscape:    "#8f4e00",        // Darker version of nord13
	chroma.LiteralStringHeredoc:   "#2d5016",        // Darker green
	chroma.LiteralStringInterpol:  "#2d5016",        // Darker green
	chroma.LiteralStringOther:     "#2d5016",        // Darker green
	chroma.LiteralStringRegex:     "#8f4e00",        // Darker version of nord13
	chroma.LiteralStringSymbol:    "#2d5016",        // Darker green
	chroma.LiteralStringAffix:     "#2d5016",        // Darker green
	chroma.LiteralStringDelimiter: "#2d5016",        // Darker green

	// Numbers - Aurora purple (darkened for contrast)
	chroma.LiteralNumber:            "#7c4d7a", // Darker version of nord15
	chroma.LiteralNumberBin:         "#7c4d7a", // Darker purple
	chroma.LiteralNumberFloat:       "#7c4d7a", // Darker purple
	chroma.LiteralNumberHex:         "#7c4d7a", // Darker purple
	chroma.LiteralNumberInteger:     "#7c4d7a", // Darker purple
	chroma.LiteralNumberIntegerLong: "#7c4d7a", // Darker purple
	chroma.LiteralNumberOct:         "#7c4d7a", // Darker purple

	// Names/Identifiers - use Frost and Polar Night colors
	chroma.Name:          "#2e3440", // nord0 - default text
	chroma.NameAttribute: "#5e81ac", // nord10 - darker blue
	chroma.NameBuiltin:   "#5e81ac", // nord10
	chroma.NameClass:     "#5e81ac", // nord10
	chroma.NameConstant:  "#5e81ac", // nord10
	chroma.NameDecorator: "#a0522d", // Darker version of nord12
	chroma.NameEntity:    "#a0522d", // Darker orange
	chroma.NameException: "#8b2635", // Darker version of nord11
	chroma.NameFunction:  "#5e81ac", // nord10
	chroma.NameLabel:     "#5e81ac", // nord10
	chroma.NameNamespace: "#5e81ac", // nord10
	chroma.NameTag:       "#5e81ac", // nord10

	// Errors - Aurora red (darkened)
	chroma.Error:        "#8b2635", // Darker version of nord11
	chroma.GenericError: "#8b2635", // Darker red

	// Other tokens
	chroma.Punctuation:       "#3b4252",      // nord1
	chroma.GenericDeleted:    "#8b2635",      // Darker red
	chroma.GenericInserted:   "#2d5016",      // Darker green
	chroma.GenericHeading:    "#5e81ac bold", // nord10
	chroma.GenericSubheading: "#5e81ac bold", // nord10
	chroma.GenericPrompt:     "#4c566a",      // nord3
	chroma.GenericTraceback:  "#8b2635",      // Darker red

	// Line numbers and highlighting
	chroma.LineNumbers:      "#4c566a",    // nord3
	chroma.LineNumbersTable: "#4c566a",    // nord3
	chroma.LineHighlight:    "bg:#e5e9f0", // nord5 - slightly darker than bg
}

var nordLight = chroma.MustNewStyle("nord-light", nordLightEntries)

func init() {
	// Register the nord-light style with Chroma
	styles.Register(nordLight)
}

// Available styles: abap, algol, algol_nu, arduino, autumn, base16-snazzy, borland, bw, colorful, doom-one, doom-one2, dracula, emacs, friendly, fruity, github, hr_high_contrast, hrdark, igor, lovelace, manni, monokai, monokailight, murphy, native, nord, onesenterprise, paraiso-dark, paraiso-light, pastie, perldoc, pygments, rainbow_dash, rrt, solarized-dark, solarized-dark256, solarized-light, swapoff, tango, trac, vim, vs, vulcan, witchhazel, xcode, xcode-dark
// Chroma theme configuration for light and dark modes
const chromaLightTheme = "nord-light" // Custom Nord light theme
const chromaDarkTheme = "nord"        // Good dark themes: nord, solarized-dark, doom-one2

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
			// When using WithClasses(true), the style doesn't affect the HTML output
			// The CSS class names are standardized, and colors come from the CSS file
			style := styles.Fallback // Any style works for generating class names

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
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
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

// GenerateChromaCSS creates a CSS file with chroma syntax highlighting styles for both light and dark themes
func GenerateChromaCSS(outputDir string) error {
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

	// Generate light theme CSS
	lightStyle := styles.Get(chromaLightTheme)
	if lightStyle == nil {
		lightStyle = styles.Fallback
	}

	// Write light theme CSS (default)
	if err := formatter.WriteCSS(cssFile, lightStyle); err != nil {
		return fmt.Errorf("failed to write light chroma CSS: %w", err)
	}

	// Generate dark theme CSS wrapped in media query
	darkStyle := styles.Get(chromaDarkTheme)
	if darkStyle == nil {
		darkStyle = styles.Fallback
	}

	// Write dark theme CSS inside media query
	cssFile.WriteString("\n\n/* Dark theme syntax highlighting */\n")
	cssFile.WriteString("@media (prefers-color-scheme: dark) {\n")

	// Create a temporary buffer to capture the dark theme CSS
	var darkBuffer strings.Builder
	if err := formatter.WriteCSS(&darkBuffer, darkStyle); err != nil {
		return fmt.Errorf("failed to generate dark chroma CSS: %w", err)
	}

	// Indent the dark theme CSS and write it
	darkCSS := darkBuffer.String()
	lines := strings.Split(darkCSS, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			cssFile.WriteString("  " + line + "\n")
		}
	}

	cssFile.WriteString("}\n")

	fmt.Printf("✨ Generated chroma CSS with light/dark themes at %s\n", cssPath)
	return nil
}
