# Goranite Migration Plan
*From Gatsby to Go Static Site Generator*

## Project Overview

Migrating javierchavarri.com from Gatsby to a custom Go static site generator
called **Goranite** (Go + Granite - solid as a rock!). This project serves dual
purposes: learning Go and creating a minimal, performant blog platform.

## Design Direction

- **Styling**: Custom professional theme with Nord color scheme - clean,
  readable aesthetic perfect for technical content
- **Content**: Markdown + frontmatter (same as current setup)
- **Deployment**: Cloudflare Pages (considering migration from Netlify)
- **Focus**: OCaml/ReasonML syntax highlighting support

## Project Structure

```
/
├── newsite/                    # Blog content and configuration
│   ├── content/
│   │   └── posts/             # Markdown blog posts
│   ├── static/
│   │   ├── css/
│   │   │   └── custom.css     # Custom professional theme
│   │   └── images/            # Static images
│   ├── config.json            # Site configuration
│   └── public/                # Generated output (git-ignored)
└── newstaticgenerator/         # Go static site generator
    ├── main.go                # CLI entry point
    ├── go.mod                 # Go module definition
    ├── internal/
    │   ├── config/
    │   │   └── config.go      # Site configuration handling
    │   ├── content/
    │   │   └── post.go        # Post struct and parsing + Nord themes
    │   ├── generator/
    │   │   └── generator.go   # Main site generation logic
    │   └── templates/
    │       └── templates.go   # Template handling
    └── templates/
        ├── index.html         # Homepage template
        └── post.html          # Individual post template
```

## Go Learning Opportunities

This project covers essential Go concepts:

1. **Modules and Project Structure** - Go's package system
2. **File I/O and Path Manipulation** - Reading markdown files, walking
   directories  
3. **Text Processing** - Parsing frontmatter, markdown processing
4. **Templates** - Go's `html/template` package
5. **CLI Development** - Using `flag` package for commands
6. **Error Handling** - Go's explicit error handling patterns
7. **Structs and Methods** - Modeling blog posts, pages, site config
8. **Testing** - Writing tests for generator functions
9. **Concurrency** (optional) - Parallelizing file processing
10. **Custom Chroma Styles** - Implementing syntax highlighting themes

## Implementation Phases

### Phase 1: Foundation ✅
- [x] Basic Go project setup with modules (Go 1.24.3)
- [x] CLI structure with commands: `build`, `serve`, `new`
- [x] Dependencies: goldmark, chroma (minimal set)
- [x] Template structure created
- [x] Site configuration setup (JSON format)

### Phase 2: Core Functionality ✅
- [x] Custom frontmatter parsing (JSON format, no external deps)
- [x] Markdown processing with goldmark
- [x] Template system using Go's `html/template`
- [x] Static file handling (CSS, images)
- [x] Syntax highlighting with `chroma`
- [x] Basic site generation working

### Phase 3: Content Migration ✅
- [x] Extract content from current Gatsby site
- [x] Convert/adapt markdown files and assets
- [x] Implement OCaml/ReasonML syntax highlighting
- [x] Create professional responsive design with custom theme
- [x] Nord color scheme implementation (light + dark themes)

### Phase 4: Advanced Features ✅
- [x] Development server with live reload
- [x] Custom Nord-light theme implementation
- [x] Dual theme support (light/dark mode)
- [x] Fixed syntax highlighting edge cases
- [ ] RSS feed generation
- [ ] Sitemap generation
- [ ] Build optimization

### Phase 5: Deployment
- [ ] Set up Cloudflare Pages
- [ ] Configure automated build process
- [ ] Domain migration strategy
- [ ] Performance testing

## CLI Commands

```bash
# Build the static site
./goranite -build

# Start development server
./goranite -serve

# Create new post
./goranite -new "Post Title"
```

## Deployment Strategy

**Cloudflare Pages** (Recommended):
- ✅ Excellent global CDN performance
- ✅ Generous free tier
- ✅ Native Go build support
- ✅ Custom domains included
- ✅ Built-in analytics

**Domain Migration**:
- Option 1: Transfer domain from Namecheap to Cloudflare (simpler)
- Option 2: Keep domain at Namecheap, point nameservers to Cloudflare

## Theme Implementation

**Decision**: Custom professional theme with Nord colors
- ✅ Complete control over design and responsiveness
- ✅ Professional, readable aesthetic
- ✅ Nord color scheme for consistency
- ✅ Dual theme support (light/dark mode)
- ✅ Excellent syntax highlighting

**Implementation**:
- Custom CSS in `newsite/static/css/custom.css`
- Nord-light custom Chroma style for light mode
- Nord official Chroma style for dark mode
- CSS media queries for automatic theme switching
- Optimized for technical content readability

## Key Features

- **Markdown ❤️ Professional Design**: Perfect match for technical blog
- **Nord Color Scheme**: Consistent, professional aesthetic
- **Dual Theme Support**: Automatic light/dark mode switching
- **Syntax Highlighting**: Built-in support for OCaml/ReasonML with proper operator highlighting
- **Minimal Dependencies**: Pure Go, no Node.js build chain
- **Fast Builds**: Go's performance for quick iteration
- **Clean Output**: Semantic HTML with professional styling

## Benefits

1. **Learning Go**: Practical project with real-world application
2. **Performance**: Fast generation, minimal output
3. **Control**: Full ownership of build process and output
4. **Simplicity**: No complex toolchains or dependencies
5. **Portability**: Deploy anywhere Go runs
6. **Professional Design**: Clean, readable theme perfect for technical content

## Current Status

**Completed:**
- ✅ Project structure setup
- ✅ Go module initialization with latest Go (1.24.3)
- ✅ Minimal dependencies (only goldmark + chroma)
- ✅ Basic CLI framework
- ✅ HTML templates (self-contained)
- ✅ Professional CSS theme with Nord colors
- ✅ JSON configuration (no YAML dependencies)
- ✅ Custom frontmatter parser (stdlib only)
- ✅ **Working site generation!** 🎉
- ✅ **Content migration completed!** 🎉
- ✅ **Custom Nord-light theme implemented!** 🎉
- ✅ **Dual theme support working!** 🎉
- ✅ **Development server working!** 🎉

**Next Steps:**
- 🚧 RSS feed generation
- 🚧 Sitemap generation
- 🚧 Deployment to Cloudflare Pages
- 🚧 Performance optimization

---

*Goranite: Solid as granite, powered by Go* 🪨
