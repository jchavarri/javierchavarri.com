# Contributing to Goranite

## Architecture Overview

Goranite is structured as a simple CLI tool with these main components:

```
newstaticgenerator/
├── main.go                 # CLI entry point
├── internal/
│   ├── config/            # Site configuration and data structures
│   ├── content/           # Post parsing, loading, and custom Chroma themes
│   └── generator/         # Site generation logic
└── templates/
    ├── index.html         # Homepage template
    └── post.html          # Individual post template
```

## Template System

Goranite uses Go's `html/template` package with **self-contained templates** for simplicity.

### Template Data Structure

Templates receive a `TemplateData` struct with:
- **Site-wide data** (from config.json): `.SiteTitle`, `.Author`, etc.
- **Content data**: `.Posts` (for index) or `.Post` (for individual posts)

### Template Files

- **`index.html`** - Homepage (receives `.Posts` slice)
- **`post.html`** - Individual post (receives `.Post` object)

### Template Syntax

```html
<!-- Site info -->
<title>{{.SiteTitle}}</title>

<!-- Loop posts (index.html) -->
{{range .Posts}}
  <h2><a href="{{.URL}}">{{.Title}}</a></h2>
{{end}}

<!-- Single post (post.html) -->
<h1>{{.Post.Title}}</h1>
<div>{{.Post.Content}}</div>
```

## Syntax Highlighting

Goranite includes custom Chroma themes for consistent Nord-based styling:

### Custom Nord-Light Theme
- Implemented in `content/post.go` as `nordLightEntries`
- Uses Nord color palette with proper contrast ratios
- Covers all syntax tokens including operators and keywords

### Theme Registration
```go
func init() {
    nordLightStyle := chroma.MustNewStyle("nord-light", nordLightEntries)
    styles.Register(nordLightStyle)
}
```

### Adding New Themes
1. Define color entries using `chroma.StyleEntries`
2. Create style with `chroma.MustNewStyle()`
3. Register with `styles.Register()` in `init()` function
4. Update theme constants in `post.go`

## Development Workflow

1. **Setup:**
   ```bash
   cd newstaticgenerator
   go mod tidy
   ```

2. **Test changes:**
   ```bash
   go run main.go -build    # Generate site
   go run main.go -serve    # Serve at localhost:8080
   ```

3. **Test with different content** in `../newsite/`

## Adding Features

### New CLI Command
1. Add flag to `main.go`
2. Add case to switch statement  
3. Implement function

### New Template Data
1. Add field to `TemplateData` struct in `config/config.go`
2. Use in templates with `{{.NewField}}`

### New Post Metadata
1. Add field to `PostMatter` struct in `content/post.go`
2. Add JSON tag: `NewField string \`json:"new_field"\``

### New Syntax Theme
1. Define `chroma.StyleEntries` with color mappings
2. Register in `init()` function
3. Update theme constants
4. Test with various code samples

## Dependencies Philosophy

We keep dependencies minimal:
- **`github.com/yuin/goldmark`** - Markdown processing
- **`github.com/alecthomas/chroma`** - Syntax highlighting  
- **Go stdlib** - Everything else

No YAML libraries, no web frameworks, no complex build tools.

## Code Style

- Use `gofmt` for formatting
- Follow Go naming conventions
- Handle errors explicitly
- Keep functions focused

## Testing

Currently manual testing. Run `go run main.go -build && go run main.go -serve` to test changes.

For syntax highlighting, test with the debug post at `/debug-syntax` which includes comprehensive examples.
