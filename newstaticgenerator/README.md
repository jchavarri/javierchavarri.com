# 🪨 Goranite

A minimal, fast static site generator built with Go. Solid as granite, powered by Go.

## Features

- ✅ **Minimal Dependencies** - Only goldmark + chroma
- ✅ **JSON Frontmatter** - No YAML dependencies
- ✅ **Syntax Highlighting** - Built-in support for 100+ languages
- ✅ **Terminal CSS** - Clean, developer-focused styling
- ✅ **Fast Builds** - Go's performance for quick iteration

## Quick Start

1. **Install Go** (1.21+)
2. **Clone and build:**
   ```bash
   git clone https://github.com/javierchavarri/goranite
   cd goranite/newstaticgenerator
   go build -o goranite
   ```

3. **Create your site:**
   ```bash
   mkdir mysite
   cd mysite
   
   # Create config
   cat > config.json << EOF
   {
     "site": {
       "title": "My Blog",
       "description": "My awesome blog",
       "author": "Your Name",
       "url": "https://myblog.com"
     },
     "build": {
       "output_dir": "public",
       "posts_per_page": 10
     }
   }
   EOF
   
   # Create directories
   mkdir -p content/posts static/css static/images
   
   # Download Terminal CSS
   curl -o static/css/terminal.css https://unpkg.com/terminal.css@0.7.4/dist/terminal.min.css
   
   # Create custom CSS
   touch static/css/custom.css
   
   # Create first post
   cat > content/posts/hello.md << EOF
   ---
   {
     "title": "Hello World",
     "date": "2024-01-15T00:00:00Z",
     "tags": ["hello", "first-post"],
     "summary": "My first post with Goranite"
   }
   ---
   
   # Hello World
   
   This is my first post with **Goranite**!
   
   ## Code Example
   
   \`\`\`go
   func main() {
       fmt.Println("Hello, Goranite!")
   }
   \`\`\`
   
   Pretty cool! 🪨
   EOF
   ```

4. **Build and serve:**
   ```bash
   ./goranite -build
   ./goranite -serve  # Visit http://localhost:8080
   ```

## Post Format

Posts use JSON frontmatter between `---` delimiters:

```markdown
---
{
  "title": "Post Title",
  "date": "2024-01-15T00:00:00Z",
  "tags": ["go", "blog"],
  "summary": "Brief description for the homepage"
}
---

# Your Content Here

Regular **markdown** content with syntax highlighting:

\`\`\`javascript
console.log("Hello, world!");
\`\`\`
```

## Commands

- `goranite -build` - Generate static site
- `goranite -serve` - Development server with auto-build
- `goranite -new "Title"` - Create new post (coming soon)
- `goranite -site /path` - Specify site directory (default: ../newsite)

## Configuration

Site configuration in `config.json`:

```json
{
  "site": {
    "title": "Your Blog Title",
    "description": "Blog description",
    "author": "Your Name",
    "url": "https://yourdomain.com"
  },
  "build": {
    "output_dir": "public",
    "posts_per_page": 10
  },
  "social": {
    "twitter": "@yourhandle",
    "github": "yourusername"
  }
}
```

## Directory Structure

```
mysite/
├── config.json           # Site configuration
├── content/
│   └── posts/            # Markdown posts
├── static/
│   ├── css/
│   │   ├── terminal.css  # Terminal CSS framework
│   │   └── custom.css    # Your custom styles
│   └── images/           # Static images
└── public/               # Generated site (git-ignore this)
```

## Why Goranite?

- **Learning Go** - Great project for Go beginners
- **No Node.js** - Pure Go, no JavaScript build chains
- **Minimal** - Only essential dependencies
- **Fast** - Go's performance for quick builds
- **Terminal Aesthetic** - Perfect for technical blogs

## Deployment

The generated `public/` directory contains static files that can be deployed to any static hosting service:

- **Netlify** - Drag and drop the `public/` folder
- **GitHub Pages** - Push to a repository with Pages enabled
- **Vercel** - Connect your repository for automatic deployments
- **Traditional hosting** - Upload via FTP/SFTP
- **CDN services** - Any service that serves static files

Most services support custom domains and HTTPS out of the box.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and architecture details.

## License

MIT License - see LICENSE file for details.
