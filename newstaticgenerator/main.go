package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var (
		buildCmd = flag.Bool("build", false, "Build the static site")
		serveCmd = flag.Bool("serve", false, "Start development server")
		newCmd   = flag.String("new", "", "Create new post with given title")
	)
	flag.Parse()

	switch {
	case *buildCmd:
		fmt.Println("🔨 Building static site...")
		if err := buildSite(); err != nil {
			log.Fatalf("Build failed: %v", err)
		}
		fmt.Println("✅ Site built successfully!")

	case *serveCmd:
		fmt.Println("🚀 Starting development server...")
		if err := serveSite(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}

	case *newCmd != "":
		fmt.Printf("📝 Creating new post: %s\n", *newCmd)
		if err := createNewPost(*newCmd); err != nil {
			log.Fatalf("Failed to create post: %v", err)
		}
		fmt.Println("✅ Post created successfully!")

	default:
		fmt.Println("🪨 Goranite - Static Site Generator")
		fmt.Println("Solid as granite, powered by Go")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  -build       Build the static site")
		fmt.Println("  -serve       Start development server")
		fmt.Println("  -new 'Title' Create new post")
		flag.PrintDefaults()
	}
}

func buildSite() error {
	// TODO: Implement site building
	fmt.Println("Building site... (not implemented yet)")
	return nil
}

func serveSite() error {
	// TODO: Implement development server
	fmt.Println("Starting server... (not implemented yet)")
	return nil
}

func createNewPost(title string) error {
	// TODO: Implement post creation
	fmt.Printf("Creating post '%s'... (not implemented yet)\n", title)
	return nil
}
