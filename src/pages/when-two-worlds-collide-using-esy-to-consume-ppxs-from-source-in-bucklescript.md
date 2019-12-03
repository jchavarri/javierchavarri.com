---
title: "When two worlds collide: using esy to consume BuckleScript ppxs from source"
date: "2019-12-02"
imghero: "https://www.javierchavarri.com/media/when-two-worlds-collide-using-esy-to-consume-ppxs-from-source-in-bucklescript-01.jpg"
tags:
  - "esy"
  - "ppx"
  - "BuckleScript"
  - "ReasonML"
---

In this article, we will take a look at the existing landscape of ppx rewriters for BuckleScript. We will also see the work involved in preparing and publishing a ppx rewriter for their authors. Finally, we will present a different approach to publish and consume ppx rewriters in [BuckleScript](https://bucklescript.github.io/) that relies on the cross-platform package manager [esy](https://esy.sh/).

![performance-of-records-in-bucklescript-01.jpeg](/media/when-two-worlds-collide-using-esy-to-consume-ppxs-from-source-in-bucklescript-01.jpg"

## What is a ppx

In Reason and OCaml there is a system called "pre-processor extensions", which are also known as ppx rewriters, or just "ppxs". As their name states, these programs pre-process the code: they run just after the compilation parsing stage has successfully ended, which results in the compiler has an abstract syntax tree (or AST) available. And they run before the type checking and the more thorough parts of the compilation process start.

![performance-of-records-in-bucklescript-02.png](/media/when-two-worlds-collide-using-esy-to-consume-ppxs-from-source-in-bucklescript-02.png")

Ppxs were originally introduced when a command line option `-ppx` was added to the OCaml compiler. This option allows to pass multiple paths to native applications binaries. These programs are then executed during the compilation of each file. Once the compilation starts, each ppx binary file defined in the list is given a serialized representation of the AST by the compiler, and produces another serialized AST. This process continues for all the ppxs, and all the modules built by the compiler.

## How ppxs work in BuckleScript today

In BuckleScript, as a fork of OCaml, ppxs have been available since the very beginning. BuckleScript exposes the command line flag `-ppx` from the `bsconfig.json` file, which is used to configure the build process (one can think of it as `webpack.config.js`). The property to set the list of ppx rewriters that will run is `ppx-flags`, so if we want to run 3 ppxs we will write something like:

```json
"ppx-flags": [
  "graphql-ppx",
  "my-cool-ppx",
  "destroy-all-code-ppx"
],
```

So, if ppxs have been working for a long time in OCaml, and BuckleScript exposes a way to take advantage of that feature, where is the issue?

## Being a BuckleScript ppx author

BuckleScript publishes regular releases of the compiler itself. Pre-built binaries for the most used platforms: macOS, Linux and Windows.

It also provides a really familiar way —at least for JavaScript developers— to make use of npm to enable the publication of source code for libraries or bindings packages. They can be written in Reason or OCaml syntax, and then they get compiled with the rest of the application that they are part of. The decision of using npm was very successful, and BuckleScript application developers have access today to an increasing ecosystem of packages, published by hundreds of authors ([more than 800 results](https://www.npmjs.com/search?q=bucklescript) for the word `bucklescript` in npm).

However, as of today, the BuckleScript compiler does not support any kind of system to build and link _native_ libraries at compile time like it does with other packages. But who, you might ask, would want to write and use native libraries for a compiler that did abstract away from native platforms, by deciding to use as target an interpreted language like JavaScript? Maybe not _many_ people, but definitely some! One of the largest group in this category are ppx authors and their users.

The absence of a "native build story for libraries" in BuckleScript poses some challenges for ppx authors. If they want to publish a ppx rewriter, they have to start asking the same questions that BuckleScript itself has to solve for the compiler releases:

- "What OS and versions do I want to support?"
- "How do I prepare pre-build artifacts for them?"
- "How do I publish the binaries?"
- "Should I have one package with all binaries, or many?"

But ppxs are generally projects of different nature (and scale!) than a compiler.

A ppx author generally has many ongoing projects, has to maintain other libraries, or even other ppxs, and probably prefers to avoid spending much time solving problems related to publishing cross-platform binaries, which generally requires a different skills and knowledge than building and designing software.

Besides, ppx authoring gets more intimidating when we consider that most BuckleScript users come from JavaScript, where all the intricacies of native applications development are abstracted away. This mismatch results in ppxs remaining less accessible for a very large part of the community, which is in my opinion a shame, considering how creative and exploration-driven this community has been.

Even with those challenges, right now the BuckleScript ppx ecosystem is seeing a healthy activity, with several ppxs being published for multiple platforms. Which speaks volumes about their authors and maintainers passion, technical skills and generosity, and the community desire to experiment and explore new ways to improve the way we build apps.

These are some examples of ppxs for BuckleScript:

- [`graphql_ppx`](https://github.com/mhallin/graphql_ppx), probably the most known ppx in the community, publishes binaries for [Linux and macOS](https://github.com/mhallin/graphql_ppx/tree/5796b3759bdf0d29112f48e43a2f0623f7466e8a/ci)
- [`decco`](https://github.com/reasonml-labs/decco), a ppx that generates (de)serializers for user-defined types, publishes as well for [Linux and macOS](https://github.com/reasonml-labs/decco/blob/0a972f75c164d52dae5f5b26928312582470ac74/.travis.yml) as well
- [`bs-emotion-ppx`](https://github.com/ahrefs/bs-emotion) offers builds for [Linux, macOS and Windows](https://github.com/ahrefs/bs-emotion/blob/master/.ci/azure-pipelines.yml), thanks to Azure pipelines, and the work done by [Ulrik Strid](https://twitter.com/UlrikStrid/) and [Jordan Walke](https://twitter.com/jordwalke) in the [hello-reason](https://github.com/esy-ocaml/hello-reason) example repo.
- There are even extreme cases like [`bs-deriving`](https://github.com/ELLIOTTCABLE/bs-deriving), that take an existing ppx in the OCaml ecosystem ([ppx-deriving](https://github.com/ocaml-ppx/ppx_deriving)) and wraps it to make it usable from BuckleScript. It provide builds for [Linux and macOS](https://github.com/ELLIOTTCABLE/bs-deriving/blob/175e9575988d30b1cbcd0c2205078bdcc65a7db1/.travis.yml).

This covers just the surface, but I hope it gives a glimpse of the amount of effort that goes into building and maintaining these packages, and make them accessible to as many people as possible.

## Using ppxs from source

What if there was a way for ppx authors to publish their ppxs with less pain, but at the same time keeping the users experience as seamless as possible when using them?

If this was possible, then a library author could use the _exact_ same process that is used for regular BuckleScript libraries or bindings in order to publish a ppx: write the source code, `npm publish`. Boom. Done.

But how would the ppx applications get built?

## esy: the first cross-language package manager

[esy](https://esy.sh/) is a package manager that helps managing either native OCaml or Reason packages, or JavaScript packages. It does so by following closely the model that was so successful for JavaScript ecosystem: using a `package.json` file (it also allows to use `esy.json` alternatively).

Users that would like to consumer ppxs from source, or that are not published in their platform of choice, could use esy to have access to them. This would involve some steps, starting by [installing esy](https://esy.sh/docs/en/getting-started.html), if you have not installed it yet.

## Example: use a ppx `foo` from source

Let's say we want to consume an existing ppx `foo` that has been already published to npm so it can be consumed from source. We would add an `esy.json` file to our BuckleScript app:

```json
{
  "name": "my_bucklescript_app",
  "dependencies": {
    "foo": "^1.0.0"
  }
}
````

Then add the `ppx-flags` setting, as shown above, in `bsconfig.json`:

```json
"ppx-flags": [
  "esy x foo"
],
```

This is the main difference with the previous approach: instead of calling directly a pre-built binary that the ppx author prepared in advance, we are just telling BuckleScript to ask esy to run it.

After this, we can just build and run our BuckleScript app normally using `bsb -make-world` or just `bsb`.

## Future steps

Having a streamlined way to consume ppxs from source would also allow a few other things:

- It would allow to start exploring ways to explore and support better composition for ppxs, like the OCaml community did in the past with ppxlib. This kind of optimizations allow to link all ppxs together so all ppxs run as part of the same step avoiding most of the serialization and deserialization of the AST.
- It would empower more BuckleScript users to create and author their own ppxs and using native tooling, by building on top of a JavaScript friendly workflow like esy.

## The downsides

Of course, this approach has downsides as well, mostly for ppx consumers. It introduces another tool to manage the build process (esy) with everything that that entails.

One mitigation for this is to include `esy && bsb` in the build commands defined in `package.json`. `esy` has a very aggresive way to cache artifacts, taking roughly 100ms or less to run when all dependencies have been built already.

## Resources

The demo repo [`hello-ppx-esy`](https://github.com/jchavarri/hello-ppx-esy) has been updated with the ideas from this post. The repo contains a very small ppx to transform the `[%gimme]` extension into the number 42 (I promise: ppxs can do much more than that! 😆 ). You can find a sample BuckleScript project [here](https://github.com/jchavarri/hello-ppx-esy/tree/e53f8e8b5046bfb661e215c8c10f4c159a4df538/test_bs).

If you are interested on learning more about ppxs in OCaml, I recommend reading the blog post: ["An introduction to OCaml PPX ecosystem"](https://tarides.com/blog/2019-05-09-an-introduction-to-ocaml-ppx-ecosystem). 

If you want to learn more about how to publish cross-platform binaries from native libraries or apps, make sure to check the [`hello-reason`](https://github.com/esy-ocaml/hello-reason) which contains a very polished pipeline setup to build binaries for the three main OS platforms.