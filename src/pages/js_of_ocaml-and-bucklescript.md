---
title: "Js_of_ocaml and BuckleScript"
subtitle: "A Comparison"
date: "2019-03-17"
tags:
  - "JavaScript"
  - "ReasonML"
  - "BuckleScript"
  - "OCaml"
---

When compiling Reason / OCaml applications to JavaScript, there are two main options: [BuckleScript](http://bucklescript.github.io/) and [Js\_of\_ocaml](http://ocsigen.org/js_of_ocaml/).

In this article, we will explain briefly the origins of each one of them, show the upsides and strong points of each solution, and lastly conclude with some guidelines on what use cases fit each one better.

![js_of_ocaml-and-bucklescript-01.png](/media/js_of_ocaml-and-bucklescript-01.png)

## History

#### Js\_of\_ocaml: leveraging bytecode stability

Back in 2013, Jérôme Vouillon and Vincent Balat published the paper [“From Bytecode to JavaScript: the Js\_of\_ocaml compiler”](https://www.irif.fr/~balat/publications/vouillon_balat-js_of_ocaml.pdf). The main idea was to use OCaml bytecode –an output that the OCaml compiler could produce already through the use of `ocamlc`– to produce JavaScript code. The main motivation was that OCaml bytecode is a very stable format, so building on top of it would decrease the maintenance costs of Js\_of\_ocaml. Another upside of using bytecode was that it would make possible to leverage all the existing ecosystem of libraries and applications and make them run in the most ubiquitous platform that ever existed: the browser.

![js_of_ocaml-and-bucklescript-02.png](/media/js_of_ocaml-and-bucklescript-02.png)

*High-level view of the OCaml compilation process, highlighting the part that is handled by Js\_of\_ocaml*

#### BuckleScript: leveraging language semantics

A couple of years later, Hongbo Zhang proposed [in the Js\_of\_ocaml repository](https://github.com/ocsigen/js_of_ocaml/issues/338) the idea of taking OCaml _rawlambda_ output to produce JavaScript code. This was a very different approach from the one taken by Js\_of\_ocaml, as _rawlambda_ is a data structure much more central in the compilation process. The motivation behind this approach was that JavaScript and OCaml share a lot of the language semantics, so by reaching into _rawlambda_, the compiler backend that would become BuckleScript could make the JavaScript output smaller, and keep it closer to the original OCaml code. This was impossible by design for Js\_of\_ocaml, which starts the conversion to JavaScript at the very end of the compilation process, because at that time a lot of information that would be valuable semantically has been erased already –like function names *(Edit: as noticed by [Louis Roché](https://medium.com/@TestCross), Js\_of\_ocaml maintains functions names just fine, as can be seen by the source maps feature)*.

![js_of_ocaml-and-bucklescript-02.png](/media/js_of_ocaml-and-bucklescript-02.png)

*High-level view of the OCaml compilation process, highlighting the part that is handled by BuckleScript, in comparison to Js\_of\_ocaml*

An so, two options to compile from OCaml to JavaScript were born. 😄

## How do they compare?

These two projects don’t compete against each other, even if it could seem so because they both produce JavaScript. In reality, their initial goals, target audiences and design decisions were so different, that this led to very different implementations and sets of trade-offs too.

The following is a non-exhaustive list of the benefits of each of them.

#### Upsides of Js\_of\_ocaml

- Integration with OCaml build system: [Dune](https://jbuilder.readthedocs.io/en/latest/) is the default tool to build native applications in OCaml. It has first-class support for Js\_of\_ocaml, which makes very easy to access all the existing features, documentation and knowledge.
- Seamless integration with `external`: for apps that need to be compiled to both native code using C externals, and also to the web using JavaScript externals, then Js\_of\_ocaml provides the best experience. The JavaScript functions can be declared with `external` and Js\_of\_ocaml will take care of linking the underlying JavaScript functions.
- Access all the existing the Ocaml libs: libraries like `compiler-libs` (which allows to work with the OCaml compiler, parsetree, etc) that are not available in BuckleScript, are very easy to use in Js\_of\_ocaml. This allows to create in-browser applications that are mind blowing, like [sketch.sh](http://sketch.sh/).
- Use all the existing ppxs with the most efficient workflow: Js\_of\_ocaml is fully compatible with ocaml-migrate-parsetree, ppx drivers, and the latest features available for pre-processors in the OCaml ecosystem.
- Source maps: a nice to have when debugging!
- In sync with the latest OCaml version: because of its non intrusive integration, Js\_of\_ocaml always keeps up with upstream changes in the OCaml compiler.

#### Upsides of BuckleScript

- Closer conversion between OCaml and JavaScript types: `array`, `string` and `boolean` have the same representation, while Js\_of\_ocaml needs to convert to `Js.js_array`, `Js.js_string`, etc. These conversions in Js\_of\_ocaml are annoying if you need to interop with JavaScript code very often.
- Easy consumption of JavaScript libraries: kind of a consequence of the former, but a very important point. Because BuckleScript keeps a closer conversion between OCaml and JavaScript types, it is much easier to interop with the existing JavaScript ecosystem. This has probably been the main reason behind the project momentum and growth, fueled mostly by JavaScript developers that want to start writing or migrating their web apps in Reason, and use BuckleScript to compile them.
- BuckleScript does not need a runtime, which leads to more performant code and a smaller payload.
- Integration with JavaScript bundlers (Webpack, Rollup, Parcel or any other JavaScript bundler).
- Can produce ES5 or ES6 output.
- Readable output: while BuckleScript [does not offer support for source maps yet](https://github.com/BuckleScript/bucklescript/issues/1699) (as of March 2019), the produced JavaScript output is really similar to the original OCaml code (same functions and variable names, etc).

## Use cases

So, when to use each one? Here are some examples of situations where you want to use one or the other.

#### Js\_of\_ocaml

An existing Reason / OCaml native application that needs to be adapted to run on the browser
A library or framework written in OCaml / Reason that needs to define externals to both C and JavaScript. For example: [Revery](https://github.com/bryphe/revery).
Applications that need access to the compiler internal data structures (parsetree, typedtree). For example: [sketch.sh](http://sketch.sh/).

#### BuckleScript

An existing JavaScript application that needs to be partially or fully migrated to OCaml / Reason.
A web application written in Reason or OCaml that needs access to many existing JavaScript libraries.
A library written in Reason / OCaml that is exposed to JavaScript consumers as well. For example: [Nact](https://nact.io/).

---

I hope you enjoyed the comparison, if you want to share your experiences with BuckleScript or Js\_of\_ocaml or have a suggestion, reach out [on Twitter](https://twitter.com/javierwchavarri/).
