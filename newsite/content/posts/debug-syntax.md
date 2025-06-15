---
{
  "title": "Debug Syntax Highlighting",
  "date": "2024-01-15T00:00:00Z",
  "summary": "Testing syntax highlighting"
}
---

# Debug Test

Here's some Go code:

```go
func main() {
    fmt.Println("Hello, World!")
}
```

And some JavaScript:

```javascript
console.log("Hello, World!");
```

Here's some OCaml:

```ocaml
(* Define a record type for a person *)
type person = {
  name: string;
  age: int;
}

(* Function to greet a person - takes a person record *)
let greet person =
  Printf.printf "Hello, %s! You are %d years old.\n" person.name person.age

(* Main execution block *)
let () =
  (* Create a person record *)
  let p = { name = "Alice"; age = 30 } in
  greet p;
  (* Use pipe operator to chain operations *)
  List.fold_left (+) 0 [1; 2; 3; 4; 5]  (* Sum a list *)
  |> Printf.printf "Sum: %d\n"          (* Print the result *)
```

And some Reason:

```reason
/* Define a record type for a person */
type person = {
  name: string,
  age: int,
};

/* Function to greet a person - arrow function syntax */
let greet = (person) =>
  Printf.printf("Hello, %s! You are %d years old.\n", person.name, person.age);

/* Main execution block */
let () = {
  /* Create a person record */
  let p = {name: "Alice", age: 30};
  greet(p);
  /* Use pipe operator to chain operations */
  [1, 2, 3, 4, 5]                       /* Create a list */
  |> List.fold_left((+), 0)             /* Sum the list */
  |> Printf.printf("Sum: %d\n");        /* Print the result */
};
```

Run the build and see what the debug output shows:

```bash
go run main.go -build
```

## 4. Alternative Approaches (After Debugging)

You're right that regex might not be the best approach. Here are better alternatives:

### Option A: Goldmark Extension
Goldmark has built-in syntax highlighting extensions:

```go
import (
    "github.com/yuin/goldmark-highlighting"
)

md := goldmark.New(
    goldmark.WithExtensions(
        extension.GFM,
        highlighting.NewHighlighting(
            highlighting.WithStyle("terminal"),
            highlighting.WithFormatOptions(
                chromahtml.WithClasses(true),
            ),
        ),
    ),
)
```

### Option B: Parse HTML with Go's html Package
```go
import "golang.org/x/net/html"

func addSyntaxHighlighting(htmlContent string) string {
    doc, err := html.Parse(strings.NewReader(htmlContent))
    if err != nil {
        return htmlContent
    }
    
    // Walk the HTML tree and find <code> elements
    var walk func(*html.Node)
    walk = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "code" {
            // Check if parent is <pre> and has language class
            // Apply syntax highlighting
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            walk(c)
        }
    }
    walk(doc)
    
    // Convert back to HTML
    // ...
}
```

## Let's Start with Debugging

Run the debug version first and share the output. This will tell us:

1. **What HTML goldmark is generating** - Is it the format we expect?
2. **Whether our regex is matching** - Are we finding the code blocks?
3. **What chroma is doing** - Is it highlighting correctly?
4. **Where the process breaks** - Regex, chroma, or somewhere else?

Once we see the debug output, we can decide whether to fix the regex or switch to a better approach like the goldmark extension.

What does the debug output show when you run it?