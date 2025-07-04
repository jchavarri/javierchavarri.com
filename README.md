# Javier Chávarri's Blog

This directory contains the content and configuration for
[javierchavarri.com](https://javierchavarri.com/) built with
[Goranite](https://github.com/jchavarri/goranite), a static site generator
written in Go.

## Prerequisites

- [Go](https://golang.org/) 1.21 or later
- [Goranite](https://github.com/jchavarri/goranite) static site generator

## Setup

1. **Install Goranite:**
   ```bash
   git clone https://github.com/jchavarri/goranite
   cd goranite
   ```

2. **Build the site:**
   ```bash
   go run main.go -build -site ../javierchavarri.com
   ```

3. **Serve locally with live reload:**
   ```bash
   go run main.go -serve -site ../javierchavarri.com
   ```
   
   Then open http://localhost:8080

## Directory Structure

```
javierchavarri.com/
├── config.json          # Site configuration (title, description, etc.)
├── content/
│   └── posts/           # Blog posts in Markdown with JSON frontmatter
├── static/
│   ├── css/            # Custom CSS files
│   └── images/         # Static images and assets
└── public/             # Generated output (created after build)
```

## Writing Posts

Posts are written in Markdown with JSON frontmatter. Create new files in
`content/posts/`:

```markdown
---
{
  "title": "Your Post Title",
  "date": "2024-01-15T00:00:00Z",
  "tags": ["tag1", "tag2", "programming"],
  "summary": "A brief summary of your post content"
}
---

# Your Post Title

Write your content here in Markdown. You can use:

- **Bold** and *italic* text
- Code blocks with syntax highlighting
- Links, images, tables, etc.

## Code Example

\```javascript
function hello() {
    console.log("Hello, world!");
}
\```
```

## Configuration

Edit `config.json` to customize:

- Site title and description
- Author information
- Social media links
- Build settings

## Deployment

The generated `public/` directory contains the complete static site ready for
deployment to any web server or CDN.

## Custom Styling

- Main styles: `static/css/custom.css`
- Syntax highlighting: Auto-generated by Goranite using Chroma
- The site supports both light and dark themes 