---
title: "Js_of_ocaml: a bundle size study"
subtitle: "An analysis about JavaScript bundle size in Js_of_ocaml with a real world app"
date: "2022-01-23"
imghero: "https://www.javierchavarri.com/media/js_of_ocaml-bundle-size-study-01.jpg"
tags:
  - "JavaScript"
  - "Js_of_ocaml"
  - "ReScript"
  - "OCaml"
---

Historically, I have noticed a recurring theme in OCaml forums and Discord, where from time to time someone mentions that [Js\_of\_ocaml](http://ocsigen.org/js_of_ocaml/) —a compiler that converts from OCaml bytecode to JavaScript— generates large output files.

However, my experience was generally the opposite. When using Js\_of\_ocaml  [production builds](https://dune.readthedocs.io/en/stable/jsoo.html?highlight=mode%20release#separate-compilation) (using `--profile=prod` flag), the resulting JavaScript artifacts were quite small, as the compiler applies very aggressive dead code optimizations.

It is true though that Js\_of\_ocaml allows to use many OCaml libraries that were written with native use cases in mind, where binary size is not generally an issue. Js\_of\_ocaml also requires some conversions between OCaml native types and JavaScript types, and the fact that it generates JavaScript from bytecode makes the whole process more opaque.

So, are Js\_of\_ocaml generated files really that large, as the rumors suggest? I ran a small experiment to find some answers.

![js_of_ocaml-bundle-size-study-01.jpg](/media/js_of_ocaml-bundle-size-study-01.jpg)
  
*Photo by [Kai Dahms](https://unsplash.com/@dilucidus) on [Unsplash](https://unsplash.com/s/photos/measure)*


## Hypothesis

The main theory that I wanted to prove is that Js\_of\_ocaml produces reasonably sized JavaScript files. I also was interested about tracking the evolution in size of these output files over time, as the application keeps being developed and grows. If the output file is small for small apps, but grows too quickly over time, it would mean Js\_of\_ocaml would not be suitable for web applications that have limited bundle size budgets, or products that should grow sustainably in the long term.

In order to answer the above question, I thought it would be nice to use one of the most efficient compilers to JavaScript that exist out there: [ReScript](https://rescript-lang.org/). Which happens to be very close to OCaml as well. In a [previous article](https://www.javierchavarri.com/js_of_ocaml-and-bucklescript/) I compared both solutions and the trade-offs between them.

## The experiment

To run the experiment, I looked for an existing ReScript application that had some functionality that is common to most web applications:
- data fetching and pushing from/to servers
- routing
- parsing
- some work with collections and data processing

I found a good candidate in [jihchi/rescript-react-realworld-example-app](https://github.com/jihchi/rescript-react-realworld-example-app). This application uses ReScript and [rescript-react](https://rescript-lang.org/docs/react/latest/introduction) —the ReScript bindings to [React.js](https://reactjs.org/)— to build another example of the ["mother of all demo apps"](https://github.com/gothinkster/realworld). This demo app is a social blogging site (i.e. a Medium.com clone) that uses a custom API for all requests, including authentication.

To compare apples to apples, the plan would be to migrate this application fully to Js\_of\_ocaml. Then, as the application consists on a dozen of screens or so, there would be an easy way produce JavaScript files and take measurements of the output files progressively, as new screen components get added to the application.

To do the experiment we would leverage [jsoo-react](https://github.com/reason-in-barcelona/jsoo-react) to replicate the behavior found in the original ReScript app. `jsoo-react` are the bindings to React.js for Js\_of\_ocaml. They were originally based on `rescript-react`, but over time grew apart, with `jsoo-react` having more emphasis on supporting OCaml syntax.

## Methodology

Both bundles are generated using Webpack v4.46.0 in production mode, to avoid differences in minimization or bundling.
The bundles are analyzed using [webpack-bundle-analyzer](https://github.com/webpack-contrib/webpack-bundle-analyzer).

The results shared below were created from these commits:

- Js\_of\_ocaml: [0138bfe](https://github.com/jchavarri/jsoo-react-realworld-example-app/commit/0138bfedddc1e57237ffe9a9a53a07aea9f73bf6)
- ReScript: [2c319d9](https://github.com/jchavarri/jsoo-react-realworld-example-app/commit/2c319d933a5025e11a0a12e6abb804130659cc7a)

For each case, all components in [the main `App` component](https://github.com/jchavarri/jsoo-react-realworld-example-app/blob/0138bfedddc1e57237ffe9a9a53a07aea9f73bf6/src/app.ml#L20-L28) were commented, and then uncommented progressively, while running these commands on each step of the way:

- Js\_of\_ocaml: `make build-prod && yarn webpack:analyze`
- ReScript: `yarn build  && yarn webpack:analyze`

## Results

The results are published together with the Webpack bundle analyzer reports in
https://jchavarri.github.io/jsoo-react-realworld-example-app/bundle-study/.

Note the `Stat` column does not provide very valuable information. The reason why it is so large in ReScript case is that Js\_of\_ocaml generated JavaScript has already been minified by the Js\_of\_ocaml compiler (because of the `profile=prod` flag) while ReScript produces human readable JavaScript and has no minified. In any case, in the end Webpack minifier runs in both cases, so the `Parsed` and `Gzipped` columns are more useful.

#### Js\_of\_ocaml

| Component   | Stat        | Parsed      | Gzipped     |
| ----------- | ----------- | ----------- | ----------- |
| No page components | 55.15 KB |	51.61 KB |	17.54 KB |
| + Home               | 68.95 KB |	65.14 KB |	21.47 KB |
| + Settings           | 75.45 KB |	71.48 KB |	23.37 KB |
| + Login              | 77.86 KB |	73.85 KB |	23.96 KB |
| + Register           | 81.01 KB |	77.07 KB |	24.55 KB |
| + CreateArticle      | 87.09 KB |	82.97 KB |	26.11 KB |
| + EditArticle        | 87.17 KB |	83.05 KB |	26.13 KB |
| + Article            | 109.45 KB |	104.55 KB |	31.67 KB |
| + Profile            | 115.52 KB |	110.29 KB |	33.38 KB |
| + Favorited          | 115.58 KB |	110.34 KB |	33.39 KB |

#### ReScript

| Component   | Stat        | Parsed      | Gzipped     |
| ----------- | ----------- | ----------- | ----------- |
| No page components | 307.61 KB |	11.16 KB |	3.71 KB |
| + Home               | 328.29 KB |	23.10 KB |	6.61 KB |
| + Settings           | 349.47 KB |	28.87 KB |	7.9 KB |
| + Login              | 359.54 KB |	31.54 KB |	8.27 KB |
| + Register           | 372.67 KB |	34.76 KB |	8.53 KB |
| + CreateArticle      | 391.85 KB |	40.7 KB |	9.46 KB |
| + EditArticle        | 392.13 KB |	40.79 KB |	9.47 KB |
| + Article            | 420.48 KB |	53.96 KB |	11.71 KB |
| + Profile            | 435.64 KB |	60.08 KB |	12.69 KB |
| + Favorited          | 435.91 KB |	60.14 KB |	12.69 KB |

#### Chart

[See interactive version](https://jchavarri.github.io/jsoo-react-realworld-example-app/bundle-study/).

![js_of_ocaml-bundle-size-study-02.png](/media/js_of_ocaml-bundle-size-study-02.png)

## Caveats and learning

#### Runtime size

It is obvious from the analysis and data above that Js\_of\_ocaml runtime is larger than ReScript runtime (51.61 KB vs 11.16 KB in the "No page components" case). Part of this is due to design decisions of both compilers, and there is no way around it. Js\_of\_ocaml needs more code as it has to provide conversion functions for types like `string`, while ReScript just does not.

But I believe part of this runtime could be reduced with some careful optimizations. Js\_of\_ocaml runtime implementation could remove some of the runtime code if it stopped using functors for internal implementations (see more below), or gave more control to users over what functionality is available. For example, there is a "synthetic" file system available in Js\_of\_ocaml that gets included in bundle through the `caml_fs_init` call, but this is not required for the large majority of applications.

All this means that if your application has very tight bundle budgets, ReScript is probably the better choice here, or just plain JavaScript.

#### `Printf` module

Unlike the original ReScript application, where JSON encoders were written manually, in the Js\_of\_ocaml version I decided to use some tool to generate encoders and decoders to work with the JSON values that result from interacting with the public API.

After some research I decided to use [ppx\_jsobject\_conv](https://github.com/little-arhat/ppx_jsobject_conv). This tool turned out to be a great choice as it leverages all the infrastructure from ppxlib, so it is robust and very easy to use.

However, one small thing was that it made some usage of the `Printf` functions. `Printf` has a complex implementation, and it increases the bundle quite significantly. Fortunately, the usages in `ppx\_jsobject\_conv` were quite limited and [could be removed](https://github.com/little-arhat/ppx_jsobject_conv/pull/8) without much hassle.

In general, it is better avoid using `Printf` functions if the bundle size budget or a Js\_of\_ocaml is limited.

#### Functors

A noticeable increase can be seen in Js\_of\_ocaml bundle sizes when the `Article` component gets added to the bundle. While ReScript app only increases by ~13KB (from 40.79 to 53.96), the Js\_of\_ocaml app increased by ~21KB (from 83.05 to 104.55).

This bump is due to [a couple of modules](https://github.com/jchavarri/jsoo-react-realworld-example-app/blob/0138bfedddc1e57237ffe9a9a53a07aea9f73bf6/src/hook.ml#L2-L3) that are created using OCaml [Set.Make](https://ocaml.org/api/Set.Make.html) functor.

Apparently, [functors](https://ocaml.org/learn/tutorials/functors.html) can not be dead code eliminated by the compiler, so all the functions that are part of the `Set` module will appear in the resulting bundle, regardless if they are used or not. The good news is that the functions only appear once, so as an application grows more (and more functions from `Set` are used, and more times the functor is called), the cost would remain constant.

This problem is something that ReScript has tackled by re-implementing modules like `Set` in its standard library [Belt](https://rescript-lang.org/docs/manual/latest/api/belt/set) in a way that they don't use functors. Maybe something similar could be done for Js\_of\_ocaml.

#### Leverage existing browser APIs

One way to keep bundle size limited is to use the browser APIs, for the cases when they are available.

For example, instead of [lwt](https://github.com/ocsigen/lwt/), the project is using [promise-jsoo](https://github.com/mnxn/promise_jsoo).

Instead of [ppx\_yojson\_conv](https://github.com/janestreet/ppx_yojson_conv) the project uses the aforementioned [ppx\_jsobject\_conv](https://github.com/little-arhat/ppx_jsobject_conv). The advantages of the latter is that it delegates the parse step from string to JSON to browser APIs like [JSON.parse](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON/parse) or [Response.json](https://developer.mozilla.org/en-US/docs/Web/API/Response/json). This removes the need of bundling additional code, and most probably leads to faster applications, as browser implementors have optimized these functions heavily over time.

#### Functions with lots of optional values

One realization that might be surprising is that all applications of functions with optional labelled arguments get compiled to a bunch of zeroes separated by commas when the arguments that are unused.

This is fine for functions that take a few arguments, but in `jsoo-react` case there was a function to create style blocks that was taking more than 300 optional labelled arguments 😅. Ultimately, the issue was [solved](https://github.com/ml-in-barcelona/jsoo-react/issues/112) by changing the API to use a list, but it is still something that might be interesting to solve at the compiler level (e.g. if Js\_of\_ocaml supported some way of creating JavaScript objects inline, in a way that doesn't impact bundle size).

## Conclusion

The measurements show that Js\_of\_ocaml version of the application has a larger initial cost, due to the runtime being larger than that of ReScript.

It also shows some larger increases due to functors, that could be fixed using alternatives like some bindings to JavaScript objects, or a reimplementation of `Set` that does not use functors.

But otherwise, for most of the incremental steps, Js\_of\_ocaml shows mostly the same bundle size increases than ReScript. The bundle size of both apps remain in the same order of magnitude, and it would be expected that as more components are added to the application, the difference become smaller.

This study also shows that regardless which compiler is used to generate JavaScript, some tooling and potentially integrations with continuous integration pipelines are desired, as it is easy to suddenly get bundle size increases in unexpected ways.

We might revisit this study in the future to incorporate improvements, in which case we will add notes to this post.

---

I hope you enjoyed the study, if you want to share any other caveats that are missed, or there is anything inaccurate or that can be improve, feel free to reach out [on Twitter](https://twitter.com/javierwchavarri/).
