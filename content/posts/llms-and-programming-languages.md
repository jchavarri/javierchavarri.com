---
{
  "title": "LLMs, programming languages and other tools: a guess at the future",
  "date": "2025-01-16T00:00:00Z",
  "tags": [
    "LLMs",
    "Programming language design",
    "TypeScript",
    "OCaml",
    "React",
    "API Design",
    "Claude"
  ],
  "summary": "How LLMs as coding assistants are changing which programming languages, patterns, and architectures make sense. From the rise of explicit code over clever abstractions to the shifting economics of testing and dependencies.",
  "image": "/images/llms-and-programming-languages-01.jpg"
}
---

The rise of LLMs as coding assistants isn't just changing how we write code.
It's fundamentally altering which languages, patterns, and architectures make
sense in the first place. After spending the last months pair-programming with
Claude (and friends), I find myself thinking about the way programming languages
and libraries might be designed in the future, and how the priorities are
changing in a world where LLMs are doing more and more code reading and writing.

Here are some ideas and speculations about the future of programming languages,
based on my experience with LLMs.

![llms-and-programming-languages-01.jpg](/images/llms-and-programming-languages-01.jpg)

*Photo by [Mudassar Ahmed](https://unsplash.com/@mudassarahmed) on
[Unsplash](https://unsplash.com/photos/brown-dirt-road-between-green-grass-field-during-foggy-day-1CQVhYGIBZM)*

## Explicit is better than implicit (for real this time)

Here's something that would have sounded heretical (to me) five years ago:
boilerplate code isn't that bad anymore.

Take TypeScript vs OCaml as a [not so
random](https://www.javierchavarri.com/beyond-typescript/) example. In OCaml,
you can often get away with zero type annotations thanks to global type
inference. Write a function, and the compiler figures out the types by looking
at how it's used across your entire codebase. Elegant, concise, beautiful.

But here's the thing: LLMs don't see your entire codebase. Many times, they just
see the file you're working on, maybe a few related files in the context window.
When an LLM encounters an OCaml function with no type annotations, it has to
guess. When it sees a TypeScript function where every parameter is explicitly
typed, it knows exactly what's going on.

```typescript
// LLM-friendly: everything is explicit in this file
function processUser(user: User, config: ProcessingConfig): ProcessedUser {
  // ...
}
```

```ocaml
(* LLM-unfriendly: types inferred from elsewhere *)
let process_user user config = 
  (* ... *)
```

OCaml offers the possibility to include [interface
files](https://ocaml.org/docs/modules#interfaces-and-implementations) (`.mli`)
to type and document the API of a module. And this is a great way to make the
code more readable and understandable. But for LLMs, it means the types are not
colocated with the implementation, and the implementation itself is not
documented.

Most probably, a language designed today for LLMs would require having the types
and docs as close as possible to the implementation.

This isn't about verbosity for its own sake. It's about having all the necessary
context **locally available**. Languages that prioritized human cognitive
efficiency through clever abstractions and global inference are suddenly at a
disadvantage.

The same principle applies to implicit imports, monkey-patching, and dynamic
dispatch. All those smart language features that made code more concise for
humans make it harder for LLMs to understand what's actually happening. Some
humans always complained about this already ([basic code is easy to
maintain](https://www.mathieu-ferment.com/posts/basic-code-is-easy-to-maintain/)),
and they have now a much stronger argument for their demands.

## Context windows vs human cognition

React (with Flux) emerged because Facebook's engineers [couldn't keep track of
state mutations across their growing
codebase](https://www.youtube.com/watch?list=PLb0IAmt7-GS188xDYE-u1ShQmFFGbrkOv&t=622&v=nYkdrAPrdcw).
A notification bubble wouldn't update when it should have, because some distant
part of the application had modified shared state, and nobody could trace
through all the dependencies. I helped debug similar issues when I arrived at
Webflow, at a time when it was built with Knockout.js, before migrating to
React.

The solution was unidirectional data flow, immutable state, make re-rendering
automated. React's entire architecture is essentially a workaround for human
cognitive limitations.

But LLMs don't have those limitations, or at least they might get much better
than the standard human over time. A context window of 200K tokens can hold more
code than most humans can reason about in a week. Suddenly, the main original
problem that React solved (tracking mutations across a large codebase) becomes
much less of an issue.

To be clear, I'm not talking about React disappearing. React has organizational
benefits, ecosystem momentum, and solves problems beyond just state tracking. It
also has a very long history, adoption and ecosystem behind it.

But the technical motivations for many of our current tools might change, and
future ideas might have a harder time emerging just because the user pain that
helped past ideas get adoption will not be there anymore.

## The short tail

LLMs are essentially sophisticated pattern matchers trained on existing code.
They naturally gravitate toward common patterns and popular libraries.

This creates a feedback loop. Popular languages and libraries get used more in
LLM-generated code, which makes them even more popular, which means they appear
more in training data, which makes LLMs even more likely to suggest them.

We've seen this dynamic before in other domains: academic papers with more
citations get cited more, social platforms with network effects become
impossible to displace.

Programming languages already had network effects. But now, niche or new
libraries or programming languages will have an even harder time to find users
and gain adoption. Authors of these libraries and languages will have to either
provide a differentiator from mainstream choices in terms of integration with
LLMs, MCP and other AI-centric tools, or become even more aggressive on their
marketing efforts (I for one am not looking forward to seeing the latter 🤣).

## Cheaper tests

Like with type annotations, different languages have always had different
"testing taxes". By "testing tax", I mean the percentage of your codebase that
needs to be tests to feel confident about correctness, and how this "tax" varies
across languages (with all other things considered equal).

For example, for the sake of this article, we could say that Go codebases need
around 30-40% of test code because the language provides less compile-time
guarantees than other languages. Or OCaml codebases might have 10-15% tests
because the type system catches many errors upfront (I'm making these numbers
up).

But if LLMs write tests as easily as implementation code (or more), this dynamic
changes completely. The "unsafe" languages lose their disadvantage because the
safety net gets generated automatically. As long as the tests run fast, of
course!

So this could be another point with an impact on the language landscape.
Languages that previously required extensive testing to ensure correctness
suddenly become more appealing, when the cost of writing comprehensive tests
approaches zero.

## Bring your own code

My initial instinct was that LLMs would make dependency bloat worse: if they
default to popular libraries, we'd see even more bloated dependency trees.

But the opposite might also happen. Why import a library when you can generate
the exact solution you need? The cost of solving your specific problem is now
mostly zero in the majority of cases.

This could lead to leaner, more bespoke codebases. The ["not invented
here"](https://en.wikipedia.org/wiki/Not_invented_here) syndrome might actually
become a virtue. Instead of pulling in a 50KB library to format dates, just ask
the LLM to write exactly the date formatting you need.

Of course, this doesn't apply to everything. You're still going to want
battle-tested libraries for cryptography, database drivers, and other critical
infrastructure. But for the long tail of utility functions and domain-specific
logic, in many cases it will be possible to generate it mostly for free.

## DSLs vs natural language

Domain-specific languages have traditionally succeeded by offering concise,
domain-specific syntax that's more expressive than general-purpose languages for
specific problems. Think SQL for data queries, CSS for styling, or GraphQL for
API specifications.

But LLMs change this equation. When you can describe what you want in natural
language and get working general-purpose code, the value proposition of learning
a specialized syntax diminishes. The incentives to master a complex build system
DSL get significantly reduced when you can tell an LLM "create a build script
that compiles TypeScript, runs tests, and deploys to AWS".

This doesn't mean DSLs will disappear entirely. They'll likely survive in
domains where performance is critical (SQL isn't going anywhere), the domain
model is really well-established (CSS has decades of refinement), and/or human
readability and maintenance matter more than initial creation.

But like with libraries, the bar for creating new DSLs just got much higher. The
convenience factor that made many internal DSLs attractive is now competing with
the convenience of natural language generation.

## Native is the new black

With LLMs being able to translate between languages easily, I suspect that the
maintenance costs of cross-platform products (say, Android and iOS apps, or
macOS, Linux, Windows apps) might become less of a decision factor when choosing
between native development or a platform-agnostic / cross-platform framework.

Maybe in the future one could have a "source of truth" codebase (say, the iOS
app). Then, maintaining the Android app is mostly a matter of asking an LLM to
translate the latest changes and adapt the specific parts that need custom
implementation, and ensuring the result works as expected.

Frameworks like Flutter or React Native might have a harder time finding users,
both because these costs are reduced and also because the "native" languages are
much more prevalent (so there's more code to train the LLMs on, results will be
better, etc).

This multi-platform development space might be the first place that starts
seeing the rise of the first AI-native programming languages, with human defined
behavior that gets compiled down to platform-specific code.

## Language evolution speculation

Looking ahead, we might see new emerging languages that are optimized for LLM
consumption rather than human ergonomics. Features like ultra-explicit syntax,
mandatory documentation and tests, and standardized patterns that LLMs can
recognize and manipulate.

We might see languages with "compiled explicitness" where code is concise for
humans but expands to be explicit for LLMs, or "layered abstraction" languages
where humans write high-level intent and LLMs fill in the implementation
details. I can't wait to see a language with `.hmn` and `.llm` file extensions!
If you are working or know about any such languages, I'd love to learn more.

## The human element

Humans will continue to work on architecting systems, making design decisions,
and understanding business requirements and translate them to technical specs.
There will be new tasks, like a lot of reviewing and maintaining LLM-generated
code. Or debugging complex edge cases and performance issues.

The value will shift from knowing how to implement to knowing what to ask for
and knowing how to verify the result.

Joel Spolsky
[said](https://www.joelonsoftware.com/2002/11/11/the-law-of-leaky-abstractions/):

> (...) the abstractions save us time working, but they don't save us time
> learning.

I believe the same applies to LLMs.

## Conclusion

It is clear that our tools have already changed, and will continue to do so. The
criteria that mattered for a large part of developers (conciseness, elegance,
cognitive load) are being supplemented by criteria that matter for LLM
collaboration (explicitness, local context, predictable patterns).

A lot of software authors working on libraries (in the traditional sense) will
have a harder time finding new users for them, but at the same time there will
be plenty of opportunities to work on new tooling to enhance the collaboration
with LLMs, improve documentation, testing and many of the other areas and
workflows involved on product development.

Languages evolve quite slowly as they carry with them a lot of weight (tooling,
libraries, existing codebases). But eventually they will catch up. I am
particularly excited about existing languages that best fit the LLM-human
collaboration new needs, and even more about potential new languages that are
created to take the new workflows to higher grounds.

---

What do you think? Are you seeing similar patterns in your own experience with
LLMs? I'd love to hear your thoughts, feel free to reach out on
[Twitter](https://twitter.com/javierwchavarri)!
