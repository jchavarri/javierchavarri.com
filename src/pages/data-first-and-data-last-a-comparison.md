---
title: "Data-first and data-last: a comparison"
date: "2019-05-05"
tags:
  - "OCaml"
  - "ReasonML"
  - "BuckleScript"
---

[BuckleScript](https://bucklescript.github.io/) is a very interesting project: it takes the compiler from one language, OCaml, and modifies it in a way  that it becomes more ergonomic for users of another language: JavaScript. Because of this special nature, it has been received with some skepticism [since the very beginning](https://github.com/ocsigen/js_of_ocaml/issues/338).

One of the most relevant decisions —and probably one of the most controversial ones as well— was [to choose a "data-first" API for Belt](https://github.com/BuckleScript/bucklescript/issues/2463) (BuckleScript's standard library), as well as introducing a ["pipe first" operator](https://reasonml.github.io/docs/en/pipe-first) (`|.` in OCaml syntax, `->` in Reason syntax) to make it easier to work with Belt functions.

The context and constraints for this decision are quite nuanced, but most of the information has been spread in short comments in forums, pull requests, StackOverflow answers, etc., but I thought it would help to have a place that goes through the problem in detail, and backing the explanations with as many examples as possible.

In this post, we will go through the reasons that led to this decision, see examples that illustrate the problem, and lastly evaluate the trade-offs of it.

## Data-last: a traditional convention in functional languages

One of the common conventions in functional languages is to always pass the main parameter —the data that will be processed by the function— as the last parameter. This is known as "data-last", or in the OCaml ecosystem, because it's idiomatic to use `t` as the main type of a module, "t-last". In the rest of the post we will refer to this as data-last and data-first, but it is the same thing.

If we are using the OCaml standard library for example, and we want to map over the values of a list, we will do something like this (in Reason syntax):

```reason
let listOne = [1, 2, 3];
let listTwo = List.map(a => a + 1, listOne); // [2, 3, 4]
```

To understand the reasoning behind this convention it is important to understand currying.

### Currying and partial application

Currying means that all functions in the language only have one parameter. Functions that take "multiple parameters" are in reality functions that return other functions.

For example, the functions `f` and `g` below are equivalent:

```reason
let f = (a, b, c) => a + b + c;
let h = f(1, 2, 3); // 6

let g = a => b => c => a + b + c;
let i = g(1, 2, 3); // 6
```

Currying enables to partially apply functions, so one can write:

```reason
let add = (a, b) => a + b;
let addTwo = add(2);
let t = addTwo(2); // 4
let s = addTwo(5); // 7
```

Continuing the list mapping example above, we could abstract the function that adds 1 to all elements by taking advantage of currying:

```reason
let addOneToList = List.map(a => a + 1);
let listA = [1, 2, 3];
let listB = addOneToList(listA); // [2, 3, 4]
```

This is a very powerful (de)composition mechanism. For example, we can abstract `a => a + 1` in a new `plusOne` function for reusability purposes:

```reason
let plusOne = a => a + 1;
let addOneToList = List.map(plusOne);
let listA = [1, 2, 3];
let listB = addOneToList(listA); // [2, 3, 4]
```

### The pipe operator `|>`

The final push for the adoption of "data-last" convention was the pipe operator, which was [originally introduced](https://blogs.msdn.microsoft.com/dsyme/2011/05/17/archeological-semiotics-the-birth-of-the-pipeline-symbol-1994/)  in [Isabelle](https://en.wikipedia.org/wiki/Isabelle_%28proof_assistant%29), a theorem prover written in StandardML.

The main problem the pipe operator was trying to solve is when applying many functions one after each other, chaining the result of one function to the parameters of the next becomes was quite verbose and tedious.

The pipe operator solved this problem, as it passes the value on its left side as the last parameter of the expression that is placed at the right side.

So:

```reason
let filtered = list |> List.filter(a => a > 1);
```

is equivalent to:

```reason
let filtered = List.filter(a => a > 1, list);
```

_Side note: in practice, it is not exactly translated with such a straight forward conversion, as there is an extra function call involved, as we will see below._

To see the impact the pipe operator can have in readability and conciseness, here is a more complex example.

Instead of writing:

```reason
let getFolderSize = folderName => {
  let filesInFolder = filesUnderFolder(folderName);
  let fileInfos = List.map(fileInfo, filesInFolder);
  let fileSizes = List.map(fi, leSize, fileInfos);
  let totalSize = List.fold((+), 0L, fileSizes);
  let fileSizeInMB = bytesToMB(totalSize);
  fileSizeInMB;
};
```

one can write:

```reason
let getFolderSize = folderName =>
  folderName
  |> filesUnderFolder
  |> List.map(fileInfo)
  |> List.map(fileSize)
  |> List.fold((+), 0)
  |> bytesToMB;
```

A significant simplication!

The pipe operator allows us to pass the result of each expression as an input to the next one, without having to name the result of each step and explicitly pass it to the next function call.

### A convention in many functional languages

We have seen how currying enables partial application, and how this feature allows to compose functions easily, which was at the origin of the data-last convention. 

Also, we saw how the pipe operator contributes to adopt this convention by passing the result of an expression as the last argument of the next.

The data-last convention is not exclusive of OCaml and Reason, many other languages like [Elm](https://package.elm-lang.org/help/design-guidelines#the-data-structure-is-always-the-last-argument), [F#](https://fsharpforfunandprofit.com/posts/partial-application/#designing-functions-for-partial-application), [Haskell](https://wiki.haskell.org/Parameter_order) or even JavaScript libraries like [Ramda](https://ramdajs.com/) adopted it.

## Data-first: a different approach

So why then did BuckleScript decide to move away from this convention, towards the "data-first" approach?

To understand it, we have to go first through a short trip through type inference and how the type checker evaluates code.

### Type inference: creating truth, one step at a time

In OCaml, type inference works from left to right, and from top to bottom. Here is a simple example that shows it:

```reason
let aList = [2, "a"];
                ^^^
```
```
Error: This expression has type string but an expression was expected of type int
```

We can see how the compiler gets to analyze the integer `2` first, so it takes that as the "truth": `aList` has type `list(int)`.

So when it encounters the second element of the list, the string `"a"`, it checks it against that truth. At that point, the compilation process fails because `string` and `int` are different types.

This might sound pretty obvious, probably because some of us might be more used to left-to-right written languages. But one could imagine a compiler that would analyze programs in a different way. Maybe. 😂

What does this have to do with "data-first" or "data-last"? A lot, as it turns out.

Let's see a small example:

```reason
let strList = ["hi"]
let res = List.map(n => n + 1, strList)
                               ^^^^^^^
```
```
Error: This expression has type list(string)
but an expression was expected of type list(int)
Type string is not compatible with type int 
```
In this example, the compiler assumes the callback `n => n + 1` is the truth, so it infers we are dealing with a `list(int)`. Then it finds a `list(string)`, and fails.

However, if we are working with an API where the data comes first (also known as "t-first"), like Belt does:

```reason
let strList = ["hi"]
let b = Belt.List.map(strList, n => n + 1)
                                    ^
```
```
Error: This expression has type string but an expression was expected of type int
```

Note the difference: in this case, **the compiler assumes the type of `strList`, `list(string)`, is the truth**, and then it fails then when the callback returns an `int` type. Note also how the error message is simpler: the compiler is not matching a `list(int)` against a `list(string)` like in the first case, because it is operating already "inside" the callback, it can match the `string` —the original type of `a`— against an `int`.

This might not seem a big deal in this small example, but for real applications where the functions and data become more complex, the errors can become quite more cryptic with the data-last approach, because the compiler is assuming the source of truth is that of the "lifted" types of the callback: list, maps, options and any other "wrapping" types that are used in those callbacks.

### Annotations needed

In some cases, the compiler might not even be able to infer the types of data-last functions.

Let's say we have a module `User` with the following implementation:

```reason
module User = {
  type t = {
    name: string,
    age: int,
  };
  let admins = [
    {name: "Jane", age: 30},
    {name: "Manuel", age: 72},
    {name: "Lin", age: 54},
  ];
}
```

Now, outside this module, we want to get a list with the ages of the `admins` users. We use OCaml standard library `List` function `List.map`:

```reason
module User = { ... }
let ages = List.map(u => u.age, User.admins);
                           ^^^
```
```
The record field age can't be found.
If it's defined in another module or file, bring it into scope by:
  - Annotating it with said module name: let baby = {MyModule.age: 3}
  - Or specifying its type: let baby: MyModule.person = {age: 3}
```

Whoa whoa there compiler... "_annotating it_"? "specifying its type"?! I was promised OCaml had such a powerful inference engine that I would never need to write any more type annotations! 😄

It seems the compiler can't figure out that we want to get the value of the field `age` from the record of type `User.t`, even if it has `users`, of type `list(User.t)` right there, next to it.

We can solve the problem by adding a type annotation, as suggested by the compiler error:

```reason
let ages = List.map((u: User.t) => u.age, User.admins);
```

This is a consequence of the way type inference works: as we saw, type checking is done left to right, so when the compiler evaluates the `map` callback `u => u.age`, in the case without type annotations, it has no information about what `u` is.

Maybe if we used the pipe operator, it would work? `User.admins` would appear first in that case. 🤔

Let's see:

```reason
let ages = User.admins |> List.map(u => u.age);
                                          ^^^
```
```
Error: The record field age can't be found.
```

Still the same issue.

This doesn't work because the pipe operator is an [infix](https://en.wikipedia.org/wiki/Infix_notation) operator, which is a fancy way of saying it's like a function that takes two parameters, with the "infix" meaning each parameter is placed at each side of the operator.

If we wrote it as a plain function `pipeOp`, the pipe operator is equivalent to something like:

```reason
let ages = pipeOp(User.admins, List.map(u => u.age));
```

`User.admins` appears first, but the type checker still analyzes the callback body _before_ evaluating the `map` function as a whole, so it still doesn't have enough information to know where `age` is coming from.

### No annotations needed

With a data-first approach to API design, the need for a required annotation goes away:

```reason
let ages = Belt.List.map(User.admins, u => u.age);
```

This compiles just fine, without annotations needed! ✨

The compiler now can infer that the `u` expression in the callback parameter has type `User.t`, and so when it sees the `u.age` expression on the right side, it can be 100% sure where it comes from and check that it is valid.

### The pipe first operator `->`

In the same vein as the pipe operator `|>`, BuckleScript introduced a pipe operator that is similar, but instead passes the resulting value of the expression in the left side as the _first_ parameter of the one on the right. 

So, for example:

```reason
let filtered = list -> Belt.List.filter(a => a > 1);
```

is equivalent to:

```reason
let filtered = Belt.List.filter(list, a => a > 1);
```

Another important difference from the `|>` operator is that `->` is not an infix operator, just syntactic sugar, so it is really as if you were writing the second form instead of the first from the compiler perspective. With the traditional pipe `|>` it is interpreted like applying a function.

## Advantages of data-first

As we have seen, data-first remains more closely to how type inference works. We also saw how this:

- Helps the compiler infer types in functions that take callbacks as parameters, without having to manually add type annotations.
- Makes error messages simpler

This might not seem like a very big deal on its own, but it has a lot of impact down the line that affects the resulting developer experience.

### IDE integration

A direct consequence of the data-first approach that benefits from the compiler having more information is that we get more help from editor extensions when writing our functions.

In the example above, as we are writing it:

```reason
let ages = Belt.List.map(User.admins, u => u.
                                             ◻️ age
                                             ◻️ name
```

The editor, as the compiler, can know that `u` is of the record `User.t` and can provide autocompletion for the fields in it. Very helpful!

The advantages of data-first when it comes to editor integration is something that language designers with a vast experience, like Anders Hejlsberg (lead architect of TypeScript), [have explained in the past](https://github.com/Microsoft/TypeScript/issues/15680#issuecomment-307571917).

### Intuitive design for functions with multiple params

One of the downsides of the data-last approach is that sometimes it makes harder to understand what a function with two operands is doing.

For example, the `String.concat` function in OCaml standard library:

```reason
let foo = String.concat("a", ["b", "c"])
```

Because we are used to left-to-right reading —like the  inference engine— we could guess the resulting value of `foo` is the string `"abc"`, but it's actually `"bca"` because the API is data-last.

Another example is the `Js.String.split` in BuckleScript:

```reason
let bar = Js.String.split("a|b|c", "|")
```

The resulting value is the string `"|"`, because the `Js` module in BuckleScript was originally designed to follow the data-last approach.

So, if we are not using currying and the pipe operator, we have to think kind of "backwards" when using these functions.

### Performance

Pipe first operator `->` is implemented as purely syntactic sugar, as mentioned above. This means that, from the compiler perspective, the usage of `->` means that no extra functions calls are involved.

This is not what happens with the pipe last operator `|>`, that gets compiled into a function call. While the compiler does a lot of optimizations to inline values whenever possible, this difference makes it hard to optimize in some specific cases.

Here is an example of a binding to a JavaScript function using BuckleScript interop capabilities that shows how this can impact the resulting code.

With pipe last, using `[@bs.send.pipe]` annotations:

```reason
[@bs.send.pipe: user]
external update: (~isAdmin: bool=?, ~age: int) => user = "";

let jane = jane |> update(~age=45);
```

Which results in the following generated JavaScript code:

```javascript
var jane$1 = (function(param) {
  return function(param$1) {
    return param$1.update(
      param !== undefined ? Caml_option.valFromOption(param) : undefined,
      45
    );
  };
})(undefined)(jane);
```

Now, with pipe first:

```reason
[@bs.send]
external update: (user, ~isAdmin: bool=?, ~age: int, unit) => user = "";

let jane = jane->update(~age=45, ());
```

Results in this JavaScript code:

```javascript
var jane$1 = jane.update(undefined, 45);
```

Quite simpler! We went from two function definitions and a curried application in both of them, to a one-liner with just one function call.

One the other hand, this shows one of the main disadvantages of the data-first approach: with data-last and optional arguments, one doesn't need to put an extra `unit` argument at the end of the function: whenever the data is passed (`jane` in the example above) the compiler knows the function can be applied, and thus set all optional values to `None`.

With data-first this is not possible: the data comes first and then all the optional values, so we are forced to always include a `unit` type as last param to make sure the compiler knows when the function has been fully applied.

### Usages in OCaml and other languages

As we saw with the pipe operator `|>`, there are many functional languages that have followed the data-last approach. The data-first convention and pipe first operator are not as common as data-last, but there are some appearances too in functional and object-oriented languages. 

- [Elixir](https://hexdocs.pm/elixir/Kernel.html#%7C%3E/2) includes in the core language a pipe operator that passes the value as the first param called pipe-forward.
- There are some [commonly used libraries](https://ocaml.janestreet.com/ocaml-core/latest/doc/core_kernel/Core_kernel/Map/) in OCaml that have adopted a data-first approach in their design. Even some OCaml imperative APIs use it, like the ones found in [`Hashtbl`](http://caml.inria.fr/pub/docs/manual-ocaml/libref/Hashtbl.html).
- In other cases, like C#, there are [proposals](https://github.com/dotnet/csharplang/issues/96) to include pipe-first operator in the future.

## Conclusion

So, if you have read until here, first, you're awesome 😄 and second, what are the conclusions?

From all the data seen above, there is no clear "better way", both data-first and data-last have their own set of trade-offs.

Data-last has:
- A long tradition in functional languages
- Great integration with partially applied functions
- More straight-forward composition
- A simpler solution for application of functions with optional labelled arguments

Data-first provides:
- Simpler compiler errors
- More accurate type inference
- Smaller compiled output
- Better IDE integration
- More intuitive for left-to-right readers

In the end, I would probably say that going one way or another largely depends on what are the values, intention and audience of the specific language or libraries.

In BuckleScript's case, I think it made sense to go with the data-first approach, as it is targeting developers that come from JavaScript —an uncurried, object-oriented language— and are finding the new language for the first time. Because JavaScript is not a curried language, these developers might not find that much value in the advantages of data-last, while the more straight forward inference and better error messages and editor integration can be very helpful.

---

Thanks for reading! I hope the goal of the post was accomplished and it helped make clearer what the rationale was behind this decision. If you want to share any feedback, leave a comment on the orange site or reach out [on Twitter](https://twitter.com/javierwchavarri/). 

Keep shipping! 🚀
