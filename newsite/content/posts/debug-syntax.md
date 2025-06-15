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

More OCaml examples:

```ocaml
(* Variant types and pattern matching *)
type shape = 
  | Circle of float
  | Rectangle of float * float
  | Triangle of float * float * float

let area = function
  | Circle r -> Float.pi *. r *. r
  | Rectangle (w, h) -> w *. h
  | Triangle (a, b, c) ->
      let s = (a +. b +. c) /. 2.0 in
      sqrt (s *. (s -. a) *. (s -. b) *. (s -. c))

(* Higher-order functions and modules *)
module StringSet = Set.Make(String)

let process_words words =
  words
  |> List.map String.lowercase_ascii
  |> List.filter (fun s -> String.length s > 3)
  |> StringSet.of_list
  |> StringSet.elements

(* Option type and error handling *)
let safe_divide x y =
  if y = 0.0 then None
  else Some (x /. y)

let calculate_average numbers =
  match numbers with
  | [] -> Error "Empty list"
  | nums ->
      let sum = List.fold_left (+.) 0.0 nums in
      let count = List.length nums |> Float.of_int in
      Ok (sum /. count)

(* Recursive data structures *)
type 'a tree = 
  | Leaf 
  | Node of 'a * 'a tree * 'a tree

let rec tree_map f = function
  | Leaf -> Leaf
  | Node (value, left, right) ->
      Node (f value, tree_map f left, tree_map f right)

(* Functors and advanced module system *)
module type COMPARABLE = sig
  type t
  val compare : t -> t -> int
end

module MakeSet(Ord: COMPARABLE) = struct
  type elt = Ord.t
  type t = elt list
  
  let empty = []
  
  let rec add x = function
    | [] -> [x]
    | h :: t as s ->
        match Ord.compare x h with
        | 0 -> s
        | n when n < 0 -> x :: s
        | _ -> h :: add x t
end
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

More ReasonML examples:

```reason
/* Variant types with modern syntax */
type httpMethod = 
  | GET
  | POST(string)
  | PUT(string, string)
  | DELETE(string);

type apiResponse('a) = 
  | Loading
  | Success('a)
  | Error(string);

/* Pattern matching with destructuring */
let handleResponse = (response) =>
  switch (response) {
  | Loading => "Please wait..."
  | Success(data) => "Got data: " ++ data
  | Error(msg) => "Error: " ++ msg
  };

/* JSX-like syntax for ReasonReact */
let make = (~name, ~age, ~children) => {
  ...component,
  render: (_self) =>
    <div className="person-card">
      <h2> {ReasonReact.string(name)} </h2>
      <p> {ReasonReact.string("Age: " ++ string_of_int(age))} </p>
      <div> children </div>
    </div>
};

/* Advanced functional programming */
module Option = {
  let map = (f, opt) =>
    switch (opt) {
    | None => None
    | Some(x) => Some(f(x))
    };
    
  let flatMap = (f, opt) =>
    switch (opt) {
    | None => None
    | Some(x) => f(x)
    };
    
  let getWithDefault = (default, opt) =>
    switch (opt) {
    | None => default
    | Some(x) => x
    };
};

/* Async/Promise handling */
let fetchUserData = (userId) => {
  Js.Promise.(
    Fetch.fetch("/api/users/" ++ userId)
    |> then_(Fetch.Response.json)
    |> then_(json => {
         let user = Decode.user(json);
         resolve(Success(user));
       })
    |> catch(error => {
         let message = Js.String.make(error);
         resolve(Error(message));
       })
  );
};

/* Belt standard library usage */
let processUsers = (users) => {
  users
  |> Belt.Array.keep(user => user.age >= 18)
  |> Belt.Array.map(user => {...user, name: String.capitalize(user.name)})
  |> Belt.Array.reduce(Belt.Map.String.empty, (acc, user) =>
       Belt.Map.String.set(acc, user.id, user)
     );
};

/* Recursive data structures with modern syntax */
type rec binaryTree('a) = 
  | Empty
  | Node({
      value: 'a,
      left: binaryTree('a),
      right: binaryTree('a),
    });

let rec insertIntoTree = (value, tree) =>
  switch (tree) {
  | Empty => Node({value, left: Empty, right: Empty})
  | Node({value: nodeValue, left, right}) =>
      if (value <= nodeValue) {
        Node({value: nodeValue, left: insertIntoTree(value, left), right});
      } else {
        Node({value: nodeValue, left, right: insertIntoTree(value, right)});
      }
  };

/* Interop with JavaScript */
[@bs.module] external moment: string => Js.Date.t = "moment";
[@bs.send] external format: (Js.Date.t, string) => string = "format";

let formatDate = (dateString) => {
  let date = moment(dateString);
  format(date, "YYYY-MM-DD");
};

/* Polymorphic variants */
type color = [
  | `Red
  | `Green  
  | `Blue
  | `RGB(int, int, int)
  | `HSL(float, float, float)
];

let colorToString = (color) =>
  switch (color) {
  | `Red => "red"
  | `Green => "green"
  | `Blue => "blue"
  | `RGB(r, g, b) => Printf.sprintf("rgb(%d, %d, %d)", r, g, b)
  | `HSL(h, s, l) => Printf.sprintf("hsl(%.1f, %.1f%%, %.1f%%)", h, s *. 100.0, l *. 100.0)
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