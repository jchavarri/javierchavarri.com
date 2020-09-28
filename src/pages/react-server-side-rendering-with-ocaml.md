---
title: "React server-side rendering with OCaml"
subtitle: "An experiment with TyXML and ReasonReact"
date: "2020-08-25"
imghero: "https://www.javierchavarri.com/media/react-server-side-rendering-with-ocaml-01.jpg"
tags:
  - "React"
  - "ReasonML"
  - "OCaml"
---

A while back, [Dan Abramov mentioned](https://twitter.com/dan_abramov/status/942859338472882176) —to my surprise— that it would be relatively easy to have React server side renderer implemented in a different language:

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">RDS doesn&#39;t need the reconciler so it&#39;s easy to rewrite by hand. It&#39;s about 1kloc: <a href="https://t.co/ymAiLBl2Il">https://t.co/ymAiLBl2Il</a></p>&mdash; Dan Abramov (@dan_abramov) <a href="https://twitter.com/dan_abramov/status/942859338472882176">December 18, 2017</a></blockquote>

I kept thinking about this regularly, and at some point started wondering how cool it would be to explore using OCaml to implement that server renderer.

But due to ~~laziness~~ lack of time, instead of rewriting it from scratch, I took an existing library that allows to build statically correct HTML ([TyXML](https://github.com/ocsigen/tyxml)), and used it to render HTML server side, that can later on be picked up and hydrated by [ReasonReact](https://reasonml.github.io/reason-react/).

The results from the experiment seem promising. The [Reason](reasonml.github.io/) syntax and the JSX extension for TyXML allow the same components to be shared across both server and client environments.

The experiment code is open source, and available in https://github.com/jchavarri/ocaml_webapp. The demo app can be accessed in https://ocaml-webapp.herokuapp.com/. All the pages in this demo app can be rendered by either the server or the client.

This blog post will go through the details on how the experiment went, what troubles were found along the way, and some of the solutions around them.

![/media/react-server-side-rendering-with-ocaml-01.jpg](/media/react-server-side-rendering-with-ocaml-01.jpg)

*Photo by [Linus Nylund](https://unsplash.com/@doto?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText) on [Unsplash](https://unsplash.com/wallpapers/nature/water?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText)*

## React hydration in a hurry

To summarize, "hydration" is a technique that allows React client-side applications to start faster, by assuming the HTML that React application code would produce was previously rendered by the server and returned as part of the page response.

Based on that assumption, React just needs to attach the event handlers to existing DOM elements, but does not have to create these elements (which is probably the slowest part in the initialization of a React application).

By not having to touch the DOM, the process of starting up a React application avoids a lot of work and can finish in less time.

#### A warning about hydration and performance

Hydration is a very nuanced topic and has several performance implications. For most real-world scenarios and applications, hydration will not be the optimal solution. There are other approaches that lead to better results in terms of performance:
- either do less work server side by rendering just enough HTML, and then have React components client-side do more rendering work
- or the other way around, do most of the work server side and spread small scripts in the client to add dynamic behavior.

However, hydration unlocks a great developer experience, so it is a very hyped topic at the moment. Companies like [Gatsby](gatsbyjs.org/) and [Vercel](https://vercel.com/) are innovating on it very quickly to work around these performance issues while keeping the same great development experience. 

If you're curious to know more, I recommend [the official documentation](https://reactjs.org/docs/react-dom.html#hydrate) and also the post ["Rendering on the Web"](https://developers.google.com/web/updates/2019/02/rendering-on-the-web) in the Google dev blog.

Now that we got these performance concerns out of the way, let's go back to the experiment.

## Why OCaml native and TyXML?

React server side rendering (SSR) and hydration is typically implemented using Node, and for good reasons.

There are limitations that are inherent to the "platform gap" between Node and the browser: components rendered in Node can't call methods or APIs available only on the browser (note the same happens in this experiment, between the OCaml native APIs and BuckleScript ones). This gap is not really obvious, and sometimes users of SSR frameworks like Gatsby [get confused by errors like `window is undefined`](https://github.com/gatsbyjs/gatsby/issues/12849). The solution involves doing runtime checks to see [if a given global is defined](https://www.gatsbyjs.com/docs/debugging-html-builds/#how-to-check-if-window-is-defined), and from there one can infer that is in one or another environment.

However! React is written in JavaScript, so Node applications that render React components can leverage a lot of previously existing libraries and tools from the extensive React and JavaScript ecosystems.

So, why this attempt to use a completely different language when all this exist already and works in JavaScript? Besides just for pure sake of experimenting 👨‍🔬 there are some other good reasons.

#### Speed 

One reason that makes worth explore rendering components with OCaml native is speed. OCaml binaries can start render some content and return it in an incredibly short time ([even less than 2ms!](https://twitter.com/_anmonteiro/status/1069738117647777792)), which makes them very appealing for serverless environments like lambda, which companies like Vercel [are migrating to due](https://vercel.com/blog/zeit-is-now-vercel) to their appeal for developers.

OCaml binaries also run pretty fast, and they do well in scenarios where a lot of short-lived small allocations are made (like with parsers or web servers). In general, one can trust that OCaml-generated binaries will run [fast](https://blog.chewxy.com/2019/02/20/go-is-average/).

#### Type system and safety

Another reason to experiment with OCaml to render components is the type system. TyXML is a library that allows to generate valid HTML. The way TyXML guarantees the validity is because it encodes in its implementation [the W3C rules](https://html.spec.whatwg.org/) for document validity.

For example, if you try to do this:

```reason
let t = <ul> <div /> </ul>;
```

The compiler will complain:

```
let t = <ul> <div /> </ul>;
        ˜˜˜˜˜˜˜˜˜˜˜˜˜˜˜˜˜˜
Type 'a = [> `Div ] is not compatible with type
  'b = [< `Li(Html_types.li_attrib) ] 
The second variant type does not allow tag(s) `Div
```

One can quickly realize that the only tag allowed inside `ul` is `li`. I _learn_ about HTML rules from TyXML while I'm coding, which is really amazing.

Note that React has similar invalid HTML detection mechanisms through [an internal function `validateDOMNesting`](https://github.com/facebook/react/blob/0b5a26a4895261894f04e50d5a700e83b9c0dcf6/packages/react-dom/src/__tests__/validateDOMNesting-test.js#L36), but there are two big differences:

- they only apply to nesting, while TyXML will also check attributes are valid
- more importantly, React checks are only done at runtime

As far as I know, neither TypeScript or Flow, or even ReasonReact, do this kind of static checks to make sure the resulting HTML is valid, although it seems that support for a similar mechanism [could be part of ReasonReact](https://github.com/reasonml/reason-react/pull/567) at some point.

## So how does a component look like?

Components rendered with TyXML can be adapted to look mostly like a ReasonReact component, with a few differences. Here's an example of a `Link.re` component from the demo application ([source](https://github.com/jchavarri/ocaml_webapp/blob/7e03cc30374e08788cccf0dd9f16eac65c48cca3/shared/Link.re)):

```reason
open Bridge;

let createElement = (~url, ~txt, ()) => {
  <a
    className="text-blue-500 hover:text-blue-800"
    href=url
    onClick={e => {
      ReactEvent.Mouse.preventDefault(e);
      ReasonReactRouter.push(url);
    }}>
    {React.string(txt)}
  </a>;
};

[@react.component]
let make = (~url, ~txt) => {
  createElement(~url, ~txt, ());
};
```

We will now go through the code of this sample component and see the challenges that the experiment brought up, before we can have a more seamless experience.

## Challenges

#### `createElement` vs `make`

This is the first and probably more obvious. TyXML offers a JSX ppx[^2]. In this ppx, the elements created from "uppercase" components convert to a call to `createElement`. For example:

```reason
let t = <Foo bar=2 />
/* will convert to: */
let t = Foo.createElement(~bar=2,());
```

While in ReasonReact, the JSX ppx makes a slighly different transformation, calling the `make` function inside the component module:

```reason
let t = <Foo bar=2 />
/* will convert to: */
let t = let t = React.createElement(Foo.make, Foo.makeProps(~bar=2, ()));
```

So how was this problem fixed? For now, each component exposes both `createElement` and `make`. Not the most elegant solution I know 😅, but probably this can be simplified in the future by bringing TyXML ppx behavior closer to what ReasonReact is expecting, in terms of naming.

#### Platform-dependent shims

Sometimes the components will need to call functions that are only available in one platform, for example, only in the browser or only on the server. To solve this, there was a small module call `bridge` that is available on both sides: [server](https://github.com/jchavarri/ocaml_webapp/blob/7e03cc30374e08788cccf0dd9f16eac65c48cca3/server/tyxml-reasonreact-bridge/bridge.ml) and [client](https://github.com/jchavarri/ocaml_webapp/blob/7e03cc30374e08788cccf0dd9f16eac65c48cca3/client/src/Bridge.re).

There are functions that are required to work around ReasonReact and TyXML handling things differently. For example, TyXML allows component children to be a list, but ReasonReact expects them to be a value of type `React.element`. So can have a function `React.list` that does nothing in TyXML, but calls the appropriate converters in ReasonReact (note there will be a performance cost for this conversion).

So, [in TyXML](https://github.com/jchavarri/ocaml_webapp/blob/7e03cc30374e08788cccf0dd9f16eac65c48cca3/server/tyxml-reasonreact-bridge/bridge.ml#L4) it would be something like[^3]:

```reason
module React = {
  ...
  let list = a => a;
};
```

And [in ReasonReact](https://github.com/jchavarri/ocaml_webapp/blob/7e03cc30374e08788cccf0dd9f16eac65c48cca3/client/src/Bridge.re#L3):

```reason
module React = {
  ...
  let list = el => el->Array.of_list->React.array;
};
```

React hooks are also part of this bridge. The functions in React API that allow to create hooks (like `useEffect` or `useMemo`) only get called _after_ the component has rendered. In the server, these components never really get mounted, we just need to get back the HTML after their render function is called.

So in the server, OCaml native can implement a shim for these functions that is part of the bridge (so the component code does not fail to build) but do nothing when they are called:

```reason
let useState: (unit => 'state) => ('state, ('state => 'state) => unit) =
  f => (f(), _ => ());

let useEffect0: (unit => option(unit => unit)) => unit = _ => ();
```

There are more examples in [the demo app](https://github.com/jchavarri/ocaml_webapp/blob/7e03cc30374e08788cccf0dd9f16eac65c48cca3/server/tyxml-reasonreact-bridge/bridge.ml#L6-L14).

#### Event handlers

Another interesting challenge involves React event handlers. By default, TyXML does not allow props like `onClick` to be passed to elements, as it has been designed originally with HTML attributes in mind. So any components using them will fail to compile.

The solution to this was to add [a small update](https://github.com/ocsigen/tyxml/commit/f3376134fbd51d50ca0720a8cddfb3919f570ea7) to TyXML ppx, so that it can handle props with React event handlers names. When the ppx finds one of these props, it will make the prop and the value passed with it disappear from the resulting code.

For example, the implementation of the `createElement` function in the `Link` component above was:

```reason
let createElement = (~url, ~txt, ()) => {
  <a
    className="text-blue-500 hover:text-blue-800"
    href=url
    onClick={e => {
      ReactEvent.Mouse.preventDefault(e);
      ReasonReactRouter.push(url);
    }}>
    {React.string(txt)}
  </a>;
};
```

In ReasonReact, it will remain like shown above. But in TyXML it will be:

```reason
let createElement = (~url, ~txt, ()) => {
  <a
    className="text-blue-500 hover:text-blue-800"
    href=url
    >
    {React.string(txt)}
  </a>;
};
```

This is also cool with regards to shims and platform-specific code. Because this attribute and its value gets eliminated, there is no need to add to shims or care about any code that goes inside it, as the native type checker will never see it.

## Conclusions and future work

So, while this prototype proves that it is possible to share some components code between environments and libraries as different as TyXML and ReasonReact, many challenges remain as seen above.

Some future work could involve:
- Continue updating TyXML to improve the integration with ReasonReact. For example, to process children in a similar way.
- Replicate in TyXML everything that React server side does to guarantee hydration. For example, print HTML comments between string children so the client knows how to hydrate them, add support for Suspense, etc.
- On the other direction, maybe extract the nice parts of TyXML HTML validation and make them available as a shared library that both TyXML and ReasonReact can consume.

---

I hope you enjoyed the post, check the demo app in https://github.com/jchavarri/ocaml_webapp, and if you want to share any feedback or have a suggestion, reach out [on Twitter](https://twitter.com/javierwchavarri/).

[^1]: Now [ReScript](https://reasonml.org/blog/bucklescript-is-rebranding).
[^2]: Pre-processor extension, see more about them [here](https://www.javierchavarri.com/when-two-worlds-collide-using-esy-to-consume-ppxs-from-source-in-bucklescript/#-what-is-a-ppx). 
[^3]: Note the server code in the demo app is written in OCaml syntax, translated below to Reason syntax.
