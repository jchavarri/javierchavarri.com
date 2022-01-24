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

Historically, I have noticed a recurring theme (meme) in both the OCaml and ReScript communities, where from time to time someone mentions that [Js\_of\_ocaml](http://ocsigen.org/js_of_ocaml/) —a compiler that converts from OCaml bytecode to JavaScript— generates large output files.

However, my experience was generally the opposite. When using Js\_of\_ocaml  [production builds](https://dune.readthedocs.io/en/stable/jsoo.html?highlight=mode%20release#separate-compilation) (using `--profile=prod` flag), the resulting JavaScript artifacts were quite thin, and the compiler applies very aggressive dead code optimizations.

Js\_of\_ocaml also requires some conversions between OCaml native types and JavaScript types, and the fact that it generates JavaScript from bytecode makes the whole process harder to track. Are Js\_of\_ocaml generated files really that large, as these rumors suggest? I did a small experiment to answer it.

![js_of_ocaml-bundle-size-study-01.jpg](/media/js_of_ocaml-bundle-size-study-01.jpg)
  
*Photo by [Kai Dahms](https://unsplash.com/@dilucidus) on [Unsplash](https://unsplash.com/s/photos/measure)*


## Hypothesis

The main theory that I wanted to proof is if Js\_of\_ocaml produces reasonably sized JavaScript files. I also was interested about tracking the evolution in size of these output files over time, as the application keeps being developed. If the output file is small for small apps, but grows too quickly over time, it would mean Js\_of\_ocaml would not be suitable for web applications that have limited bundle size budgets, and are developed for the long term.

In order to answer the above question, I thought it would be nice to use one of the most efficients compilers to JavaScript that exist out there: [ReScript](https://rescript-lang.org/). Which happens to be very close to OCaml as well :)


## The experiment

To run the experiment, I looked for an existing ReScript application that had some functionality that is common to most web applications:
- data fetching and pushing from/to servers
- routing
- parsing
- some work with collections and data processing

I found a good candidate in [jihchi/rescript-react-realworld-example-app](https://github.com/jihchi/rescript-react-realworld-example-app). This application uses ReScript and [rescript-react](https://rescript-lang.org/docs/react/latest/introduction) —the ReScript bindings to [React.js](https://reactjs.org/)— to build another example of the ["mother of all demo apps"](https://github.com/gothinkster/realworld). This demo app is a social blogging site (i.e. a Medium.com clone) that uses a custom API for all requests, including authentication.

To compare apples to apples, the plan would be to migrate this application fully to Js\_of\_ocaml, and use [jsoo-react](https://github.com/reason-in-barcelona/jsoo-react) to replicate the behavior on the original ReScript app.

Then, as the application consists on a dozen of screens or so, there would be an easy way produce JavaScript files and take measurements of the output files progressively, as new screens were added to the application.

## Methodology

Both bundles are generated using Webpack v4.46.0 in production mode, to avoid differences in minimization or bundling.
The bundles are analyzed using [webpack-bundle-analyzer](https://github.com/webpack-contrib/webpack-bundle-analyzer).

The results below were created on these commits:

- Js\_of\_ocaml: [0138bfe](https://github.com/jchavarri/jsoo-react-realworld-example-app/commit/0138bfedddc1e57237ffe9a9a53a07aea9f73bf6)
- ReScript: [2c319d9](https://github.com/jchavarri/jsoo-react-realworld-example-app/commit/2c319d933a5025e11a0a12e6abb804130659cc7a)

For each case, all components in [the main `App` component](https://github.com/jchavarri/jsoo-react-realworld-example-app/blob/0138bfedddc1e57237ffe9a9a53a07aea9f73bf6/src/app.ml#L20-L28) were comment, and then uncomment progressively, while running these commands on each step of the way:

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

![js_of_ocaml-bundle-size-study-02.png](/media/js_of_ocaml-bundle-size-study-02.png)

## Caveats and learning 

#### `Printf` module

Unlike the original ReScript application, in Js\_of\_ocaml I decided to use some tool to generate encoders and decoders to work with JSON to interact with the API.

After some research I decided to use [ppx\_jsobject\_conv](https://github.com/little-arhat/ppx_jsobject_conv). This tool turned out to be a great choice as it leverages all the infrastructure from ppxlib so it is very robust and friendly to use.

However, one small thing was that it used `Printf` module. `Printf` has a complex implementation, and it increases the bundle quite significantly. Fortunately, the usages in `ppx\_jsobject\_conv` were quite limited and [could be removed](https://github.com/little-arhat/ppx_jsobject_conv/pull/8).

In general, it is better avoid using `Printf` functions if the bundle size budget or a Js\_of\_ocaml is limited.

#### Functors

There is a noticeable increase when adding `Article` component to the bundle. While ReScript app only increases by ~13KB (from 40.79 to 53.96), the Js\_of\_ocaml app increased by ~21KB (from 83.05 to 104.55).

This bump is due to [a couple of modules](https://github.com/jchavarri/jsoo-react-realworld-example-app/blob/0138bfedddc1e57237ffe9a9a53a07aea9f73bf6/src/hook.ml#L2-L3) that are created using OCaml [Set.Make](https://ocaml.org/api/Set.Make.html) functor.

Apparently, functors can not be dead code eliminated by the compiler, so all the functions that are part of the `Set` module will appear in the resulting bundle, regardless if they are not used. The good news is that the functions only appear once, so as an application grows more (and more functions from `Set` are used, and more times the functor is called), the cost would remain the same.

This problem is something that ReScript has tackled by re-implementing modules like `Set` in its standard library [Belt](https://rescript-lang.org/docs/manual/latest/api/belt/set) in a way that they don't use functors. Maybe something similar could be done for Js\_of\_ocaml.

#### Leverage existing browser APIs

One way to keep bundle size limited was to use the browser APIs when they are available.

For example, instead of [lwt](https://github.com/ocsigen/lwt/), the project is using [promise-jsoo](https://github.com/mnxn/promise_jsoo).

Instead of [ppx_yojson_conv](https://github.com/janestreet/ppx_yojson_conv) the project uses the aforementioned [ppx\_jsobject\_conv](https://github.com/little-arhat/ppx_jsobject_conv). The advantages of the latter is that it allows to parse from string to JSON using browser APIs like [JSON.parse](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON/parse) or [Response.json](https://developer.mozilla.org/en-US/docs/Web/API/Response/json), which both removes the need of bundling additional code, and most probably leads to faster applications, as browser implementors have optimized these functions very heavily.

#### Functions with lots of optional values

Something that was interesting during this experiment is the realization that all applications of functions with optional labelled arguments get compiled to a bunch of zeroes separated by commas for the arguments that are unused.

This is alright for functions that take a few arguments, but in `jsoo-react` case there was a function to create style blocks that was taking more than 300 arguments 😅. Ultimately, the issue was [reported and solved](https://github.com/ml-in-barcelona/jsoo-react/issues/112), but it is still something that might be interesting to solve at the compiler level (e.g. if Js\_of\_ocaml supported some way of creating JavaScript objects in a way that doesn't impact bundle size).

#### PPXs

One really nice thing is that the whole OCaml [PPX](https://tarides.com/blog/2019-05-09-an-introduction-to-ocaml-ppx-ecosystem) ecosystem is available as first class citizen in Js\_of\_ocaml applications, without any impact on resulting JavaScript output size.

## Conclusion

The measurements show that Js\_of\_ocaml version of the application has a larger initial cost, due to the runtime being larger than that of ReScript.

It also shows some larger increases due to functors, that could be fixed using alternatives like some bindings to JavaScript objects, or a reimplementation of `Set` that does not use functors.

But for many of the incremental steps, Js\_of\_ocaml shows mostly the same bundle size increase than ReScript, which I think is very promising result. The bundle size of both apps remain in the same order of magnitude, and it would be expected that as more components are added to the application, the difference become smaller.

We might revisit this study in the future to incorporate improvements, in which case we will add notes to this post.

---

I hope you enjoyed the study, if you want to share any other caveats that are missed, or there is anything inaccurate or that can be improve, feel free to reach out [on Twitter](https://twitter.com/javierwchavarri/).
