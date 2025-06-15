---
{
  "title": "Beyond TypeScript: Differences Between Typed Languages",
  "date": "2023-09-08T00:00:00Z",
  "tags": [
    "TypeScript",
    "OCaml",
    "JavaScript",
    "Melange"
  ],
  "summary": "A technical post about beyond typescript: differences between typed languages, covering TypeScript, OCaml, JavaScript, Melange"
}
---


For the past six years, I have been working with OCaml, most of this time has
been spent writing code at [Ahrefs](https://ahrefs.com/) to process [a lot of
data](https://ahrefs.com/big-data) and show it to users in a way that makes
sense.

OCaml is a language designed with types in mind. It took me some time to learn
the language, its syntax, and semantics, but once I did, I noticed a significant
difference in the way I would write code and colaborate with others.

Maintaining codebases became much easier, regardless of their size. And
day-to-day work felt more like having a super pro sidekick that helped me
identify issues in the code as I refactored it. This was a very different
feeling from what I had experienced with TypeScript and Flow.

Most of the differences, especially those related to the type system, are quite
subtle. Therefore, it is not easy to explain them without experiencing them
firsthand while working with a real-world codebase.

However, in this post, I will attempt to compare some of the things you can do
in OCaml, and explain them from the perspective of a TypeScript developer.

Before every snippet of code, we will provide links like this:
([try](https://melange.re/unstable/playground/?language=OCaml&code=bGV0ICgpID0gcHJpbnRfZW5kbGluZSAiaGVsbG8gd29ybGQi&live=off)).
These links will go either to the [TypeScript
playground](https://www.typescriptlang.org/play) for TypeScript snippets, or to
the [Melange playground](https://melange.re/unstable/playground), for OCaml
snippets. [Melange](https://melange.re/) is a backend for the OCaml compiler
that emits JavaScript.

Without further ado, let's go!

![beyond-typescript-01.jpg](/images/beyond-typescript-01.jpg)

*Photo by [Bernice Tong](https://unsplash.com/@bernicehtong) on
[Unsplash](https://unsplash.com/photos/VPTSmbGba7Q?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)*


## Syntax

OCaml's syntax is very minimal (and, in my opinion, quite nice once you get used
to it), but it is also quite different from the syntax in mainstream languages
like JavaScript, C, or Java.

Here is a simple snippet of code in OCaml syntax
([try](https://melange.re/unstable/playground/?language=OCaml&code=bGV0IHJlYyByYW5nZSBhIGIgPQogIGlmIGEgPiBiIHRoZW4gW10KICBlbHNlIGEgOjogcmFuZ2UgKGEgKyAxKSBiCgpsZXQgbXlfcmFuZ2UgPSByYW5nZSAwIDEw&live=off)):

```ocaml
let rec range a b =
  if a > b then []
  else a :: range (a + 1) b

let my_range = range 0 10
```

OCaml is built on a mathematical foundation called [lambda
calculus](http://www.inf.fu-berlin.de/lehre/WS03/alpi/lambda.pdf). In lambda
calculus, function definitions and applications don't use parentheses. So it was
natural to design OCaml with similar syntax to that of lambda calculus.

However, the syntax might be too foreign for someone used to JavaScript.
Luckily, there is a way to write OCaml programs using a different syntax which
is much closer to the JavaScript one. This syntax is called
[Reason](https://reasonml.github.io/) syntax, and it will make it much easier to
get started with OCaml if you are familiar with JavaScript.

Let's translate the example above into Reason syntax (you can translate any
OCaml program to Reason syntax from [the
playground](https://melange.re/unstable/playground/?language=Reason&code=bGV0IHJlYyByYW5nZSA9IChhLCBiKSA9PgogIGlmIChhID4gYikgewogICAgW107CiAgfSBlbHNlIHsKICAgIFthLCAuLi5yYW5nZShhICsgMSwgYildOwogIH07CgpsZXQgbXlSYW5nZSA9IHJhbmdlKDAsIDEwKTs%3D&live=off)!):

```reason
let rec range = (a, b) =>
  if (a > b) {
    [];
  } else {
    [a, ...range(a + 1, b)];
  };

let myRange = range(0, 10);
```

This syntax is fully supported throughout the entire OCaml ecosystem, and you
can use it to build:

- native applications if you need fast startups or high speed of execution
- or compile to JavaScript if you need to run your application in the browser.

To use Reason syntax, you just need to name your source file with the `.re`
extension instead of `.ml`, and you're good to go.

Since Reason syntax is widely supported and is closer to TypeScript than OCaml
syntax, we will use Reason syntax for all code snippets throughout the rest of
the article. Although understanding OCaml syntax has some advantages, such as
allowing us to understand a larger body of source code, blog posts, and
tutorials, there is absolutely no rush to do so, and you can always learn it at
any time in the future. If you're curious, we'll provide links to the Melange
playground for every snippet, so you can switch syntaxes to see how a Reason
program looks in OCaml syntax, or vice versa.

## Data types

OCaml has great support for data types, which are types that allow values to be
contained within them. They are sometimes called [algebraic data
types](https://en.wikipedia.org/wiki/Algebraic_data_type) (ADTs).

One example is tuples, which can be used to represent a point in a 2-dimensional
space
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSBwb2ludCA9IChmbG9hdCwgZmxvYXQpOwoKbGV0IHAxOiBwb2ludCA9ICgxLjIsIDQuMyk7&live=off)):

```reason
type point = (float, float);

let p1: point = (1.2, 4.3);

```

One difference with TypeScript is that OCaml tuples are their own type,
different from lists or arrays, whereas in TypeScript, tuples are a subtype of
arrays.

Let’s see this in practice. This is a valid TypeScript program
([try](https://www.typescriptlang.org/play?#code/DYUwLgBGCuAOoC4IG0DOYBOBLAdgcwBoJ1t8BdCAXhQCIAzAewZqJoCMBDDGsgbgCh+oSKBxUIACg5ISuPMjIBKKgD4IHAHSi8YABYCh4CNHGiJMeCEVA)):

```typescript
let tuple: [string, string] = ["foo", "bar"];

let len = (a: string[]) => a.length;

let u = len(tuple)
```

Note how the `len` function is annotated to take an array of strings as input,
but then we apply it and pass `tuple`, which has a type `[string, string]`.

In OCaml, this will fail to compile
([try](https://melange.re/unstable/playground/?language=Reason&code=bGV0IHR1cGxlOiAoc3RyaW5nLCBzdHJpbmcpID0gKCJmb28iLCAiYmFyIik7CgpsZXQgbGVuID0gKGE6IGFycmF5KHN0cmluZykpID0%2BIEFycmF5Lmxlbmd0aChhKTsKCmxldCB1ID0gbGVuKHR1cGxlKTsK&live=off)):

```reason
let tuple: (string, string) = ("foo", "bar");

let len = (a: array(string)) => Array.length(a);

let u = len(tuple)
//          ^^^^^
// Error This expression has type (string, string)
// but an expression was expected of type array(string)
```

Another data type is records. Records are similar to tuples, but each
"container" in the type is labeled.
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSBwb2ludCA9IHsKICB4OiBmbG9hdCwKICB5OiBmbG9hdCwKfTsKCmxldCBwMTogcG9pbnQgPSB7eDogMS4yLCB5OiA0LjN9Ow%3D%3D&live=off)):

```reason
type point = {
  x: float,
  y: float,
};

let p1: point = {x: 1.2, y: 4.3};
```

Records are similar to object types in TypeScript, but there are subtle
differences in how the type system works with these types. In TypeScript, object
types are structural, which means a function that works over an object type can
be applied to another object type as long as they share some properties. Here's
an example
([try](https://www.typescriptlang.org/play?#code/JYOwLgpgTgZghgYwgAgCoHsAm7kG8BQyyYwYANhAFzIDOYUoA5gNyHKYQ0IMAOJ6IanQYgWbAJ4Q4UaiACuAWwBG0VgF98+UJFiIUAZQAW6KDozY8bEuSq16TVkQ5de-QXZFiN+BALrFSCmQAXmQACjAsdGojEzMogEoQgD5kXxAadAoAOjJ0Rgio7OsKBNYfPzBiKOpzHFDcAJtqACIYdHQWgBp2Tm5gPmABVqVpbuRJaWoAJgAGaYBGZG8SiELsBKA)):

```typescript
interface Todo {
  title: string;
  description: string;
  year: number;
}

interface ShorterTodo {
  title: string;
  description: string;
}

const title = (todo: ShorterTodo) => console.log(todo.title);

const todo: Todo = { title: "foo", description: "bar", year: 2021 }

title(todo)
```

In OCaml, you have a choice. Record types are nominal, so a function that takes
a record type can only take values of that type. Let's look at the same example
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSB0b2RvID0gewogIHRpdGxlOiBzdHJpbmcsCiAgZGVzY3JpcHRpb246IHN0cmluZywKICB5ZWFyOiBpbnQsCn07Cgp0eXBlIHNob3J0ZXJUb2RvID0gewogIHRpdGxlOiBzdHJpbmcsCiAgZGVzY3JpcHRpb246IHN0cmluZywKfTsKCmxldCB0aXRsZSA9ICh0b2RvOiBzaG9ydGVyVG9kbykgPT4gSnMubG9nKHRvZG8udGl0bGUpOwoKbGV0IHRvZG86IHRvZG8gPSB7dGl0bGU6ICJmb28iLCBkZXNjcmlwdGlvbjogImJhciIsIHllYXI6IDIwMjF9OwoKdGl0bGUodG9kbyk7&live=off)):

```reason
type todo = {
  title: string,
  description: string,
  year: int,
};

type shorterTodo = {
  title: string,
  description: string,
};

let title = (todo: shorterTodo) => Js.log(todo.title);

let todo: todo = {title: "foo", description: "bar", year: 2021};

title(todo);
//    ^^^^
// Error This expression has type todo but an expression was expected of
// type shorterTodo
```

But if we want to use structural types, OCaml objects also offer that option.
Here is an example using `Js.t` object types in Melange
([try](https://melange.re/unstable/playground/?language=Reason&code=bGV0IHByaW50VGl0bGUgPSB0b2RvID0%2BIHsKICBKcy5sb2codG9kbyMjdGl0bGUpOwp9OwoKbGV0IHRvZG8gPSB7InRpdGxlIjogImZvbyIsICJkZXNjcmlwdGlvbiI6ICJiYXIiLCAieWVhciI6IDIwMjF9OwpwcmludFRpdGxlKHRvZG8pOwpsZXQgc2hvcnRlclRvZG8gPSB7InRpdGxlIjogImZvbyIsICJkZXNjcmlwdGlvbiI6ICJiYXIifTsKcHJpbnRUaXRsZShzaG9ydGVyVG9kbyk7&live=off)):

```reason
let printTitle = todo => {
  Js.log(todo##title);
};

let todo = {"title": "foo", "description": "bar", "year": 2021};
printTitle(todo);
let shorterTodo = {"title": "foo", "description": "bar"};
printTitle(shorterTodo);
```

To conclude the topic of ADTs, one of the most useful tools in the OCaml toolbox
are variants, also known as sum types or [tagged
unions](https://en.wikipedia.org/wiki/Tagged_union).

The simplest variants are similar to TypeScript
[enums](https://www.typescriptlang.org/docs/handbook/enums.html)
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSBzaGFwZSA9CiAgfCBQb2ludAogIHwgQ2lyY2xlCiAgfCBSZWN0YW5nbGU7&live=off)):

```reason
type shape =
  | Point
  | Circle
  | Rectangle;
```

The individual names of the values of a variant are called *constructors* in
OCaml. In the example above, the constructors are `Point`, `Circle`, and
`Rectangle`. Constructors in OCaml have a different meaning than the reserved
word
[`constructor`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes/constructor)
in JavaScript.

Unlike TypeScript enums, OCaml does not require prefixing variant values with
the type name. The type inference system will automatically infer them as long
as the type is in scope.

This TypeScript code
([try](https://www.typescriptlang.org/play?#code/KYOwrgtgBAygFgQwA7CgbwFBSgBQPYCWIALgDRZQDCBATgMYA2w52ASsHcQiAOZMYBfDBibEoAZ0QpxUALxQA2hXjJgAOnxEyyqeur0mLWLrXtO3PswwBdANxA)):

```ts
enum Shape {
  Point,
  Circle,
  Rectangle
}

let shapes = [
  Shape.Point,
  Shape.Circle,
  Shape.Rectangle,
];
```

Can be written like
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSBzaGFwZSA9CiAgfCBQb2ludAogIHwgQ2lyY2xlCiAgfCBSZWN0YW5nbGU7CgpsZXQgc2hhcGVzID0gW1BvaW50LCBDaXJjbGUsIFJlY3RhbmdsZV07&live=off)):

```reason
type shape =
  | Point
  | Circle
  | Rectangle;

let shapes = [Point, Circle, Rectangle];
```

Another difference is that, unlike TypeScript enums, OCaml variants can hold
data for each constructor. Let's improve the `shape` type to include more
information about each constructor
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSBwb2ludCA9IChmbG9hdCwgZmxvYXQpOwp0eXBlIHNoYXBlID0KICB8IFBvaW50KHBvaW50KQogIHwgQ2lyY2xlKHBvaW50LCBmbG9hdCkgLyogY2VudGVyIGFuZCByYWRpdXMgKi8KICB8IFJlY3QocG9pbnQsIHBvaW50KTsgLyogbG93ZXItbGVmdCBhbmQgdXBwZXItcmlnaHQgY29ybmVycyAqLwo%3D&live=off)):

```reason
type point = (float, float);
type shape =
  | Point(point)
  | Circle(point, float) /* center and radius */
  | Rect(point, point); /* lower-left and upper-right corners */
```

Something like this is possible in TypeScript using [discriminated
unions](https://www.typescriptlang.org/docs/handbook/typescript-in-5-minutes-func.html#discriminated-unions)
([try](https://www.typescriptlang.org/play?#code/C4TwDgpgBACg9gSwHbCgXigbysAhgcwC4oByeZYEgbigGM44AnAEwGdiBtJAVwFsAjCIwA0UHgKEBdKAF8qAKFCQoAYQSNaAG2gZseIqTUbt1OhBRDO4wSLF8bkmo1zME3dnYmNZCpdABKELSoujgExCSBwaaacADuQgAyEABmwFb2QqLWUjTcYJCM-gj4ABbpUFyZtjmM0nKK4NAAyqW4yhjkKFAAPqrqWtB9UcBUQA)):

```typescript
type Point = { tag: 'Point'; coords: [number, number] };
type Circle = { tag: 'Circle'; center: [number, number]; radius: number };
type Rect = { tag: 'Rect'; lowerLeft: [number, number]; upperRight: [number, number] };
type Shape = Point | Circle | Rect;
```

The TypeScript representation is slightly more verbose than the OCaml one, as we
need to use object literals with a `tag` property to achieve the same effect. On
top of that, there are greater advantages of variants that we will see just
right next.

## Pattern matching

Pattern matching is one of the killer features of OCaml, along with the
inference engine (which we will discuss in the next section).

Let's take the `shape` type we defined in the previous example. Pattern matching
allows us to conditionally act on values of any type in a concise way. For
example
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSBwb2ludCA9IChmbG9hdCwgZmxvYXQpOwp0eXBlIHNoYXBlID0KICB8IFBvaW50KHBvaW50KQogIHwgQ2lyY2xlKHBvaW50LCBmbG9hdCkgLyogY2VudGVyIGFuZCByYWRpdXMgKi8KICB8IFJlY3QocG9pbnQsIHBvaW50KTsgLyogbG93ZXItbGVmdCBhbmQgdXBwZXItcmlnaHQgY29ybmVycyAqLwoKbGV0IGFyZWEgPSBzaGFwZSA9PgogIHN3aXRjaCAoc2hhcGUpIHsKICB8IFBvaW50KF8pID0%2BIDAuMAogIHwgQ2lyY2xlKF8sIHIpID0%2BIEZsb2F0LnBpICouIHIgKiogMi4wCiAgfCBSZWN0KCh4MSwgeTEpLCAoeDIsIHkyKSkgPT4KICAgIGxldCB3ID0geDIgLS4geDE7CiAgICBsZXQgaCA9IHkyIC0uIHkxOwogICAgdyAqLiBoOwogIH07&live=off)):

```reason
type point = (float, float);
type shape =
  | Point(point)
  | Circle(point, float) /* center and radius */
  | Rect(point, point); /* lower-left and upper-right corners */

let area = shape =>
  switch (shape) {
  | Point(_) => 0.0
  | Circle(_, r) => Float.pi *. r ** 2.0
  | Rect((x1, y1), (x2, y2)) =>
    let w = x2 -. x1;
    let h = y2 -. y1;
    w *. h;
  };
```

Here is the equivalent code in TypeScript
([try](https://www.typescriptlang.org/play?#code/C4TwDgpgBACg9gSwHbCgXigbysAhgcwC4oByeZYEgbigGM44AnAEwGdiBtJAVwFsAjCIwA0UHgKEBdKAF8qAKFCQoAYQSNaAG2gZseIqTUbt1OhBRDO4wSLF8bkmo1zME3dnYmNZCpdABKELSoujgExCSBwaaacADuQgAyEABmwFb2QqLWUjTcYJCM-gj4ABbpUFyZtjmM0nKK4NAAyqW4yhjkKFAAPqrqWtB9UcAK8vRIrKi4jBC46FAAFKxtkMSt7RAAlMS16AB8WPJQUKxxCMC0pUsrmwB0+ltHJye0uKzQZIgoJITHLydZsBuIwkFAAAx3cEKAFvD6GAYmP4AwEQYGgqAAWVwwFKdxgAEkoAAqLE4vFgeLLVYQO7OVzuUQAJi2MJecM+I1+-1hcEmqDiC1ukDu+UKxTKwA44OkAFpTjS7rEEoxkmlpY4eey+VMoNcMMLaWKhBLyhwAIxyhX3ZVJVJSy1sgFAkFgwWk0pOqDMVK4biadJak4AemDUAAopMQdAIAAPNruYAIABu0CuQQA1sh8KIIKmwbi4NwyjhSghWHR3tAVkXNMwxHmhFBBFBZrh08wg3Qdag4wmpimICpSpndo3vAaaV7UeiwX2-QPU8PM2yZPI5EA)):

```typescript
type Point = { tag: 'Point'; coords: [number, number] };
type Circle = { tag: 'Circle'; center: [number, number]; radius: number };
type Rect = { tag: 'Rect'; lowerLeft: [number, number]; upperRight: [number, number] };
type Shape = Point | Circle | Rect;

const area = (shape: Shape): number => {
  switch (shape.tag) {
    case 'Point':
      return 0.0;
    case 'Circle':
      return Math.PI * Math.pow(shape.radius, 2);
    case 'Rect':
      const w = shape.upperRight[0] - shape.lowerLeft[0];
      const h = shape.upperRight[1] - shape.lowerLeft[1];
      return w * h;
    default:
      // Ensure exhaustive checking, even though this case should never be reached
      const exhaustiveCheck: never = shape;
      return exhaustiveCheck;
  }
};
```

We can observe how in OCaml, the values inside each constructor can be extracted
directly from each branch of the `switch` statement. On the other hand, in
TypeScript, we need to first check the tag, and then access the other properties
of the object. Additionally, ensuring coverage of all cases in TypeScript using
the `never` type can be more verbose, and functions may be more error-prone if
we forget to handle it. In OCaml, exhaustiveness is ensured when using variants,
and covering all cases requires no extra effort.

The best thing about pattern matching is that it can be used for anything: basic
types like `string` or `int`, records, lists, etc.

Here is another example using pattern matching with lists
([try](https://melange.re/unstable/playground/?language=Reason&code=bGV0IHJlYyBzdW1MaXN0ID0gbHN0ID0%2BCiAgc3dpdGNoIChsc3QpIHsKICAvKiBCYXNlIGNhc2U6IGFuIGVtcHR5IGxpc3QgaGFzIGEgc3VtIG9mIDAuICovCiAgfCBbXSA9PiAwCiAgLyogU3BsaXQgdGhlIGxpc3QgaW50byBoZWFkIGFuZCB0YWlsLiAqLwogIHwgW2hlYWQsIC4uLnRhaWxdID0%2BCiAgICAvKiBSZWN1cnNpdmVseSBzdW0gdGhlIHRhaWwgb2YgdGhlIGxpc3QuICovCiAgICBoZWFkICsgc3VtTGlzdCh0YWlsKQogIH07CgpsZXQgbnVtYmVycyA9IFsxLCAyLCAzLCA0LCA1XTsKbGV0IHJlc3VsdCA9IHN1bUxpc3QobnVtYmVycyk7CmxldCAoKSA9IEpzLmxvZyhyZXN1bHQpOwo%3D&live=off)):

```reason
let rec sumList = lst =>
  switch (lst) {
  /* Base case: an empty list has a sum of 0. */
  | [] => 0
  /* Split the list into head and tail. */
  | [head, ...tail] =>
    /* Recursively sum the tail of the list. */
    head + sumList(tail)
  };

let numbers = [1, 2, 3, 4, 5];
let result = sumList(numbers);
let () = Js.log(result);
```

## Type annotations are optional

If we wanted to write some identity function in TypeScript, we would do
something like
([try](https://www.typescriptlang.org/play?#code/MYewdgzgLgBAlgEwFwwDwBUB8AKAbgQwBsV0BKGAXkxnUpgMMuoYCgWAzAVzGCjnBicIAUwCSCbIhQYcDEuSo1yAbxYx1MAE7ConTWBgBtRNgCMABlIBdFgF8gA)):

```typescript
const id: <T>(val: T) => T = val => val

function useId(id: <T>(val: T) => T) {
    return [id(10)]
}
```

While TypeScript
[generics](https://www.typescriptlang.org/docs/handbook/2/generics.html) are
very powerful, they lead to really verbose type annotations. As soon as our
functions start taking more parameters, or increasing in complexity, the type
signatures length increases accordingly.

Plus, the generic annotations have to be carried over to any other functions
that compose with the original ones, making maintenance quite cumbersome in some
cases.

In OCaml, the type system is based on [unification of
types](https://www.cs.cornell.edu/courses/cs3110/2011sp/Lectures/lec26-type-inference/type-inference.htm).
This differs from TypeScript, and allow to infer types for functions (even with
generics) without the need of type annotations.

For example, here is how we would write the above snippet in OCaml
([try](https://melange.re/unstable/playground/?language=Reason&code=bGV0IGlkID0gdmFsdWUgPT4gdmFsdWU7CgpsZXQgdXNlSWQgPSBpZCA9PiBbaWQoMTApXTs%3D&live=off)):

```reason
let id = value => value;

let useId = id => [id(10)];
```

The compiler can infer correctly the type of `useId` is `(int => 'a) =>
list('a)`.

With OCaml, type annotations are optional. But we can still add type annotations
anywhere optionally, if we think it will be useful for documentation purposes
([try](https://melange.re/unstable/playground/?language=Reason&code=bGV0IGlkOiAnYSA9PiAnYSA9IHZhbHVlID0%2BIHZhbHVlOwoKbGV0IHVzZUlkOiAoaW50ID0%2BICdhKSA9PiBsaXN0KCdhKSA9IGlkID0%2BIFtpZCgxMCldOw%3D%3D&live=off)):

```reason
let id: 'a => 'a = value => value;

let useId: (int => 'a) => list('a) = id => [id(10)];
```

I can not emphasize enough how the simplification seen above, which only
involves a single function, can affect a codebase with hundreds, or thousands of
more complex functions in it.

## Immutability

JavaScript is a language where mutability is pervasive, and working with
immutable data structures often require using third party libraries or other
complex solutions.

Trying to obtain real immutable values in TypeScript is quite challenging.
Historically, it has been hard to prevent mutation of properties inside objects,
which was mitigated with `as const`.

But still, the way the type system has to be flexible to adapt for the dynamism
of JavaScript can lead to "leaks" in immutable values.

Let's see an example
([try](https://www.typescriptlang.org/play?#code/JYOwLgpgTgZghgYwgAgLIFcxwEYBsIBqcu6EAPACoB8yA3gFDJPIBuxpAXMhQNz0C+9eqEixEKAJIBbKZhz4iJctTqNmUCHAAmAexC4Anq3YQuvAUIR6AzmGTAu02VjyETZW1FABzGgF46YyUuACIAC2AQ5H4+eisQW2QpLgwXBXdPH397PikAOjYlZADwuDCQniA)):

```typescript
interface MutableValue<T> {
    value: T;
}

interface ImmutableValue<T> {
    readonly value: T;
}

const i: ImmutableValue<string> = { value: "hi" };

const m: MutableValue<string> = i;
m.value = "hah";
```

As you can see, even when being strict about defining the immutable nature of
the value `i` using TypeScript expressiveness, it is fairly easy to mutate
values of that type if they happen to be passed to a function that expects a
type similar in shape, but without the `readonly` flag.

In OCaml, immutability is the default, and it's guaranteed. Records are
immutable (like tuples, lists, and most basic types), but even if we can define
mutable fields in them, something like the previous TypeScript leak is not
possible
([try](https://melange.re/unstable/playground/?language=Reason&code=dHlwZSBpbW11dGFibGVWYWx1ZSgnYSkgPSB7dmFsdWU6ICdhfQp0eXBlIG11dGFibGVWYWx1ZSgnYSkgPSB7bXV0YWJsZSB2YWx1ZSA6ICdhfQoKbGV0IGk6IGltbXV0YWJsZVZhbHVlKHN0cmluZykgPSB7IHZhbHVlOiAiaGkiIH07CgpsZXQgbTogbXV0YWJsZVZhbHVlKHN0cmluZykgPSBpOwptLnZhbHVlID0gImhhaCI7&live=off)):

```reason
type immutableValue('a) = {value: 'a}
type mutableValue('a) = {mutable value : 'a}

let i: immutableValue(string) = { value: "hi" };

let m: mutableValue(string) = i;
m.value = "hah";
```

When trying to assign `i` to `m` we get an error: `This expression has type
immutableValue(string) but an expression was expected of type
mutableValue(string)`.

## No imports

This might not be as impactful of a feature as the ones we just went through,
but it is really nice that in OCaml there is no need to manually import values
from other modules.

In TypeScript, to use some function `bar` defined in a module located in
`../../foo.ts`, we have to write:

```typescript
import {bar} from "../../foo.ts";
let t = bar();
```

In OCaml, libraries and modules in your project are all available for your
program to use, so we would just write:

```reason
let t = Foo.bar()
```

The compiler will figure out how to find the paths to the module.

## Currying

[Currying](https://en.wikipedia.org/wiki/Currying) is the technique of
translating the evaluation of a function that takes multiple arguments into
evaluating a sequence of functions, each with a single argument. It is a feature
that might be more desirable for those looking into learning more about
functional programming.

While it is possible to use currying in TypeScript, but it becomes quite verbose
([try](https://www.typescriptlang.org/play?#code/MYewdgzgLgBAtgSwB4wLwwBQEMBcNoBOCYA5gJRoB8mARnoceVTDTANQwBEX7MWA3AChQkWDQCm4gGZp4yDJwggsAa3EATGMRgAhHQEV8WAK7BxnMgonSLQkdBjAsBAiFjopxsMCgJwmCgBvQRhQx3AHGXREJAV1BAAHLTAYAAtjODhjCFsQsIJxKGMCFKkFJxc3XIBfDDJ+IA)):

```typescript
const mix = (a: string) => (b: string) => b + " " + a;
const beef = mix("soaked in BBQ sauce")("beef");
const carrot = function () {
    const f = mix("dip in hummus");
    return f("carrot");
}();
```

In OCaml, all functions are curried by default. This is how a similar code would
look like
([try](https://melange.re/unstable/playground/?language=Reason&code=bGV0IG1peCA9IChhLCBiKSA9PiBiICsrICIgIiArKyBhOwpsZXQgYmVlZiA9IG1peCgic29ha2VkIGluIEJCUSBzYXVjZSIsICJiZWVmIik7CmxldCBjYXJyb3QgPSB7CiAgbGV0IGYgPSBtaXgoImRpcCBpbiBodW1tdXMiKTsKICBmKCJjYXJyb3QiKTsKfTsK&live=off)):

```reason
let mix = (a, b) => b ++ " " ++ a;
let beef = mix("soaked in BBQ sauce", "beef");
let carrot = {
  let f = mix("dip in hummus");
  f("carrot");
};
```

## Build native apps that run fast

One of the best parts of OCaml is how flexible it is in the amount of places
your code can run. Your applications written in OCaml can run natively on
multiple devices, with very fast starts, as there is no need to start a virtual
machine.

The nice thing is that OCaml does not compromise expressiveness or ergonomics to
obtain really fast execution times. As this [study
shows](http://blog.gmarceau.qc.ca/2009/05/speed-size-and-dependability-of.html),
the language hits a great balance between verbosity (Y axis) and performance (X
axis). It provides features like garbage collection or a powerful type system as
we have seen, while producing small, fast binaries.

## Write your client and server with the same language

This is not a particular feature of OCaml, as JavaScript has allowed to write
applications that run in the server and the client for years. But I want to
mention it because with OCaml one can obtain the upsides of sharing the same
language across boundaries, together with a precise type system, a fast
compiler, and an expressive and consistent functional language.

At Ahrefs, we work with the same language in frontend and backend, including
tooling like build system and package manager (we wrote about it
[here](https://tech.ahrefs.com/ahrefs-is-now-built-with-melange-b14f5ec56df4)).
Having the OCaml compiler know about all our code allows us to support several
number of applications and systems with a reasonably sized team, working across
different timezones.

---

I hope you enjoyed the article. If you want to learn more about OCaml as a
TypeScript developer I can recommend the Melange documentation site, which has
plenty of information about how to get started. This page in particular,
[Melange for X developers](https://melange.re/v1.0.0/melange-for-x-developers/),
summarizes some of the things we have discussed, and expanding on others.

If you want to share any feedback or comments, please comment [on
Twitter](https://twitter.com/javierwchavarri/), or join the [Reason
Discord](https://discord.gg/reasonml) to ask questions or share your progress on
any project or idea built with OCaml.
